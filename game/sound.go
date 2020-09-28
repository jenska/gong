package game

import (
	"io/ioutil"

	"github.com/hajimehoshi/ebiten/audio"
	"github.com/hajimehoshi/ebiten/audio/wav"
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
	sl = &soundLibrary{}

	var err error
	sl.audioContext, err = audio.NewContext(sampleRate)
	if err != nil {
		panic(err)
	}

	newPlayer := func(fileName string) *audio.Player {
		if buffer, err := ioutil.ReadFile(fileName); err != nil {
			panic(err)
		} else if stream, err := wav.Decode(sl.audioContext, audio.BytesReadSeekCloser(buffer)); err != nil {
			panic(err)
		} else if player, err := audio.NewPlayer(sl.audioContext, stream); err != nil {
			panic(err)
		} else {
			return player
		}
	}
	sl.volume = 1.0
	sl.players = make(map[sound]*audio.Player)
	sl.players[ping] = newPlayer("assets/ping.wav")
	sl.players[pong] = newPlayer("assets/pong.wav")
	sl.players[lost] = newPlayer("assets/lost.wav")
	sl.players[win] = newPlayer("assets/win.wav")
	sl.players[menuSelect] = newPlayer("assets/menu_select.wav")
}

func playSound(s sound) {
	if audioPlayer, ok := sl.players[s]; ok {
		audioPlayer.SetVolume(sl.volume)
		audioPlayer.Rewind()
		audioPlayer.Play()
	}
}
