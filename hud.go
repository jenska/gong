package gong

import (
	"fmt"
	_ "image/png" // for scanline png
	"log"

	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/text"
	"golang.org/x/image/font"
)

const (
	fontSize      = 30
	smallFontSize = fontSize / 2
)

type hud struct {
	arcadeFont      font.Face
	smallArcadeFont font.Face
	image           *ebiten.Image
	message         string
	hints           []string
	scores          string
}

var (
	tt        *truetype.Font
	scanlines *ebiten.Image
)

func init() {
	var err error
	tt, err = truetype.Parse(fonts.ArcadeN_ttf)
	if err != nil {
		log.Fatal(err)
	}
	scanlines, _, err = ebitenutil.NewImageFromFile("scanlines.png", ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}
}

func newHUD() *hud {
	h := &hud{}

	h.arcadeFont = truetype.NewFace(tt, &truetype.Options{
		Size: fontSize, DPI: 72, Hinting: font.HintingFull,
	})
	h.smallArcadeFont = truetype.NewFace(tt, &truetype.Options{
		Size: smallFontSize, DPI: 72, Hinting: font.HintingFull,
	})

	h.image, _ = ebiten.NewImage(windowWidth, windowHeight, ebiten.FilterDefault)
	sw, sh := scanlines.Size()
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(windowWidth/float64(sw), windowHeight/float64(sh))
	h.image.DrawImage(scanlines, op)
	return h
}

func (h *hud) update(g *Gong) {
	h.message = ""
	h.hints = []string{}
	h.scores = ""
	switch g.state {
	case start:
		h.hints = []string{"", "      GONG!      ", "", "H -> HELP         ", "V -> TWO PLAYERS  ", "A -> AI VS PLAYER ", "B -> AI VS AI     "}
	case controls:
		h.hints = []string{"", "PLAYER 1:   ", "W -> UP     ", "S -> DOWN   ", "", "PLAYER 2:   ",
			"ARROW UP    ", "ARROW DOWN  ", "", "ESC -> EXIT "}
		h.message = "Press SPACE for main menu"
	case pause:
		h.hints = []string{"PAUSED", "", "SPACE -> RESUME ", "R     -> RESTART"}
	case play, interrupt:
		h.message = "SPACE -> PAUSE"
		h.scores = fmt.Sprintf("PLAYER 1: %d <-> PLAYER 2: %d", g.score[leftPlayer], g.score[rightPlayer])
	case gameOver:
		winner := "PLAYER 1 WINS"
		if g.score[rightPlayer] > g.score[leftPlayer] {
			winner = "PLAYER 2 WINS"
		}
		h.scores = fmt.Sprintf("PLAYER 1: %d <-> PLAYER 2: %d", g.score[leftPlayer], g.score[rightPlayer])
		h.hints = []string{"", "GAME OVER!", winner, "", "SPACE -> RESTART"}
	}
}

func (h *hud) draw(screen *ebiten.Image) {
	width, height := screen.Size()
	x := (width - len(h.message)*smallFontSize) / 2
	text.Draw(screen, h.message, h.smallArcadeFont, x+10, height-4-2*smallFontSize, ghostColor)
	text.Draw(screen, h.message, h.smallArcadeFont, x, height-4-2*smallFontSize, objectColor)

	for row, hint := range h.hints {
		x = (width - len(hint)*fontSize) / 2
		text.Draw(screen, hint, h.arcadeFont, x+10, (row+4)*fontSize, ghostColor)
		text.Draw(screen, hint, h.arcadeFont, x, (row+4)*fontSize, objectColor)
	}

	x = (width - len(h.scores)*smallFontSize) / 2
	text.Draw(screen, h.scores, h.smallArcadeFont, x+10, 60, ghostColor)
	text.Draw(screen, h.scores, h.smallArcadeFont, x, 60, objectColor)
	screen.DrawImage(h.image, nil)
}
