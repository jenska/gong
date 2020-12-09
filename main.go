package main

import (
	ebiten "github.com/hajimehoshi/ebiten/v2"
	"github.com/jenska/gong/game"
)

func main() {
	if err := ebiten.RunGame(game.NewGong()); err != nil {
		panic(err)
	}
}
