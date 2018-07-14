package main

import (
	"github.com/maxproske/lyns-rhythm-dungeon/game"
	"github.com/maxproske/lyns-rhythm-dungeon/ui2d"
)

func main() {
	// Make new game
	game := game.NewGame(1)
	go func() {
		game.Run()
	}()

	// Make our UI
	ui := ui2d.NewUI(game.InputChan, game.LevelChans[0])
	ui.Run()
}
