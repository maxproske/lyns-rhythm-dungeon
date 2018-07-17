package ui2d

import (
	"github.com/maxproske/lyns-rhythm-dungeon/game"
	"github.com/veandco/go-sdl2/sdl"
)

// DrawBurst renders a short pattern of arrows on the battle UI
// TODO(Max): Draw the attacker's burst first. (player -> character)
func (ui *ui) DrawBurst(c, defender *game.Character) {

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

	// Draw versus context
	yOffset := int32(60)
	xCenter := int32(24)
	// Draw attacker
	c1SrcRect := ui.textureIndex[c.Rune][0]
	ui.renderer.Copy(ui.textureAtlas, &c1SrcRect, &sdl.Rect{playfieldRect.X + xCenter, playfieldRect.Y - yOffset, 32, 32})
	// Draw attcker weapon
	if c.Weapon != nil {
		weaponSrcRect := ui.textureIndex[c.Weapon.Rune][0]
		ui.renderer.Copy(ui.textureAtlas, &weaponSrcRect, &sdl.Rect{playfieldRect.X + playfieldRect.W/2 - 32/2, playfieldRect.Y - yOffset, 32, 32})
	}
	// Draw defender
	c2SrcRect := ui.textureIndex[defender.Rune][0]
	ui.renderer.Copy(ui.textureAtlas, &c2SrcRect, &sdl.Rect{playfieldRect.X + playfieldRect.W - 32 - xCenter, playfieldRect.Y - yOffset, 32, 32})
	// Draw hitpoints
	hpXOffset := int32(-16)
	hpYOffset := int32(40)
	hpSrcRect := ui.textureIndex['<'][0]
	attackerHPDstRect := sdl.Rect{playfieldRect.X + xCenter + hpXOffset, playfieldRect.Y - yOffset + hpYOffset, 32, 32}
	defenderHPDstRect := sdl.Rect{playfieldRect.X + playfieldRect.W - 32 - xCenter + hpXOffset, playfieldRect.Y - yOffset + hpYOffset, 32, 32}
	ui.renderer.Copy(ui.textureAtlas, &hpSrcRect, &attackerHPDstRect)
	ui.renderer.Copy(ui.textureAtlas, &hpSrcRect, &defenderHPDstRect)

	//digitSrcRect := ui.textureIndex['0']
	ui.drawHitpoints(c.Hitpoints, &attackerHPDstRect)
	ui.drawHitpoints(defender.Hitpoints, &defenderHPDstRect)

	// Draw receptors
	srcRect := ui.noteskinIndex[game.Receptor][0]
	for i := 0; i < game.NumKeys; i++ {
		dstRect := sdl.Rect{int32(i*24) + offsetX, int32(0) + offsetY, 24, 20}
		if c.Name == "You" {
			ui.noteskinAtlas.SetColorMod(255, 255, 255)
		} else {
			ui.noteskinAtlas.SetColorMod(255, 0, 0)
		}
		ui.renderer.Copy(ui.noteskinAtlas, &srcRect, &dstRect)
	}

	for noteIndex, columnIndex := range c.Burst.Notes {
		// Get note colour
		noteskinIndex := 0 // Uncoloured for monsters
		if c.Name == "You" {
			noteskinIndex = colors[(noteIndex+c.Burst.Combo)%game.NumKeys] // Coloured for player
			if noteIndex < c.Stamina {
				ui.noteskinAtlas.SetColorMod(255, 255, 255)
			} else {
				ui.noteskinAtlas.SetColorMod(64, 64, 64)
			}
		} else {
			if noteIndex < c.Stamina {
				ui.noteskinAtlas.SetColorMod(255, 0, 0)
			} else {
				ui.noteskinAtlas.SetColorMod(64, 0, 0)
			}
		}
		noteskinRune := getRuneFromNoteskinIndex(noteskinIndex)
		srcRect := ui.noteskinIndex[noteskinRune][0]
		dstRect := sdl.Rect{int32(columnIndex*24) + offsetX, int32((noteIndex+1)*20) + offsetY, 24, 20}
		ui.renderer.Copy(ui.noteskinAtlas, &srcRect, &dstRect)
	}
}

func (ui *ui) drawHitpoints(value int, heartRect *sdl.Rect) {
	digits := ui.getSliceFromInt(value)
	for i, digit := range digits {
		digitSrcRect := ui.textureIndex['0'][digit]
		digitXInterval := int32((len(digits) - 1 - i) * 10) // Print digits from left to right
		digitXOffset := int32(23)
		digitYOffset := int32(-12)
		digitDstRect := sdl.Rect{heartRect.X + digitXOffset + digitXInterval, heartRect.Y + digitYOffset, heartRect.W, heartRect.H}
		ui.renderer.Copy(ui.textureAtlas, &digitSrcRect, &digitDstRect)
	}
}

func (ui *ui) getSliceFromInt(value int) []int {
	if value < 10 {
		return []int{value}
	}
	powerOfTen := 0
	if value >= 10000 {
		powerOfTen = 5
	} else if value >= 1000 {
		powerOfTen = 4
	} else if value >= 100 {
		powerOfTen = 3
	} else if value >= 10 {
		powerOfTen = 2
	}
	digits := make([]int, powerOfTen)
	for i := 0; i < powerOfTen; i++ {
		digit := value % 10
		digits[i] = digit
		value /= 10
	}
	return digits
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
	ui.DrawBurst(level.Battle.C1, level.Battle.C2)
}
