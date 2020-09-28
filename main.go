package main

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/jenska/gong/game"
)

func main() {
	if err := ebiten.RunGame(game.NewGong()); err != nil {
		panic(err)
	}
}
