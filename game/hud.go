package game

import (
	"fmt"
	_ "image/png" // load scanline png
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
	image   *ebiten.Image
	message string
	hints   []string
	scores  string
}

var (
	tt              *truetype.Font
	scanlines       *ebiten.Image
	arcadeFont      font.Face
	smallArcadeFont font.Face
)

func init() {
	var err error
	tt, err = truetype.Parse(fonts.ArcadeN_ttf)
	if err != nil {
		log.Fatal(err)
	}
	scanlines, _, err = ebitenutil.NewImageFromFile("game/scanlines.png", ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}
	arcadeFont = truetype.NewFace(tt, &truetype.Options{
		Size: fontSize, DPI: 72, Hinting: font.HintingFull,
	})
	smallArcadeFont = truetype.NewFace(tt, &truetype.Options{
		Size: smallFontSize, DPI: 72, Hinting: font.HintingFull,
	})
}

func newHUD() *hud {
	h := &hud{}
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

	player1, player2 := "PLAYER 1", "PLAYER 2"
	if g.isComputer1 {
		player2, player1 = "PLAYER", "COMPUTER"
	}
	if g.isComputer2 {
		player2 = "COMPUTER"
	}

	switch g.state {
	case start:
		h.hints = []string{
			"",
			"      GONG!      ",
			"",
			"H -> HELP         ",
			"V -> TWO PLAYERS  ",
			"A -> AI VS PLAYER ",
			"B -> AI VS AI     "}
	case controls:
		h.hints = []string{
			"",
			"PLAYER 1:  ",
			"W -> UP    ",
			"S -> DOWN  ",
			"",
			"PLAYER 2:  ",
			"ARROW UP   ",
			"ARROW DOWN ",
			"",
			"ESC -> EXIT"}
		h.message = "SPACE -> main menu"
	case pause:
		h.hints = []string{
			"PAUSED",
			"",
			"SPACE -> RESUME ",
			"R     -> RESTART"}
	case play, interrupt:
		h.message = "SPACE -> PAUSE"
		h.scores = fmt.Sprintf("%s: %d <-> %s: %d", player1, g.score1, player2, g.score2)

	case gameOver:
		winner := player1
		if g.score2 > g.score1 {
			winner = player2
		}
		h.scores = fmt.Sprintf("%s: %d <-> %s: %d", player1, g.score1, player2, g.score2)
		h.hints = []string{"", "GAME OVER!", winner + " WINS", "", "SPACE -> RESTART"}
	}
}

func (h *hud) draw(screen *ebiten.Image) {
	drawText(windowHeight-4-2*smallFontSize, h.message, smallArcadeFont, smallFontSize, screen)
	for row, hint := range h.hints {
		drawText((row+4)*fontSize, hint, arcadeFont, fontSize, screen)
	}
	drawText(60, h.scores, smallArcadeFont, smallFontSize, screen)
	screen.DrawImage(h.image, nil)
}

func drawText(y int, str string, font font.Face, size int, screen *ebiten.Image) {
	x := (windowWidth - len(str)*size) / 2
	text.Draw(screen, str, font, x+10, y, ghostColor)
	text.Draw(screen, str, font, x, y, objectColor)
}
