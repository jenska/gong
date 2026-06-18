package game

import "github.com/hajimehoshi/ebiten/v2"

const (
	ghostShift  = 5
	trailLength = 3
)

type sprite struct {
	x, y       float64
	trailX     [trailLength]float64
	trailY     [trailLength]float64
	trailCount int
	trailNext  int

	image      *ebiten.Image
	ghostImage *ebiten.Image
	visible    bool
}

func (s *sprite) recordPosition() {
	if !s.visible {
		s.trailCount = 0
		s.trailNext = 0
		return
	}
	s.trailX[s.trailNext] = s.x
	s.trailY[s.trailNext] = s.y
	s.trailNext = (s.trailNext + 1) % trailLength
	s.trailCount = min(s.trailCount+1, trailLength)
}

func (s *sprite) draw(screen *ebiten.Image) {
	if !s.visible {
		return
	}

	var op ebiten.DrawImageOptions
	oldest := 0
	if s.trailCount == trailLength {
		oldest = s.trailNext
	}
	for i := range s.trailCount {
		index := (oldest + i) % trailLength
		op.GeoM.Reset()
		op.ColorScale.Reset()
		op.GeoM.Translate(s.trailX[index], s.trailY[index])
		op.ColorScale.ScaleAlpha(float32(i) * 0.1)
		screen.DrawImage(s.image, &op)
	}

	op.GeoM.Reset()
	op.ColorScale.Reset()
	op.GeoM.Translate(s.x+ghostShift, s.y)
	screen.DrawImage(s.ghostImage, &op)
	op.GeoM.Translate(-ghostShift, 0)
	screen.DrawImage(s.image, &op)
}

func (s *sprite) intersects(other *sprite) bool {
	w1, h1 := s.image.Size()
	w2, h2 := other.image.Size()

	if other.x < s.x+float64(w1) && s.x < other.x+float64(w2) && other.y < s.y+float64(h1) {
		return s.y < other.y+float64(h2)
	}
	return false
}
