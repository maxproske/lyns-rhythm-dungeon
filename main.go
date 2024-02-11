package main

import (
	"runtime"

	"github.com/maxproske/lyns-rhythm-dungeon/game"
	"github.com/maxproske/lyns-rhythm-dungeon/ui2d"
)

func main() {
	// Make new game
	game := game.NewGame(1)
	go game.Run()

	// Make our UI
	runtime.LockOSThread()
	ui := ui2d.NewUI(game.InputChan, game.LevelChans[0])
	ui.Run()
}
