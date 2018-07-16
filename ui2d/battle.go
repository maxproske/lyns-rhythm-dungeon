package ui2d

import (
	"github.com/maxproske/lyns-rhythm-dungeon/game"
	"github.com/veandco/go-sdl2/sdl"
)

// DrawBurst renders a short pattern of arrows on the battle UI
// TODO(Max): Draw the attacker's burst first. (player -> character)
func (ui *ui) DrawBurst(c *game.Character) {

	// Dim the lights
	ui.renderer.Copy(ui.dimOverlay, nil, nil) // Stretch to fit

	colors := [4]int{1, 4, 2, 4}

	// For now, always draw the players's burst first
	offsetX := int32(ui.winWidth/2 - (game.NumKeys*20)/2) // Cast int to int32 since we will always use it as int32
	offsetY := int32(ui.winHeight/2 - (game.NumKeys*24)/2)

	// Draw black playfield with white border
	padding := int32(24)
	borderWidth := int32(1)
	playfieldRect := sdl.Rect{offsetX - padding, offsetY - padding/2, game.NumKeys*24 + padding*2, int32(c.Burst.MaxCombo+1)*20 + padding*3}
	battleBackgroundRect := sdl.Rect{playfieldRect.X - borderWidth, playfieldRect.Y - borderWidth, playfieldRect.W + borderWidth*2, playfieldRect.H + borderWidth*2}
	if c.Name == "You" {
		ui.renderer.Copy(ui.battleBorderPlayer, nil, &battleBackgroundRect)
	} else {
		ui.renderer.Copy(ui.battleBorderMonster, nil, &battleBackgroundRect)

	}

	ui.renderer.Copy(ui.playfieldBackground, nil, &playfieldRect)

	// Draw receptors
	srcRect := ui.noteskinIndex[game.Receptor][0]
	for i := 0; i < game.NumKeys; i++ {
		dstRect := sdl.Rect{int32(i*24) + offsetX, int32(0) + offsetY, 24, 20}
		ui.renderer.Copy(ui.noteskinAtlas, &srcRect, &dstRect)
	}

	for noteIndex, columnIndex := range c.Burst.Notes {
		// Get note colour
		noteskinIndex := 0 // Uncoloured for monsters
		if c.Name == "You" {
			noteskinIndex = colors[(noteIndex+c.Burst.Combo)%game.NumKeys] // Coloured for player
			ui.noteskinAtlas.SetColorMod(255, 255, 255)
		} else {
			ui.noteskinAtlas.SetColorMod(255, 0, 0)
		}
		noteskinRune := getRuneFromNoteskinIndex(noteskinIndex)
		srcRect := ui.noteskinIndex[noteskinRune][0]
		dstRect := sdl.Rect{int32(columnIndex*24) + offsetX, int32((noteIndex+1)*20) + offsetY, 24, 20}
		ui.renderer.Copy(ui.noteskinAtlas, &srcRect, &dstRect)
	}
}

func getRuneFromNoteskinIndex(i int) rune {
	switch i {
	case 0:
		return game.Receptor
	case 1:
		return game.Red
	case 2:
		return game.Blue
	case 4:
		return game.Yellow
	default:
		panic("Noteskin index out of bounds.")
	}
}

// DrawBattle renders the battle screen
func (ui *ui) DrawBattle(level *game.Level) {
	// Draw the attacker's burst first
	ui.DrawBurst(level.Battle.C1)
}
