package game

import (
	"fmt"
	"math"
	mrand "math/rand" // alias to avoid confusion with built-in rand
	"testing"
	"time"
)

func createTestGame() *Game {
	game := &Game{
		LevelChans: make([]chan *Level, 1),
		InputChan:  make(chan *Input),
		Levels:     make(map[string]*Level),
	}
	game.LevelChans[0] = make(chan *Level)
	game.CurrentLevel = createTestLevel()
	game.Levels["test"] = game.CurrentLevel
	return game
}

func createTestLevel() *Level {
	// Create a simple 3x3 test level
	level := &Level{}
	level.Map = make([][]Tile, 3)
	for i := range level.Map {
		level.Map[i] = make([]Tile, 3)
		for j := range level.Map[i] {
			level.Map[i][j] = Tile{Rune: DirtFloor}
		}
	}
	level.Player = &Player{
		Character: Character{
			Entity:       Entity{Pos: Pos{X: 1, Y: 1}, Name: "You", Rune: '@'},
			Hitpoints:    100,
			MaxStamina:   10,
			Stamina:      10,
			Speed:        1.0,
			ActionPoints: 0,
			SightRange:   10,
		},
	}
	level.Monsters = make(map[Pos]*Monster)
	level.Items = make(map[Pos][]*Item)
	level.Debug = make(map[Pos]bool)
	level.Battle = &Battle{}
	level.Events = make([]string, 10)
	level.Portals = make(map[Pos]*LevelPos)
	return level
}

func TestGameCreation(t *testing.T) {
	game := createTestGame()
	if game == nil {
		t.Fatal("Expected game to be created")
	}
	if game.InputChan == nil {
		t.Error("Input channel should be initialized")
	}
	if len(game.LevelChans) != 1 {
		t.Errorf("Expected 1 level channel, got %d", len(game.LevelChans))
	}
	if game.CurrentLevel == nil {
		t.Error("Current level should be initialized")
	}
}

func TestPos(t *testing.T) {
	p1 := Pos{X: 1, Y: 2}
	p2 := Pos{X: 1, Y: 2}
	if p1 != p2 {
		t.Error("Equal positions should be equal")
	}
}

func TestLevel(t *testing.T) {
	level := createTestLevel()
	if level.Player == nil {
		t.Fatal("Player should be initialized")
	}
	if len(level.Map) != 3 || len(level.Map[0]) != 3 {
		t.Error("Map size should be 3x3")
	}
	if level.Map[1][1].Rune != DirtFloor {
		t.Error("Center tile should be dirt floor")
	}
}

func TestMonsterCreation(t *testing.T) {
	pos := Pos{X: 1, Y: 1}
	rat := NewRat(pos)
	if rat.Pos != pos {
		t.Error("Rat position not set correctly")
	}
	if rat.Name != "Rat" {
		t.Error("Rat name not set correctly")
	}
	if rat.Hitpoints != 4 {
		t.Error("Rat hitpoints not set correctly")
	}

	spider := NewSpider(pos)
	if spider.Speed <= rat.Speed {
		t.Error("Spider should be faster than rat")
	}
}

func TestItemOperations(t *testing.T) {
	level := createTestLevel()
	pos := Pos{X: 0, Y: 0}

	// Test item creation and placement
	sword := NewSword(pos)
	level.Items[pos] = append(level.Items[pos], sword)

	if len(level.Items[pos]) != 1 {
		t.Error("Item should be added to level")
	}

	if level.Items[pos][0].Name != "Sword" {
		t.Error("Item name not set correctly")
	}

	// Test multiple items in same position
	potion := NewPotion(pos)
	level.Items[pos] = append(level.Items[pos], potion)

	if len(level.Items[pos]) != 2 {
		t.Error("Multiple items should be allowed in same position")
	}
}

func TestCharacterStats(t *testing.T) {
	player := &Player{
		Character: Character{
			Entity: Entity{
				Pos:  Pos{X: 1, Y: 1},
				Name: "TestPlayer",
				Rune: '@',
			},
			Hitpoints:    100,
			MaxStamina:   10,
			Stamina:      10,
			Speed:        1.0,
			ActionPoints: 0,
			SightRange:   10,
		},
	}

	if player.Hitpoints != 100 {
		t.Errorf("Expected hitpoints to be 100, got %d", player.Hitpoints)
	}

	// Test damage calculation
	oldHP := player.Hitpoints
	damage := 20
	player.Hitpoints -= damage
	if player.Hitpoints != oldHP-damage {
		t.Errorf("Expected hitpoints to be %d, got %d", oldHP-damage, player.Hitpoints)
	}

	// Test stamina limits
	if player.Stamina > player.MaxStamina {
		t.Error("Stamina should not exceed MaxStamina")
	}

	// Test position
	newPos := Pos{X: 2, Y: 2}
	player.Pos = newPos
	if player.Pos != newPos {
		t.Errorf("Expected position to be %v, got %v", newPos, player.Pos)
	}
}

// func TestBurst(t *testing.T) {
// 	burst := &Burst{
// 		Notes:    []int{1, 2, 3, 4},
// 		MaxCombo: 4,
// 		Combo:    0,
// 	}

// 	if len(burst.Notes) != burst.MaxCombo {
// 		t.Errorf("Expected notes length %d to match MaxCombo %d", len(burst.Notes), burst.MaxCombo)
// 	}

// 	// Test combo increment
// 	burst.Combo++
// 	if burst.Combo != 1 {
// 		t.Error("Failed to increment combo")
// 	}

// 	// Test combo max limit
// 	for i := 0; i < burst.MaxCombo; i++ {
// 		burst.Combo++
// 	}
// 	if burst.Combo > burst.MaxCombo {
// 		t.Error("Combo should not exceed MaxCombo")
// 	}

// 	// Test combo reset
// 	burst.Combo = 0
// 	if burst.Combo != 0 {
// 		t.Error("Failed to reset combo")
// 	}
// }

func TestMonsterBehavior(t *testing.T) {
	level := createTestLevel()
	monsterPos := Pos{X: 0, Y: 0}
	rat := NewRat(monsterPos)
	level.Monsters[monsterPos] = rat

	// Test monster initialization in level
	if monster, exists := level.Monsters[monsterPos]; !exists {
		t.Error("Monster should exist in level")
	} else if monster.Name != "Rat" {
		t.Error("Monster name mismatch")
	}

	// Test monster stamina mechanics
	initialStamina := rat.Stamina
	rat.Stamina--
	if rat.Stamina != initialStamina-1 {
		t.Error("Failed to decrease monster stamina")
	}

	// Test monster items
	if len(rat.Items) == 0 {
		t.Error("Monster should have items")
	}

	// Test action points accumulation
	initialAP := rat.ActionPoints
	rat.ActionPoints += rat.Speed
	if rat.ActionPoints != initialAP+rat.Speed {
		t.Error("Failed to accumulate action points")
	}
}

func TestLevelOperations(t *testing.T) {
	level := createTestLevel()

	// Test tile visibility
	pos := Pos{X: 1, Y: 1}
	level.Map[pos.Y][pos.X].Visible = true
	if !level.Map[pos.Y][pos.X].Visible {
		t.Error("Tile should be visible")
	}

	// Test tile seen state
	level.Map[pos.Y][pos.X].Seen = true
	if !level.Map[pos.Y][pos.X].Seen {
		t.Error("Tile should be marked as seen")
	}

	// Test event handling
	event := "Test event"
	level.Events = append(level.Events, event)
	if len(level.Events) == 0 || level.Events[len(level.Events)-1] != event {
		t.Error("Failed to add event to level")
	}

	// Test debug position marking
	debugPos := Pos{X: 0, Y: 0}
	level.Debug[debugPos] = true
	if !level.Debug[debugPos] {
		t.Error("Failed to mark debug position")
	}
}

func TestTileOperations(t *testing.T) {
	level := createTestLevel()
	pos := Pos{X: 0, Y: 0}

	// Test wall placement
	level.Map[pos.Y][pos.X].Rune = StoneWall
	if level.Map[pos.Y][pos.X].Rune != StoneWall {
		t.Error("Failed to place stone wall")
	}

	// Test door states
	level.Map[pos.Y][pos.X].Rune = ClosedDoor
	if level.Map[pos.Y][pos.X].Rune != ClosedDoor {
		t.Error("Failed to place closed door")
	}

	level.Map[pos.Y][pos.X].Rune = OpenDoor
	if level.Map[pos.Y][pos.X].Rune != OpenDoor {
		t.Error("Failed to change door state to open")
	}
}

func TestBattle(t *testing.T) {
	char1 := &Character{
		Entity: Entity{
			Pos:  Pos{X: 0, Y: 0},
			Name: "Fighter1",
			Rune: '@',
		},
		Hitpoints:    100,
		MaxStamina:   10,
		Stamina:      10,
		Speed:        1.0,
		ActionPoints: 0,
		SightRange:   10,
	}

	char2 := &Character{
		Entity: Entity{
			Pos:  Pos{X: 1, Y: 1},
			Name: "Fighter2",
			Rune: 'R',
		},
		Hitpoints:    50,
		MaxStamina:   5,
		Stamina:      5,
		Speed:        1.5,
		ActionPoints: 0,
		SightRange:   8,
	}

	battle := &Battle{
		C1: char1,
		C2: char2,
	}

	// Test battle initialization
	if battle.C1 == nil || battle.C2 == nil {
		t.Error("Battle characters should be initialized")
	}

	// Test character stats in battle
	if battle.C1.Stamina > battle.C1.MaxStamina {
		t.Error("Character 1 stamina exceeds maximum")
	}
	if battle.C2.Stamina > battle.C2.MaxStamina {
		t.Error("Character 2 stamina exceeds maximum")
	}
}

func TestInputHandling(t *testing.T) {
	game := createTestGame()

	// Test movement input handling
	initialPos := game.CurrentLevel.Player.Pos
	moveInput := &Input{Typ: Up}

	// Simulate processing the input
	if moveInput.Typ == Up {
		newPos := Pos{X: initialPos.X, Y: initialPos.Y - 1}
		if game.CurrentLevel.Map[newPos.Y][newPos.X].Rune == DirtFloor {
			game.CurrentLevel.Player.Pos = newPos
			if game.CurrentLevel.Player.Pos.Y >= initialPos.Y {
				t.Error("Player should have moved up")
			}
		}
	}

	// Test item handling
	itemPos := game.CurrentLevel.Player.Pos
	potion := NewPotion(itemPos)
	itemInput := &Input{
		Typ:  TakeItem,
		Item: potion,
	}

	// Simulate processing the item input
	if itemInput.Typ == TakeItem && itemInput.Item != nil {
		initialItems := len(game.CurrentLevel.Player.Items)
		game.CurrentLevel.Player.Items = append(game.CurrentLevel.Player.Items, itemInput.Item)
		if len(game.CurrentLevel.Player.Items) != initialItems+1 {
			t.Error("Failed to add item to player inventory")
		}
	}

	// Test quit game input
	quitInput := &Input{Typ: QuitGame}
	if quitInput.Typ != QuitGame {
		t.Error("Failed to create quit game input")
	}
}

func TestPatternGeneration(t *testing.T) {
	burst := &Burst{
		Notes:    make([]int, 4),
		MaxCombo: 4,
		Combo:    0,
	}

	// Test note pattern generation
	for i := range burst.Notes {
		burst.Notes[i] = i + 1
		if burst.Notes[i] != i+1 {
			t.Errorf("Expected note %d at position %d", i+1, i)
		}
	}

	// Test combo system
	burst.Combo = 2
	if burst.Combo > burst.MaxCombo {
		t.Error("Combo should not exceed MaxCombo")
	}

	// Test combo reset
	burst.Combo = 0
	if burst.Combo != 0 {
		t.Error("Failed to reset combo")
	}
}

func TestLevelGeneration(t *testing.T) {
	level := createTestLevel()

	// Test map boundaries
	mapHeight := len(level.Map)
	mapWidth := len(level.Map[0])

	// Add walls around the border
	for x := 0; x < mapWidth; x++ {
		level.Map[0][x].Rune = StoneWall
		level.Map[mapHeight-1][x].Rune = StoneWall
	}
	for y := 0; y < mapHeight; y++ {
		level.Map[y][0].Rune = StoneWall
		level.Map[y][mapWidth-1].Rune = StoneWall
	}

	// Verify walls
	if level.Map[0][0].Rune != StoneWall {
		t.Error("Border should be wall")
	}

	// Test portal placement
	portalPos := Pos{X: 1, Y: 1}
	targetLevel := createTestLevel()
	targetPos := Pos{X: 2, Y: 2}

	level.Portals = make(map[Pos]*LevelPos)
	level.Portals[portalPos] = &LevelPos{
		Level: targetLevel,
		Pos:   targetPos,
	}

	if portal, exists := level.Portals[portalPos]; !exists {
		t.Error("Portal should exist at position")
	} else if portal.Pos != targetPos {
		t.Error("Portal target position mismatch")
	}
}

func TestMonsterAI(t *testing.T) {
	level := createTestLevel()

	// Place monster
	monsterPos := Pos{X: 1, Y: 1}
	monster := NewRat(monsterPos)
	level.Monsters[monsterPos] = monster

	// Test monster pass action (should subtract speed from action points)
	monster.ActionPoints = 2.0
	expectedAP := monster.ActionPoints - monster.Speed
	monster.Pass()
	if monster.ActionPoints != expectedAP {
		t.Errorf("Expected AP to be %f after pass, got %f", expectedAP, monster.ActionPoints)
	}

	// Test action points accumulation during update
	monster.ActionPoints = 0
	monster.Update(level)
	if monster.ActionPoints != monster.Speed {
		t.Errorf("Expected AP to be %f after update, got %f", monster.Speed, monster.ActionPoints)
	}

	// Test monster autoplay during battle
	monster.Stamina = 3
	monster.Burst = &Burst{
		Notes:    []int{1, 2, 3},
		MaxCombo: 3,
		Combo:    0,
	}

	// Set up battle conditions
	level.LastEvent = Attack
	initialNotes := len(monster.Burst.Notes)
	initialStamina := monster.Stamina

	// Run autoplay in a goroutine with timeout
	done := make(chan bool)
	go func() {
		monster.Autoplay(level)
		done <- true
	}()

	// Wait for autoplay to finish or timeout
	select {
	case <-done:
		if len(monster.Burst.Notes) >= initialNotes {
			t.Error("Monster should have consumed notes during autoplay")
		}
		if monster.Stamina >= initialStamina {
			t.Error("Monster should have consumed stamina during autoplay")
		}
	case <-time.After(2 * time.Second):
		t.Error("Autoplay test timed out")
	}
}

func TestGameEvents(t *testing.T) {
	level := createTestLevel()

	// Clear existing events
	level.Events = make([]string, 10)
	level.EventPos = 0

	// Test move event
	level.AddEvent("Player moved to (1,1)")
	if level.Events[0] != "Player moved to (1,1)" {
		t.Error("Move event should be recorded")
	}

	// Test door event
	doorPos := Pos{X: 1, Y: 0}
	level.Map[doorPos.Y][doorPos.X].Rune = ClosedDoor
	level.Map[doorPos.Y][doorPos.X].Rune = OpenDoor
	level.AddEvent("Door opened at (1,0)")

	if level.Map[doorPos.Y][doorPos.X].Rune != OpenDoor {
		t.Error("Door state should be open")
	}
	if level.Events[1] != "Door opened at (1,0)" {
		t.Error("Door event not recorded correctly")
	}

	// Test attack event
	level.AddEvent("Player attacked monster")
	if level.Events[2] != "Player attacked monster" {
		t.Error("Attack event not recorded correctly")
	}

	// Test event position wrapping
	for i := 0; i < 15; i++ {
		level.AddEvent(fmt.Sprintf("Event %d", i))
	}
	if level.EventPos >= len(level.Events) {
		t.Error("Event position should wrap around")
	}
}

func TestPathfinding(t *testing.T) {
	level := createTestLevel()

	// Create a more complex map for pathfinding
	level.Map = make([][]Tile, 5)
	for i := range level.Map {
		level.Map[i] = make([]Tile, 5)
		for j := range level.Map[i] {
			level.Map[i][j] = Tile{Rune: DirtFloor}
		}
	}

	// Add some walls to test pathfinding
	level.Map[2][1].Rune = StoneWall
	level.Map[2][2].Rune = StoneWall
	level.Map[2][3].Rune = StoneWall

	start := Pos{X: 1, Y: 1}
	end := Pos{X: 3, Y: 3}

	path := level.astar(start, end)

	if len(path) == 0 {
		t.Error("Path should exist between start and end positions")
	}

	// Verify path doesn't go through walls
	for _, pos := range path {
		if level.Map[pos.Y][pos.X].Rune == StoneWall {
			t.Error("Path should not go through walls")
		}
	}

	// Test path length is reasonable (Manhattan distance)
	expectedMinLength := int(math.Abs(float64(end.X-start.X)) + math.Abs(float64(end.Y-start.Y)))
	if len(path) < expectedMinLength {
		t.Errorf("Path length %d is shorter than minimum possible length %d", len(path), expectedMinLength)
	}
}

func TestEntityCollision(t *testing.T) {
	level := createTestLevel()
	level.Battle = &Battle{} // Initialize battle field

	// Test monster collision with wall
	monsterPos := Pos{X: 0, Y: 0}
	monster := NewRat(monsterPos)
	level.Monsters[monsterPos] = monster

	// Place wall
	wallPos := Pos{X: 1, Y: 0}
	level.Map[wallPos.Y][wallPos.X].Rune = StoneWall

	// Try to move monster into wall
	oldPos := monster.Pos
	if canWalk(level, wallPos) {
		monster.Move(wallPos, level)
		if monster.Pos != oldPos {
			t.Error("Monster should not move through walls")
		}
	}

	// Test monster collision with player
	playerPos := Pos{X: 2, Y: 0}
	level.Player.Pos = playerPos
	level.LastEvent = NoEvent

	// Move monster towards player
	monster.Move(playerPos, level)

	if level.LastEvent != Attack {
		t.Error("Moving monster into player should trigger attack")
	}

	if level.Battle.C1 == nil || level.Battle.C2 == nil {
		t.Error("Battle participants should be set")
	}
}

func TestDamageResolution(t *testing.T) {
	level := createTestLevel()

	// Setup combatants with proper initialization
	attacker := &Character{
		Entity:     Entity{Name: "Attacker"},
		Hitpoints:  100,
		Stamina:    10,
		PatternRNG: mrand.New(mrand.NewSource(time.Now().UnixNano())),
	}

	defender := &Character{
		Entity:    Entity{Name: "Defender"},
		Hitpoints: 50,
		Stamina:   5,
	}

	// Set up battle context
	level.Battle.C1 = attacker
	level.Battle.C2 = defender

	// Test attack initiation
	level.Attack(attacker, defender)

	if level.LastEvent != Attack {
		t.Error("Attack should set LastEvent to Attack")
	}

	if attacker.Burst == nil {
		t.Error("Attacker should have a burst pattern generated")
	}

	// Test damage resolution
	initialHP := defender.Hitpoints
	level.ResolveDamage()

	if defender.Hitpoints >= initialHP {
		t.Error("Damage resolution should reduce defender hitpoints")
	}

	// Verify battle event was logged
	foundEvent := false
	for _, event := range level.Events {
		if event != "" && (event == "Attacker hit the Defender for 1 damage." ||
			event == "The Attacker hits you for 1 damage.") {
			foundEvent = true
			break
		}
	}
	if !foundEvent {
		t.Error("Damage event should be logged")
	}
}

func TestLineOfSight(t *testing.T) {
	level := createTestLevel()

	// Create a map for line of sight testing
	mapSize := 7
	level.Map = make([][]Tile, mapSize)
	for i := range level.Map {
		level.Map[i] = make([]Tile, mapSize)
		for j := range level.Map[i] {
			level.Map[i][j] = Tile{Rune: DirtFloor, Visible: false, Seen: false}
		}
	}

	// Place player in center with limited sight range
	centerPos := Pos{X: 3, Y: 3}
	level.Player.Pos = centerPos
	level.Player.SightRange = 2

	// Add some walls within sight range
	wallPositions := []Pos{
		{X: 2, Y: 2},
		{X: 4, Y: 4},
	}
	for _, pos := range wallPositions {
		if pos.X >= 0 && pos.X < mapSize && pos.Y >= 0 && pos.Y < mapSize {
			level.Map[pos.Y][pos.X].Rune = StoneWall
		}
	}

	// Calculate line of sight
	level.lineOfSight()

	// Test cases
	testCases := []struct {
		pos      Pos
		expected bool
		desc     string
	}{
		{centerPos, true, "Player position should be visible"},
		{Pos{X: 3, Y: 4}, true, "Adjacent tile should be visible"},
		{Pos{X: 1, Y: 1}, false, "Tile behind wall should not be visible"},
		{Pos{X: 6, Y: 6}, false, "Tile beyond sight range should not be visible"},
	}

	for _, tc := range testCases {
		if tc.pos.X >= 0 && tc.pos.X < mapSize && tc.pos.Y >= 0 && tc.pos.Y < mapSize {
			actual := level.Map[tc.pos.Y][tc.pos.X].Visible
			if actual != tc.expected {
				t.Errorf("%s: expected visibility %v at pos {%d,%d}, got %v",
					tc.desc, tc.expected, tc.pos.X, tc.pos.Y, actual)
			}
		}
	}

	// Test seen tiles persistence
	visiblePos := Pos{X: 3, Y: 4}
	if visiblePos.X >= 0 && visiblePos.X < mapSize && visiblePos.Y >= 0 && visiblePos.Y < mapSize {
		level.Map[visiblePos.Y][visiblePos.X].Visible = true
		level.Map[visiblePos.Y][visiblePos.X].Seen = true

		// Move player away and recalculate
		level.Player.Pos = Pos{X: 1, Y: 1}
		level.lineOfSight()

		if !level.Map[visiblePos.Y][visiblePos.X].Seen {
			t.Error("Previously seen tiles should remain seen")
		}
	}
}

func TestAdvancedPathfinding(t *testing.T) {
	level := createTestLevel()

	// Create a more complex map for pathfinding
	level.Map = make([][]Tile, 7)
	for i := range level.Map {
		level.Map[i] = make([]Tile, 7)
		for j := range level.Map[i] {
			level.Map[i][j] = Tile{Rune: DirtFloor}
		}
	}

	// Create a maze-like structure
	// #######
	// #S    #
	// ### ###
	// #     #
	// # ### #
	// #    E#
	// #######

	// Add walls
	for x := 0; x < 7; x++ {
		level.Map[0][x].Rune = StoneWall // Top wall
		level.Map[6][x].Rune = StoneWall // Bottom wall
	}
	for y := 0; y < 7; y++ {
		level.Map[y][0].Rune = StoneWall // Left wall
		level.Map[y][6].Rune = StoneWall // Right wall
	}

	// Add internal walls
	for x := 0; x < 3; x++ {
		level.Map[2][x].Rune = StoneWall
	}
	for x := 4; x < 7; x++ {
		level.Map[2][x].Rune = StoneWall
	}
	for x := 1; x < 4; x++ {
		level.Map[4][x].Rune = StoneWall
	}

	start := Pos{X: 1, Y: 1}
	end := Pos{X: 5, Y: 5}

	path := level.astar(start, end)

	if len(path) == 0 {
		t.Error("Should find path through maze")
	}

	// Verify path doesn't go through walls
	for _, pos := range path {
		if level.Map[pos.Y][pos.X].Rune == StoneWall {
			t.Error("Path should not go through walls")
		}
	}

	// Test no path possible
	level.Map[3][1].Rune = StoneWall
	level.Map[3][2].Rune = StoneWall
	level.Map[3][3].Rune = StoneWall
	level.Map[3][4].Rune = StoneWall
	level.Map[3][5].Rune = StoneWall

	noPath := level.astar(start, end)
	if len(noPath) > 0 {
		t.Error("Should not find path when none exists")
	}
}

func TestPatternMechanics(t *testing.T) {
	// Test pattern generation
	char := &Character{
		Entity:     Entity{Name: "Test Character"},
		PatternRNG: mrand.New(mrand.NewSource(42)), // Fixed seed for deterministic testing
	}
	// Test stream generation with different lengths
	for _, length := range []int{1, 4, 8} {
		stream := char.MakeStream(length)
		if len(stream) != length {
			t.Errorf("Expected stream length %d, got %d", length, len(stream))
		}
		// Verify each note is within valid range (0 to NumKeys-1)
		for i, note := range stream {
			if note < 0 || note >= NumKeys {
				t.Errorf("Note at position %d is out of valid range: %d", i, note)
			}
		}
	}

	// Test burst mechanics with combo logic
	burst := &Burst{
		Notes:    []int{1, 2, 3, 0},
		MaxCombo: 4,
		Combo:    0,
	}

	// Test combo tracking
	for i := 0; i < burst.MaxCombo; i++ {
		initialCombo := burst.Combo
		burst.Combo++
		if burst.Combo != minInt(initialCombo+1, burst.MaxCombo) {
			t.Errorf("Expected combo to be %d, got %d", minInt(initialCombo+1, burst.MaxCombo), burst.Combo)
		}
	}

	// Test note consumption
	initialNotes := len(burst.Notes)
	if len(burst.Notes) > 0 {
		burst.Notes = burst.Notes[1:]
		if len(burst.Notes) != initialNotes-1 {
			t.Error("Failed to consume note from burst")
		}
	}

	// Test empty burst handling
	emptyBurst := &Burst{
		Notes:    []int{},
		MaxCombo: 4,
		Combo:    0,
	}
	if len(emptyBurst.Notes) > 0 {
		t.Error("Empty burst should have no notes")
	}
}

func TestComplexBattleSequence(t *testing.T) {
	level := createTestLevel()

	// Setup attacker with pattern generation
	attacker := &Character{
		Entity:     Entity{Name: "Attacker"},
		Hitpoints:  100,
		Stamina:    10,
		MaxStamina: 10,
		PatternRNG: mrand.New(mrand.NewSource(42)),
	}

	// Setup defender
	defender := &Character{
		Entity:     Entity{Name: "Defender"},
		Hitpoints:  20,
		Stamina:    5,
		MaxStamina: 5,
	}

	// Initialize battle
	level.Battle = &Battle{C1: attacker, C2: defender}
	level.Attack(attacker, defender)

	if attacker.Burst == nil {
		t.Fatal("Battle should generate burst pattern")
	}

	if level.LastEvent != Attack {
		t.Error("Battle should start with Attack event")
	}

	// Simulate rhythm sequence
	initialNotes := len(attacker.Burst.Notes)
	for i := 0; i < 3 && len(attacker.Burst.Notes) > 0; i++ {
		// Simulate successful hit
		attacker.Burst.Combo++
		attacker.Burst.Notes = attacker.Burst.Notes[1:]
		attacker.Stamina--

		// Resolve damage after some hits
		if i == 2 {
			level.ResolveDamage()
			if defender.Hitpoints >= 20 {
				t.Error("Damage should be applied after successful hits")
			}
		}
	}

	// Verify notes were consumed
	if len(attacker.Burst.Notes) >= initialNotes {
		t.Error("Notes should be consumed during battle sequence")
	}

	// Test stamina consumption
	if attacker.Stamina >= 10 {
		t.Error("Stamina should be consumed during battle")
	}

	// Test battle event logging
	hasAttackEvent := false
	hasDamageEvent := false
	for _, event := range level.Events {
		if event != "" {
			if event == "The Attacker attacks you." {
				hasAttackEvent = true
			} else if event == "The Attacker hits you for 1 damage." {
				hasDamageEvent = true
			}
		}
	}
	if !hasAttackEvent {
		t.Error("Attack event should be logged")
	}
	if !hasDamageEvent {
		t.Error("Damage event should be logged")
	}
}

// func TestItemManipulation(t *testing.T) {
// 	level := createTestLevel()
// 	player := level.Player
// 	pos := player.Pos

// 	// Test taking and equipping weapon
// 	sword := NewSword(pos)
// 	level.Items[pos] = append(level.Items[pos], sword)
// 	level.MoveItem(sword, &player.Character)
// 	equip(&player.Character, sword)

// 	// Test taking and equipping armor
// 	helmet := NewHelmet(pos)
// 	level.Items[pos] = append(level.Items[pos], helmet)
// 	level.MoveItem(helmet, &player.Character)
// 	equip(&player.Character, helmet)

// 	// Verify equipment state
// 	if player.Weapon != sword {
// 		t.Error("Sword should be equipped")
// 	}
// 	if player.Helmet != helmet {
// 		t.Error("Helmet should be equipped")
// 	}

// 	// Test dropping equipped items
// 	level.DropItem(sword, &player.Character)
// 	if player.Weapon != nil {
// 		t.Error("Weapon should be unequipped after dropping")
// 	}
// 	if len(level.Items[pos]) == 0 {
// 		t.Error("Dropped item should be on ground")
// 	}

// 	// Test multiple items on ground
// 	potion := NewPotion(pos)
// 	credits := NewCredits(pos)
// 	level.Items[pos] = append(level.Items[pos], potion, credits)

// 	// Test taking multiple items
// 	for _, item := range []struct {
// 		item *Item
// 		name string
// 	}{
// 		{potion, "potion"},
// 		{credits, "credits"},
// 	} {
// 		initialInventorySize := len(player.Items)
// 		level.MoveItem(item.item, &player.Character)
// 		if len(player.Items) != initialInventorySize+1 {
// 			t.Errorf("%s should be added to inventory", item.name)
// 		}
// 	}

// 	// Test error cases in separate functions to handle panics
// 	t.Run("Drop non-existent item", func(t *testing.T) {
// 		defer func() {
// 			if r := recover(); r == nil {
// 				t.Error("Dropping non-existent item should panic")
// 			}
// 		}()
// 		nonExistentItem := NewSword(pos)
// 		level.DropItem(nonExistentItem, &player.Character)
// 	})

// 	t.Run("Move non-existent item", func(t *testing.T) {
// 		defer func() {
// 			if r := recover(); r == nil {
// 				t.Error("Moving non-existent item should panic")
// 			}
// 		}()
// 		invalidPos := Pos{X: 999, Y: 999}
// 		invalidItem := NewSword(invalidPos)
// 		level.MoveItem(invalidItem, &player.Character)
// 	})
// }

func TestMonsterAbilities(t *testing.T) {
	level := createTestLevel()

	// Test monster with different attributes
	boss := &Monster{
		Character: Character{
			Entity:     Entity{Name: "Boss", Pos: Pos{X: 2, Y: 2}},
			Hitpoints:  20,
			MaxStamina: 8,
			Stamina:    8,
			Speed:      2.0,
		},
	}
	level.Monsters[boss.Pos] = boss

	// Test pattern generation
	boss.PatternRNG = mrand.New(mrand.NewSource(42))
	streamLength := 5
	stream := boss.MakeStream(streamLength)
	if len(stream) != streamLength {
		t.Errorf("Expected burst length %d, got %d", streamLength, len(stream))
	}

	// Test attack power
	initialHP := level.Player.Hitpoints
	level.Battle = &Battle{C1: &boss.Character, C2: &level.Player.Character}
	level.Attack(&boss.Character, &level.Player.Character)
	level.ResolveDamage()
	if level.Player.Hitpoints >= initialHP {
		t.Error("Boss attack should deal damage")
	}
}

// Helper function to get minimum of two integers
func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
