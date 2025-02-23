package game

import "testing"

func TestNewSword(t *testing.T) {
	pos := Pos{X: 1, Y: 1}
	sword := NewSword(pos)
	if sword.Typ != Weapon {
		t.Error("Sword should be of type Weapon")
	}
	if sword.Name != "Sword" {
		t.Error("Incorrect sword name")
	}
	if sword.power <= 0 {
		t.Error("Sword should have positive power")
	}
}

func TestNewHelmet(t *testing.T) {
	pos := Pos{X: 2, Y: 2}
	helmet := NewHelmet(pos)
	if helmet.Typ != Helmet {
		t.Error("Helmet should be of type Helmet")
	}
	if helmet.Name != "Helmet" {
		t.Error("Incorrect helmet name")
	}
}

func TestNewCredits(t *testing.T) {
	pos := Pos{X: 3, Y: 3}
	credits := NewCredits(pos)
	if credits.Typ != Other {
		t.Error("Credits should be of type Other")
	}
	if credits.Name != "Credits" {
		t.Error("Incorrect credits name")
	}
}

func TestNewPotion(t *testing.T) {
	pos := Pos{X: 4, Y: 4}
	potion := NewPotion(pos)
	if potion.Typ != Other {
		t.Error("Potion should be of type Other")
	}
	if potion.Name != "Health Potion" {
		t.Error("Incorrect potion name")
	}
	if potion.power != 16.0 {
		t.Error("Incorrect potion power")
	}
}

func TestEquipmentEffects(t *testing.T) {
	pos := Pos{X: 1, Y: 1}

	// Test sword power
	sword := NewSword(pos)
	if sword.power <= 0 {
		t.Error("Sword should have positive power value")
	}

	// Test multiple weapons
	strongSword := &Item{
		Entity: Entity{Name: "Strong Sword", Pos: pos},
		Typ:    Weapon,
		power:  2.0,
	}

	char := &Character{
		Entity:    Entity{Name: "Test Character"},
		Hitpoints: 100,
		Items:     []*Item{}, // Initialize items slice
	}

	// Add items to inventory before equipping
	char.Items = append(char.Items, sword)

	// Test weapon switching
	equip(char, sword)
	if char.Weapon != sword {
		t.Error("Failed to equip first sword")
	}

	// Add strong sword to inventory before equipping
	char.Items = append(char.Items, strongSword)
	equip(char, strongSword)
	if char.Weapon != strongSword {
		t.Error("Failed to switch to stronger sword")
	}
	if char.Weapon == sword {
		t.Error("Old sword should be unequipped")
	}
}

func TestArmorMechanics(t *testing.T) {
	pos := Pos{X: 1, Y: 1}
	helmet := NewHelmet(pos)

	// Test basic armor values
	if helmet.power <= 0 {
		t.Error("Helmet should provide positive power value")
	}

	// Test armor equipping
	char := &Character{
		Entity:    Entity{Name: "Test Character"},
		Hitpoints: 100,
		Items:     []*Item{}, // Initialize items slice
	}

	// Add helmet to inventory before equipping
	char.Items = append(char.Items, helmet)
	equip(char, helmet)
	if char.Helmet != helmet {
		t.Error("Failed to equip helmet")
	}

	strongHelmet := &Item{
		Entity: Entity{Name: "Strong Helmet", Pos: pos},
		Typ:    Helmet,
		power:  helmet.power * 2,
	}

	// Add strong helmet to inventory before equipping
	char.Items = append(char.Items, strongHelmet)
	equip(char, strongHelmet)
	if char.Helmet != strongHelmet {
		t.Error("Failed to equip stronger helmet")
	}
}

func TestItemInteractions(t *testing.T) {
	pos := Pos{X: 1, Y: 1}
	potion := NewPotion(pos)
	credits := NewCredits(pos)

	// Test item stacking
	items := []*Item{potion, credits}
	level := createTestLevel()
	level.Items[pos] = items

	// Test item pickup order
	if level.Items[pos][0] != potion {
		t.Error("Items should maintain order in stack")
	}

	// Test item removal
	level.Items[pos] = level.Items[pos][1:]
	if len(level.Items[pos]) != 1 || level.Items[pos][0] != credits {
		t.Error("Item removal should preserve remaining items")
	}

	// Test empty stack cleanup
	level.Items[pos] = level.Items[pos][:0]
	if len(level.Items[pos]) != 0 {
		t.Error("Empty item stack should be cleared")
	}
}
