package game

import (
	"bytes"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
)

const (
	sampleRate = 44100

	ping sound = iota
	pong
	lost
	win
	menuSelect
)

type (
	sound        int
	soundLibrary struct {
		audioContext *audio.Context
		volume       float64
		players      map[sound]*audio.Player
	}
)

var sl *soundLibrary

func init() {
	sl = &soundLibrary{
		audioContext: audio.NewContext(sampleRate),
		volume:       1,
		players:      make(map[sound]*audio.Player, 5),
	}

	newPlayer := func(fileName string) *audio.Player {
		buffer, err := content.ReadFile(fileName)
		if err != nil {
			panic(err)
		}
		stream, err := wav.Decode(sl.audioContext, bytes.NewReader(buffer))
		if err != nil {
			panic(err)
		}
		player, err := audio.NewPlayer(sl.audioContext, stream)
		if err != nil {
			panic(err)
		}
		return player
	}
	sl.players[ping] = newPlayer("assets/ping.wav")
	sl.players[pong] = newPlayer("assets/pong.wav")
	sl.players[lost] = newPlayer("assets/lost.wav")
	sl.players[win] = newPlayer("assets/win.wav")
	sl.players[menuSelect] = newPlayer("assets/menu_select.wav")
}

func playSound(s sound) {
	if audioPlayer, ok := sl.players[s]; ok {
		audioPlayer.SetVolume(sl.volume)
		if err := audioPlayer.Rewind(); err != nil {
			panic(err)
		}
		audioPlayer.Play()
	}
}
