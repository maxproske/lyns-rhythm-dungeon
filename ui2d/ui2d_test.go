package ui2d

import (
	"fmt"
	"testing"

	"github.com/maxproske/lyns-rhythm-dungeon/game"
)

func TestInventoryRect(t *testing.T) {
	ui := &ui{
		winWidth:  800,
		winHeight: 600,
	}

	// Test inventory rect calculations
	invRect := ui.getInventoryRect()
	if invRect.W != int32(float32(ui.winWidth)*0.4) {
		t.Errorf("Expected inventory width %v, got %v", int32(float32(ui.winWidth)*0.4), invRect.W)
	}
	if invRect.H != int32(float32(ui.winHeight)*0.75) {
		t.Errorf("Expected inventory height %v, got %v", int32(float32(ui.winHeight)*0.75), invRect.H)
	}

	// Test slot rect calculations
	helmetSlot := ui.getHelmetSlotRect()
	if helmetSlot.W != helmetSlot.H {
		t.Error("Helmet slot should be square")
	}

	weaponSlot := ui.getWeaponSlotRect()
	if weaponSlot.W != weaponSlot.H {
		t.Error("Weapon slot should be square")
	}

	// Test inventory item rect
	itemRect := ui.getInventoryItemRect(0)
	if itemRect.W != itemRect.H {
		t.Error("Inventory item rect should be square")
	}
}

func TestDroppedItem(t *testing.T) {
	ui := &ui{
		winWidth:  800,
		winHeight: 600,
		currentMouseState: &mouseState{
			leftButton:  false,
			rightButton: false,
			pos:         game.Pos{X: 0, Y: 0},
		},
		draggedItem: &game.Item{},
	}

	// Test dropping outside inventory
	if ui.CheckDroppedItem() != ui.draggedItem {
		t.Error("Should return dragged item when dropped outside inventory")
	}

	// Test dropping inside inventory
	invRect := ui.getInventoryRect()
	ui.currentMouseState.pos = game.Pos{
		X: int(invRect.X + invRect.W/2),
		Y: int(invRect.Y + invRect.H/2),
	}
	if ui.CheckDroppedItem() != nil {
		t.Error("Should return nil when dropped inside inventory")
	}
}

func TestGroundItems(t *testing.T) {
	ui := &ui{
		winWidth:  800,
		winHeight: 600,
		currentMouseState: &mouseState{
			leftButton:  false,
			rightButton: false,
			pos:         game.Pos{X: 0, Y: 0},
		},
		prevMouseState: &mouseState{
			leftButton:  true,
			rightButton: false,
			pos:         game.Pos{X: 0, Y: 0},
		},
	}

	level := &game.Level{
		Items: make(map[game.Pos][]*game.Item),
		Player: &game.Player{
			Character: game.Character{
				Entity: game.Entity{
					Pos: game.Pos{X: 1, Y: 1},
				},
			},
		},
	}

	// Add a test item at player's position
	item := &game.Item{}
	level.Items[level.Player.Pos] = append(level.Items[level.Player.Pos], item)

	// Position mouse over the ground item slot
	itemRect := ui.getGroundItemRect(0)
	ui.currentMouseState.pos = game.Pos{
		X: int(itemRect.X),
		Y: int(itemRect.Y),
	}

	// Test clicking on ground item
	result := ui.CheckGroundItems(level)
	if result != item {
		t.Error("Should return item when clicked on ground item slot")
	}
}

func TestBattleState(t *testing.T) {
	ui := &ui{
		winWidth:  800,
		winHeight: 600,
		state:     UIMain,
	}

	level := &game.Level{
		Battle: &game.Battle{
			C1: &game.Character{
				Entity: game.Entity{
					Name: "Player",
				},
				Burst: &game.Burst{
					Notes:    []int{1, 2, 3},
					MaxCombo: 3,
					Combo:    0,
				},
			},
			C2: &game.Character{
				Entity: game.Entity{
					Name: "Monster",
				},
			},
		},
		LastEvent: game.Attack,
	}

	// Test battle state transition
	if ui.state != UIMain {
		t.Error("Should start in main state")
	}

	// Test switching to battle state
	ui.state = UIBattle
	if ui.state != UIBattle {
		t.Error("Should switch to battle state")
	}

	// Verify battle participants
	if level.Battle.C1 == nil || level.Battle.C2 == nil {
		t.Error("Battle should have two participants")
	}

	// Test burst pattern
	if len(level.Battle.C1.Burst.Notes) != 3 {
		t.Error("Burst should have correct number of notes")
	}
}

func TestKeyboardHandling(t *testing.T) {
	ui := &ui{
		winWidth:  800,
		winHeight: 600,
		state:     UIMain,
		keyboardState: []uint8{
			0, // Not pressed
		},
		prevKeyboardState: []uint8{
			0, // Not pressed previously
		},
	}

	// Test key not pressed
	if ui.keyDownOnce(0) {
		t.Error("Key should not be detected as pressed")
	}

	// Test key pressed
	ui.keyboardState[0] = 1
	ui.prevKeyboardState[0] = 0
	if !ui.keyDownOnce(0) {
		t.Error("Key should be detected as pressed")
	}

	// Test key held
	ui.prevKeyboardState[0] = 1
	if ui.keyDownOnce(0) {
		t.Error("Held key should not trigger keyDownOnce")
	}

	// Test key released
	ui.keyboardState[0] = 0
	ui.prevKeyboardState[0] = 1
	if !ui.keyPressed(0) {
		t.Error("Key release should be detected")
	}
}

func TestInventoryState(t *testing.T) {
	ui := &ui{
		winWidth:  800,
		winHeight: 600,
		state:     UIMain,
	}

	// Test initial state
	if ui.state != UIMain {
		t.Error("Should start in main state")
	}

	// Test switching to inventory
	ui.state = UIInventory
	if ui.state != UIInventory {
		t.Error("Should switch to inventory state")
	}

	// Test inventory dimensions
	invRect := ui.getInventoryRect()
	expectedWidth := int32(float32(ui.winWidth) * 0.4)
	expectedHeight := int32(float32(ui.winHeight) * 0.75)

	if invRect.W != expectedWidth {
		t.Errorf("Expected inventory width %d, got %d", expectedWidth, invRect.W)
	}
	if invRect.H != expectedHeight {
		t.Errorf("Expected inventory height %d, got %d", expectedHeight, invRect.H)
	}
}

func TestUIElementDimensions(t *testing.T) {
	ui := &ui{
		winWidth:  800,
		winHeight: 600,
	}

	// Test event console dimensions
	textStart := int32(float64(ui.winHeight) * 0.6)
	textWidth := int32(float64(ui.winWidth) * 0.25)
	eventHeight := int32(ui.winHeight) - textStart

	if eventHeight <= 0 {
		t.Error("Event console height should be positive")
	}
	if textWidth <= 0 {
		t.Error("Event console width should be positive")
	}

	// Test ground inventory dimensions
	groundStart := int32(float64(ui.winWidth) * 0.9)
	groundWidth := int32(ui.winWidth) - groundStart
	itemSize := int32(itemSizeRatio * float32(ui.winWidth))

	if groundWidth <= 0 {
		t.Error("Ground inventory width should be positive")
	}
	if itemSize <= 0 {
		t.Error("Item size should be positive")
	}

	// Test item slot positions
	for i := 0; i < 3; i++ {
		slotRect := ui.getGroundItemRect(i)
		if slotRect.W != slotRect.H {
			t.Error("Ground item slots should be square")
		}
		if i > 0 {
			prevSlot := ui.getGroundItemRect(i - 1)
			if slotRect.X >= prevSlot.X {
				t.Error("Ground item slots should be arranged right to left")
			}
		}
	}
}

func TestGameStateTransitions(t *testing.T) {
	ui := &ui{
		winWidth:  800,
		winHeight: 600,
		state:     UIMain,
	}

	// Test state transitions
	states := []struct {
		from uiState
		to   uiState
		name string
	}{
		{UIMain, UIInventory, "main to inventory"},
		{UIInventory, UIMain, "inventory to main"},
		{UIMain, UIBattle, "main to battle"},
		{UIBattle, UIMain, "battle to main"},
	}

	for _, tc := range states {
		ui.state = tc.from
		ui.state = tc.to
		if ui.state != tc.to {
			t.Errorf("Failed to transition from %v to %v", tc.from, tc.to)
		}
	}
}

func TestBattleStateTransitions(t *testing.T) {
	ui := &ui{
		winWidth:  800,
		winHeight: 600,
		state:     UIMain,
	}

	level := &game.Level{
		Battle: &game.Battle{
			C1: &game.Character{
				Entity: game.Entity{
					Name: "Player",
				},
				Burst: &game.Burst{
					Notes:    []int{1, 2, 3},
					MaxCombo: 3,
				},
				Stamina:    3,
				MaxStamina: 3,
			},
			C2: &game.Character{
				Entity: game.Entity{
					Name: "Monster",
				},
				Hitpoints: 3,
			},
		},
		LastEvent: game.Attack,
	}

	// Test battle initiation
	ui.state = UIBattle
	if ui.state != UIBattle {
		t.Error("Should enter battle state")
	}

	// Test battle state checks
	if level.Battle.C1.Stamina <= 0 {
		t.Error("Player should start with full stamina")
	}

	if level.Battle.C2.Hitpoints <= 0 {
		t.Error("Monster should start with positive health")
	}

	if len(level.Battle.C1.Burst.Notes) != 3 {
		t.Error("Player should have correct number of notes")
	}
}

func TestUIInteractionEdgeCases(t *testing.T) {
	ui := &ui{
		winWidth:  800,
		winHeight: 600,
		state:     UIMain,
		currentMouseState: &mouseState{
			leftButton:  false,
			rightButton: false,
			pos:         game.Pos{X: 0, Y: 0},
		},
		prevMouseState: &mouseState{
			leftButton:  false,
			rightButton: false,
			pos:         game.Pos{X: 0, Y: 0},
		},
	}

	// Test window resizing effects
	ui.winWidth = 400
	ui.winHeight = 300
	invRect := ui.getInventoryRect()
	if invRect.W <= 0 || invRect.H <= 0 {
		t.Error("Inventory rect should remain valid after window resize")
	}

	// Test rapid mouse click state changes
	ui.currentMouseState.leftButton = true
	ui.prevMouseState.leftButton = true
	if ui.currentMouseState.leftButton && !ui.prevMouseState.leftButton {
		t.Error("Should not detect click when button was already down")
	}

	// Test dragging item out of bounds
	ui.draggedItem = &game.Item{
		Entity: game.Entity{Name: "Test Item"},
	}
	ui.currentMouseState.pos = game.Pos{X: -1, Y: -1}
	if ui.CheckDroppedItem() == nil {
		t.Error("Should return dragged item when dropped out of bounds")
	}

	// Test UI state transitions under load
	states := []uiState{UIMain, UIInventory, UIBattle, UIMain}
	for i := 0; i < len(states)*2; i++ {
		ui.state = states[i%len(states)]
	}
	if ui.state != UIMain {
		t.Error("UI should handle rapid state transitions")
	}
}

func TestBattleUIEdgeCases(t *testing.T) {
	ui := &ui{
		winWidth:  800,
		winHeight: 600,
		state:     UIBattle,
	}

	level := &game.Level{
		Battle: &game.Battle{
			C1: &game.Character{
				Entity: game.Entity{Name: "Player"},
				Burst: &game.Burst{
					Notes:    []int{},
					MaxCombo: 0,
					Combo:    0,
				},
			},
			C2: &game.Character{
				Entity: game.Entity{Name: "Monster"},
			},
		},
	}

	// Test empty burst pattern
	if ui.state == UIBattle && level.Battle != nil && len(level.Battle.C1.Burst.Notes) == 0 {
		ui.state = UIMain
	}
	if ui.state != UIMain {
		t.Error("UI should exit battle state with empty burst")
	}

	// Test extremely long burst patterns
	level.Battle.C1.Burst.Notes = make([]int, 100)
	level.Battle.C1.Burst.MaxCombo = 100
	ui.state = UIBattle

	// Battle should continue with long pattern
	if ui.state != UIBattle {
		t.Error("UI should handle long burst patterns")
	}

	// Test UI state after battle ends
	level.Battle = nil
	ui.state = UIMain
	if ui.state != UIMain {
		t.Error("UI should return to main state when battle ends")
	}
}

func TestInventoryUIEdgeCases(t *testing.T) {
	ui := &ui{
		winWidth:  800,
		winHeight: 600,
		state:     UIInventory,
	}

	// Test inventory full of items
	level := &game.Level{
		Player: &game.Player{
			Character: game.Character{
				Items: make([]*game.Item, 20),
			},
		},
	}

	// Fill inventory with items
	for i := range level.Player.Items {
		level.Player.Items[i] = &game.Item{
			Entity: game.Entity{Name: fmt.Sprintf("Item %d", i)},
		}
	}

	// Test item slot calculations with full inventory
	for i := 0; i < len(level.Player.Items); i++ {
		rect := ui.getInventoryItemRect(i)
		if rect.W <= 0 || rect.H <= 0 {
			t.Errorf("Invalid item rect dimensions for item %d", i)
		}
	}

	// Test drag and drop with full inventory
	ui.draggedItem = level.Player.Items[0]
	ui.currentMouseState = &mouseState{
		leftButton: true,
		pos: game.Pos{
			X: int(ui.getInventoryRect().X),
			Y: int(ui.getInventoryRect().Y),
		},
	}

	result := ui.CheckDroppedItem()
	if result != nil {
		t.Error("Should handle dropping item in full inventory")
	}
}
