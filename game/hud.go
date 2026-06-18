package game

import (
	"bytes"
	"embed"
	"fmt"
	"image"
	"image/color"
	_ "image/png" // load scanline png

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

const (
	fontSize       = 30
	smallFontSize  = fontSize / 2
	bigFontSize    = fontSize * 2
	ghostTextShift = 10
)

type hud struct {
	scanlineImage *ebiten.Image
	textImage     *ebiten.Image
	message       string
	hints         []string
	scores        string
	splash        string
	snapshot      hudSnapshot
	initialized   bool
	dirty         bool
}

type hudSnapshot struct {
	state                    gameState
	score1, score2           int
	isComputer1, isComputer2 bool
	volumePercent            int
}

var (
	//go:embed assets/*
	content embed.FS

	arcadeFontSource *text.GoTextFaceSource
	scanlines        *ebiten.Image
	arcadeFonts      map[int]*text.GoTextFace
)

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func init() {
	binFont, err := content.ReadFile("assets/arcade_n.ttf")
	must(err)
	arcadeFontSource, err = text.NewGoTextFaceSource(bytes.NewReader(binFont))
	must(err)

	binScanlines, err := content.ReadFile("assets/scanlines.png")
	must(err)
	scanlineImage, _, err := image.Decode(bytes.NewReader(binScanlines))
	must(err)
	scanlines = ebiten.NewImageFromImage(scanlineImage)

	arcadeFonts = make(map[int]*text.GoTextFace, 3)
	for _, size := range []int{smallFontSize, fontSize, bigFontSize} {
		arcadeFonts[size] = &text.GoTextFace{
			Source: arcadeFontSource,
			Size:   float64(size),
		}
	}
}

func newHUD() *hud {
	h := &hud{
		scanlineImage: ebiten.NewImage(windowWidth, windowHeight),
		textImage:     ebiten.NewImage(windowWidth, windowHeight),
		dirty:         true,
	}
	sw, sh := scanlines.Size()
	var op ebiten.DrawImageOptions
	op.GeoM.Scale(windowWidth/float64(sw), windowHeight/float64(sh))
	h.scanlineImage.DrawImage(scanlines, &op)
	return h
}

func (h *hud) update(g *Gong) {
	snapshot := hudSnapshot{
		state:         g.state,
		score1:        g.score1,
		score2:        g.score2,
		isComputer1:   g.isComputer1,
		isComputer2:   g.isComputer2,
		volumePercent: int(sl.volume * 100),
	}
	if h.initialized && snapshot == h.snapshot {
		return
	}
	h.snapshot = snapshot
	h.initialized = true
	h.dirty = true

	h.message = ""
	h.hints = nil
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
			fmt.Sprintf("W/S -> Volume %3d%% ", snapshot.volumePercent),
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
	if h.dirty {
		h.textImage.Clear()
		drawText(windowHeight-4-2*smallFontSize, h.message, smallFontSize, h.textImage)
		for row, hint := range h.hints {
			drawText((row+4)*fontSize, hint, fontSize, h.textImage)
		}
		drawText(60, h.scores, smallFontSize, h.textImage)
		drawText((windowHeight+bigFontSize)/2, h.splash, bigFontSize, h.textImage)
		h.dirty = false
	}
	screen.DrawImage(h.textImage, nil)
	screen.DrawImage(h.scanlineImage, nil)
}

func drawText(y int, str string, size int, screen *ebiten.Image) {
	draw := func(x int, color color.Color) {
		var op text.DrawOptions
		op.GeoM.Translate(float64(x), float64(y-size))
		op.PrimaryAlign = text.AlignCenter
		op.ColorScale.ScaleWithColor(color)
		text.Draw(screen, str, arcadeFonts[size], &op)
	}
	draw(windowWidth/2+ghostTextShift, ghostColor)
	draw(windowWidth/2, objectColor)
}
