package main

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/jenska/gong/game"
)

func main() {
	ebiten.SetWindowSize(800, 600)
	ebiten.SetWindowTitle("Gong! - The Go Pong")
	if err := ebiten.RunGame(game.NewGong()); err != nil {
		panic(err)
	}
}
