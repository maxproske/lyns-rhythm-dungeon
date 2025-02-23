package game

import (
	"fmt"
	"math"
	mrand "math/rand" // alias to avoid confusion with built-in rand
	"os"
	"path/filepath"
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
	// Create 15x15 test level to match sight range
	level := &Level{}
	level.Map = make([][]Tile, 15)
	for i := range level.Map {
		level.Map[i] = make([]Tile, 15)
		for j := range level.Map[i] {
			level.Map[i][j] = Tile{Rune: DirtFloor}
		}
	}
	level.Player = &Player{
		Character: Character{
			Entity: Entity{
				Pos:  Pos{X: 7, Y: 7},
				Name: "You",
				Rune: '@',
			},
			Hitpoints:  100,
			SightRange: 3,
			Speed:      1.0,
			PatternRNG: mrand.New(mrand.NewSource(0)),
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
	if len(level.Map) != 15 || len(level.Map[0]) != 15 {
		t.Error("Map size should be 15x15")
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
	monster.Speed = 1.5 // Explicitly set speed for test
	monster.ActionPoints = 0
	monster.Update(level)
	if monster.ActionPoints != 0.5 {
		t.Errorf("Expected AP to be 0.500000 after update, got %f", monster.ActionPoints)
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

func TestNewGame(t *testing.T) {
	// Create a temporary directory for test files
	tmpDir := t.TempDir()

	// Create the game/maps directory structure
	mapsDir := filepath.Join(tmpDir, "game", "maps")
	if err := os.MkdirAll(mapsDir, 0755); err != nil {
		t.Fatalf("Failed to create maps directory: %v", err)
	}

	// Create a test world.txt file
	worldContent := "test_level"
	worldPath := filepath.Join(mapsDir, "world.txt")
	if err := os.WriteFile(worldPath, []byte(worldContent), 0644); err != nil {
		t.Fatalf("Failed to write world.txt: %v", err)
	}

	// Create a test map file
	mapContent := "####\n#@.#\n####"
	mapPath := filepath.Join(mapsDir, "test_level.map")
	if err := os.WriteFile(mapPath, []byte(mapContent), 0644); err != nil {
		t.Fatalf("Failed to write test_level.map: %v", err)
	}

	// Save the original working directory and switch to the temp dir
	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current working directory: %v", err)
	}
	defer os.Chdir(oldWd) // Ensure we revert to the original directory after the test

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change working directory: %v", err)
	}

	// Now initialize the game which should use our test files
	game := NewGame(1)
	if game == nil {
		t.Fatal("NewGame returned nil")
	}

	// Perform assertions
	if len(game.LevelChans) != 1 {
		t.Errorf("Expected 1 level channel, got %d", len(game.LevelChans))
	}
	if game.CurrentLevel == nil {
		t.Error("CurrentLevel is nil")
	}
	if game.CurrentLevel.Player == nil {
		t.Error("Player not initialized in CurrentLevel")
	}

	// Verify line of sight calculation
	visibleFound := false
	for y, row := range game.CurrentLevel.Map {
		for x := range row {
			if game.CurrentLevel.Map[y][x].Visible {
				visibleFound = true
				break
			}
		}
		if visibleFound {
			break
		}
	}
	if !visibleFound {
		t.Error("Line of sight not calculated, no visible tiles")
	}
}

func TestDropItem(t *testing.T) {
	level := createTestLevel()
	player := level.Player
	item := NewSword(Pos{})
	player.Items = append(player.Items, item)
	pos := player.Pos
	level.DropItem(item, &player.Character)

	if len(player.Items) != 0 {
		t.Error("DropItem should remove the item from the character's inventory")
	}
	if len(level.Items[pos]) != 1 {
		t.Error("DropItem should add the item to the level's items")
	}
	expectedEvent := "You dropped 1x Sword"
	found := false
	for _, e := range level.Events {
		if e == expectedEvent {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Event '%s' not found", expectedEvent)
	}
}

func TestMoveItem(t *testing.T) {
	level := createTestLevel()
	player := level.Player
	pos := player.Pos
	item := NewSword(pos)
	level.Items[pos] = append(level.Items[pos], item)

	level.MoveItem(item, &player.Character)

	if len(level.Items[pos]) != 0 {
		t.Error("MoveItem should remove the item from the level")
	}
	if len(player.Items) != 1 {
		t.Error("MoveItem should add the item to the player's inventory")
	}
	expectedEvent := "You picked up 1x Sword"
	found := false
	for _, e := range level.Events {
		if e == expectedEvent {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Event '%s' not found", expectedEvent)
	}
}

func TestKillMonster(t *testing.T) {
	level := createTestLevel()
	pos := Pos{X: 0, Y: 0}

	// Create monster and add item
	monster := NewRat(pos)
	initialItems := len(monster.Items)                   // Check default items
	monster.Items = append(monster.Items, NewBones(pos)) // Add test item

	level.Monsters[pos] = monster
	level.Kill(&monster.Character)

	// Verify monster removed
	if _, exists := level.Monsters[pos]; exists {
		t.Error("Monster should be removed from the level")
	}

	// Check total items dropped (default + test item)
	expectedItems := initialItems + 1
	if len(level.Items[pos]) != expectedItems {
		t.Errorf("Expected %d items dropped, got %d", expectedItems, len(level.Items[pos]))
	}
}

func TestKillPlayer(t *testing.T) {
	level := createTestLevel()
	player := level.Player
	player.Items = append(player.Items, NewSword(player.Pos))

	level.Kill(&player.Character)

	if player.Speed != 0 {
		t.Error("Player speed should be 0 after death")
	}
	if player.Rune != 'x' {
		t.Error("Player rune should be 'x' after death")
	}
	if len(player.Items) != 1 {
		t.Error("Player's items should remain in inventory after death")
	}
	if len(level.Items[player.Pos]) != 0 {
		t.Error("Player's items are not dropped on death")
	}
}

func TestCheckDoor(t *testing.T) {
	// Create larger test level (15x15)
	level := &Level{
		Map: make([][]Tile, 15),
		Player: &Player{
			Character: Character{
				Entity: Entity{ // Properly nested Pos
					Pos: Pos{X: 7, Y: 7}, // Now in Entity
				},
				SightRange: 3,
			},
		},
		Monsters: make(map[Pos]*Monster),
		Items:    make(map[Pos][]*Item),
		Portals:  make(map[Pos]*LevelPos),
	}

	// Initialize 15x15 map with dirt floors
	for y := range level.Map {
		level.Map[y] = make([]Tile, 15)
		for x := range level.Map[y] {
			level.Map[y][x] = Tile{Rune: DirtFloor}
		}
	}

	// Place closed door at valid position
	doorPos := Pos{X: 7, Y: 6}
	level.Map[doorPos.Y][doorPos.X] = Tile{
		Rune:        Pending,
		OverlayRune: ClosedDoor,
	}

	checkDoor(level, doorPos)

	// Verify door opened
	if level.Map[doorPos.Y][doorPos.X].OverlayRune != OpenDoor {
		t.Error("checkDoor should open closed door")
	}

	// Verify line of sight updated
	visibleFound := false
	for y := range level.Map {
		for x := range level.Map[y] {
			if level.Map[y][x].Visible {
				visibleFound = true
				break
			}
		}
		if visibleFound {
			break
		}
	}
	if !visibleFound {
		t.Error("checkDoor should update visibility")
	}
}

func TestCheckTrap(t *testing.T) {
	// Create proper test level
	level := &Level{
		Map: make([][]Tile, 15),
		Player: &Player{
			Character: Character{
				Entity: Entity{
					Pos:  Pos{X: 7, Y: 7}, // Center of 15x15 map
					Name: "You",
					Rune: '@',
				},
				Hitpoints:  100, // Explicitly set hitpoints
				SightRange: 3,
			},
		},
		Monsters: make(map[Pos]*Monster),
		Items:    make(map[Pos][]*Item),
		Portals:  make(map[Pos]*LevelPos),
	}

	// Initialize 15x15 map
	for y := range level.Map {
		level.Map[y] = make([]Tile, 15)
		for x := range level.Map[y] {
			level.Map[y][x] = Tile{Rune: DirtFloor}
		}
	}

	// Set trap at player's position
	trapPos := level.Player.Pos
	level.Map[trapPos.Y][trapPos.X] = Tile{
		Rune:        Pending,
		OverlayRune: ClosedTrap,
	}

	checkTrap(level, trapPos)

	// Verify trap state
	if level.Map[trapPos.Y][trapPos.X].OverlayRune != OpenTrap {
		t.Error("checkTrap should open trap")
	}

	// Verify player state
	if level.Player.Hitpoints > 0 {
		t.Error("checkTrap should reduce player hitpoints to 0")
	}
	if level.Player.Rune != 'x' {
		t.Error("checkTrap should change player rune to 'x'")
	}
}

func TestMoveWithPortal(t *testing.T) {
	// Create proper test environment
	game := &Game{
		Levels: make(map[string]*Level),
	}

	// Create source level (15x15)
	srcLevel := &Level{
		Map: make([][]Tile, 15),
		Player: &Player{
			Character: Character{
				Entity: Entity{
					Pos:  Pos{X: 7, Y: 7},
					Name: "You",
					Rune: '@',
				},
				SightRange: 3,
			},
		},
		Portals: make(map[Pos]*LevelPos),
	}

	// Initialize source map
	for y := range srcLevel.Map {
		srcLevel.Map[y] = make([]Tile, 15)
		for x := range srcLevel.Map[y] {
			srcLevel.Map[y][x] = Tile{Rune: DirtFloor}
		}
	}

	// Create destination level (15x15)
	dstLevel := &Level{
		Map: make([][]Tile, 15),
		Player: &Player{
			Character: Character{
				Entity: Entity{
					Name: "You",
					Rune: '@',
				},
			},
		},
	}
	for y := range dstLevel.Map {
		dstLevel.Map[y] = make([]Tile, 15)
		for x := range dstLevel.Map[y] {
			dstLevel.Map[y][x] = Tile{Rune: DirtFloor}
		}
	}

	// Set up portal system
	portalPos := Pos{X: 8, Y: 7} // Valid position in 15x15 grid
	targetPos := Pos{X: 7, Y: 7} // Center of destination level
	srcLevel.Portals[portalPos] = &LevelPos{
		Level: dstLevel,
		Pos:   targetPos,
	}

	game.CurrentLevel = srcLevel
	game.Levels["source"] = srcLevel
	game.Levels["dest"] = dstLevel

	// Perform move
	game.Move(portalPos)

	// Verify transition
	if game.CurrentLevel != dstLevel {
		t.Error("Move should transition to destination level")
	}
	if game.CurrentLevel.Player.Pos != targetPos {
		t.Error("Move should set player to target position")
	}
}

func TestHandleInputMovement(t *testing.T) {
	game := createTestGame()
	initialPos := game.CurrentLevel.Player.Pos
	input := &Input{Typ: Up}
	game.handleInput(input)
	newPos := game.CurrentLevel.Player.Pos
	expectedPos := Pos{initialPos.X, initialPos.Y - 1}
	if newPos != expectedPos {
		t.Errorf("Expected player position %v, got %v", expectedPos, newPos)
	}
}

func TestHandleInputTakeItem(t *testing.T) {
	game := createTestGame()
	level := game.CurrentLevel
	pos := level.Player.Pos
	item := NewSword(pos)
	level.Items[pos] = append(level.Items[pos], item)

	input := &Input{Typ: TakeItem, Item: item}
	game.handleInput(input)

	if len(level.Player.Items) != 1 || level.Player.Items[0] != item {
		t.Error("Item not added to player's inventory")
	}
	if len(level.Items[pos]) != 0 {
		t.Error("Item not removed from level")
	}
}

func TestHandleInputDropItem(t *testing.T) {
	game := createTestGame()
	level := game.CurrentLevel
	player := level.Player
	item := NewSword(player.Pos)
	player.Items = append(player.Items, item)

	input := &Input{Typ: DropItem, Item: item}
	game.handleInput(input)

	if len(player.Items) != 0 {
		t.Error("Item not removed from inventory")
	}
	if len(level.Items[player.Pos]) != 1 || level.Items[player.Pos][0] != item {
		t.Error("Item not added to level")
	}
}

func TestHandleInputCombatBurst(t *testing.T) {
	game := createTestGame()
	level := game.CurrentLevel
	player := level.Player
	monsterPos := Pos{7, 6}
	monster := NewRat(monsterPos)
	level.Monsters[monsterPos] = monster
	level.Attack(&player.Character, &monster.Character)
	player.Burst = &Burst{Notes: []int{2}, MaxCombo: 1, Combo: 0} // Assuming 2 corresponds to Up
	player.Stamina = 1

	input := &Input{Typ: Up}
	game.handleInput(input)

	if len(player.Burst.Notes) != 0 {
		t.Error("Burst note not consumed")
	}
	if player.Stamina != 0 {
		t.Error("Stamina not decremented")
	}
}

func TestHandleInputCloseWindow(t *testing.T) {
	game := createTestGame()
	initialChannels := len(game.LevelChans)
	testChan := game.LevelChans[0]
	input := &Input{Typ: CloseWindow, LevelChannel: testChan}
	game.handleInput(input)

	if len(game.LevelChans) != initialChannels-1 {
		t.Error("Level channel not removed")
	}
}

func TestRunProcessesInput(t *testing.T) {
	game := createTestGame()
	levelChan := make(chan *Level, 1)
	game.LevelChans = []chan *Level{levelChan}

	go game.Run()
	defer close(game.InputChan)

	initialPos := game.CurrentLevel.Player.Pos
	game.InputChan <- &Input{Typ: Up}

	updatedLevel := <-levelChan
	newPos := updatedLevel.Player.Pos
	expectedPos := Pos{initialPos.X, initialPos.Y - 1}
	if newPos != expectedPos {
		t.Errorf("Expected player position %v after Run, got %v", expectedPos, newPos)
	}
}

func TestRunUpdatesMonsters(t *testing.T) {
	game := createTestGame()
	level := game.CurrentLevel
	monsterPos := Pos{7, 6}
	monster := NewRat(monsterPos)
	monster.ActionPoints = 0.0
	level.Monsters[monsterPos] = monster

	levelChan := make(chan *Level, 1)
	game.LevelChans = []chan *Level{levelChan}

	go game.Run()
	defer close(game.InputChan)

	<-levelChan // Discard initial state sent on Run start

	// Send dummy input to trigger game loop processing
	game.InputChan <- &Input{Typ: None}

	updatedLevel := <-levelChan // Get state after update
	updatedMonster, exists := updatedLevel.Monsters[monsterPos]
	if !exists {
		t.Fatal("Monster not found in level")
	}
	// Correct the expected value from 1.5 to 0.5
	if updatedMonster.ActionPoints != 0.5 {
		t.Errorf("Expected monster ActionPoints 0.5, got %f", updatedMonster.ActionPoints)
	}
}
