package main

import (
	"runtime"

	"github.com/maxproske/lyns-rhythm-dungeon/game"
	"github.com/maxproske/lyns-rhythm-dungeon/ui2d"
)

func main() {
	// Make new game
	numWindows := 1
	game := game.NewGame(numWindows)

	// Make our UIs
	for i := 0; i < numWindows; i++ {
		go func(i int) {
			runtime.LockOSThread() // Goroutines must stay on the same thread for the window to draw and handle input
			ui := ui2d.NewUI(game.InputChan, game.LevelChans[i])
			ui.Run()
		}(i) // Loop will finish quickly, so pass i in
	}

	game.Run()
}
