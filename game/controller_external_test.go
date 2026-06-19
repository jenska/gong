package game_test

import "github.com/jenska/gong/game"

type stationaryController struct{}

func (stationaryController) Name() string {
	return "STATIONARY"
}

func (stationaryController) Control(game.GameView) game.Control {
	return game.Control{}
}

var _ game.Controller = stationaryController{}
