package game

import (
	"embed"
	_ "embed"
	"fmt"
	"image"
	_ "image/png" // load scanline png

	"github.com/golang/freetype/truetype"
	ebiten "github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
)

const (
	fontSize       = 30
	smallFontSize  = fontSize / 2
	bigFontSize    = fontSize * 2
	ghostTextShift = 10
)

type hud struct {
	image   *ebiten.Image
	message string
	hints   []string
	scores  string
	splash  string
}

var (
	//go:embed assets/*
	content embed.FS

	tt          *truetype.Font
	scanlines   *ebiten.Image
	arcadeFonts = make(map[int]font.Face)
)

func assert(err error) {
	if err != nil {
		panic(err)
	}
}

func init() {
	var err error
	binFont, err := content.ReadFile("assets/arcade_n.ttf")
	assert(err)
	tt, err = truetype.Parse(binFont)
	assert(err)
	imageFile, _ := content.Open("assets/scanlines.png")
	image, _, _ := image.Decode(imageFile)
	scanlines = ebiten.NewImageFromImage(image)
	assert(err)

	arcadeFonts[fontSize] = truetype.NewFace(tt, &truetype.Options{
		Size: fontSize, DPI: 72, Hinting: font.HintingFull,
	})
	arcadeFonts[smallFontSize] = truetype.NewFace(tt, &truetype.Options{
		Size: smallFontSize, DPI: 72, Hinting: font.HintingFull,
	})
	arcadeFonts[bigFontSize] = truetype.NewFace(tt, &truetype.Options{
		Size: bigFontSize, DPI: 72, Hinting: font.HintingFull,
	})

}

func newHUD() *hud {
	h := &hud{}
	h.image = ebiten.NewImage(windowWidth, windowHeight)
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
	h.splash = ""

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
			"GONG!",
			"",
			"V   -> TWO PLAYERS ",
			"A   -> AI VS PLAYER",
			"B   -> AI VS AI    ",
			"",
			"F   -> Fullscreen  ",
			fmt.Sprintf("W/S -> Volume %3d%% ", int(sl.volume*100)),
			"H   -> HELP        ",
			"ESC -> Exit        "}
	case controls:
		h.hints = []string{
			"",
			"PLAYER 1:  ",
			"W -> UP    ",
			"S -> DOWN  ",
			"",
			"PLAYER 2:  ",
			"ARROW UP   ",
			"ARROW DOWN "}
		h.message = "SPACE -> main menu"
	case pause:
		h.hints = []string{
			"PAUSED",
			"",
			"SPACE -> RESUME ",
			"R     -> RESTART"}
	case interrupt:
		h.splash = fmt.Sprintf("%d : %d", g.score1, g.score2)
		fallthrough
	case play:
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
	drawText(windowHeight-4-2*smallFontSize, h.message, smallFontSize, screen)
	for row, hint := range h.hints {
		drawText((row+4)*fontSize, hint, fontSize, screen)
	}
	drawText(60, h.scores, smallFontSize, screen)
	drawText((windowHeight+bigFontSize)/2, h.splash, bigFontSize, screen)
	screen.DrawImage(h.image, nil)
}

func drawText(y int, str string, size int, screen *ebiten.Image) {
	x := (windowWidth - len(str)*size) / 2
	text.Draw(screen, str, arcadeFonts[size], x+ghostTextShift, y, ghostColor)
	text.Draw(screen, str, arcadeFonts[size], x, y, objectColor)
}
