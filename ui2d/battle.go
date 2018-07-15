package ui2d

import (
	"github.com/maxproske/lyns-rhythm-dungeon/game"
	"github.com/veandco/go-sdl2/sdl"
)

// DrawBurst renders a short pattern of arrows on the battle UI
// TODO(Max): Draw the attacker's burst first. (player -> character)
func (ui *ui) DrawBurst(p *game.Player) {

	colors := [4]int{1, 4, 2, 4}

	// For now, always draw the players's burst first
	offsetX := int32(ui.winWidth/2 - int(float64(ui.winWidth)*0.1)) // Cast int to int32 since we will always use it as int32
	offsetY := int32(ui.winHeight/2 - int(float64(ui.winHeight)*0.1))

	// Play playfield
	playfieldRect := sdl.Rect{offsetX, offsetY, game.NumKeys * 24, int32(len(p.Burst)+1) * 20}
	ui.renderer.Copy(ui.playfieldBackground, nil, &playfieldRect)

	// Draw receptors
	srcRect := ui.noteskinIndex[game.Receptor][0]
	for i := 0; i < game.NumKeys; i++ {
		dstRect := sdl.Rect{int32(i*24) + offsetX, int32(0) + offsetY, 24, 20}
		ui.renderer.Copy(ui.noteskinAtlas, &srcRect, &dstRect)
	}

	// fmt.Println("Burst:", p.Burst)

	for noteIndex, columnIndex := range p.Burst {
		// Get note colour
		noteskinIndex := colors[noteIndex%game.NumKeys]
		noteskinRune := getRuneFromNoteskinIndex(noteskinIndex)
		srcRect := ui.noteskinIndex[noteskinRune][0]
		dstRect := sdl.Rect{int32(columnIndex*24) + offsetX, int32((noteIndex+1)*20) + offsetY, 24, 20}
		ui.renderer.Copy(ui.noteskinAtlas, &srcRect, &dstRect)
	}
	//TODO(Max): Append and take away from the slice using Go magic when a note is hit
}

func getRuneFromNoteskinIndex(i int) rune {
	switch i {
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
	invRect := ui.getBattleRect()
	ui.renderer.Copy(ui.battleBackground, nil, invRect)

	// Draw the player's burst first
	ui.DrawBurst(level.Player)
}

func (ui *ui) getBattleRect() *sdl.Rect {
	invWidth := int32(float32(ui.winWidth) * 0.4)
	invHeight := int32(float32(ui.winHeight) * 0.75)
	offsetX := (int32(ui.winWidth) - invWidth) / 2
	offsetY := (int32(ui.winHeight) - invHeight) / 2
	return &sdl.Rect{offsetX, offsetY, invWidth, invHeight}
}
