package main

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/jenska/gong"
)

func main() {
	ebiten.SetWindowSize(800, 600)
	ebiten.SetWindowTitle("Gong! - The Go Pong")
	if err := ebiten.RunGame(gong.NewGong()); err != nil {
		panic(err)
	}
}
