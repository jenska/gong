package gong

import "github.com/hajimehoshi/ebiten"

type sprite struct {
	x, y    float64
	xbuffer []float64
	ybuffer []float64

	image      *ebiten.Image
	ghostImage *ebiten.Image
	visible    bool
}

func (s *sprite) bufferSize() int {
	if s.xbuffer == nil {
		s.xbuffer = make([]float64, 0)
		s.ybuffer = make([]float64, 0)
	}
	return len(s.xbuffer)
}

func (s *sprite) draw(screen *ebiten.Image) {
	if s.visible {
		var op *ebiten.DrawImageOptions

		for i := 0; i < s.bufferSize(); i++ {
			op = &ebiten.DrawImageOptions{}
			op.GeoM.Translate(s.xbuffer[i], s.ybuffer[i])
			op.ColorM.Scale(1.0, 1.0, 1.0, float64(i)/10.0)
			screen.DrawImage(s.image, op)
		}

		op = &ebiten.DrawImageOptions{}
		op.GeoM.Translate(s.x+5, s.y)
		screen.DrawImage(s.ghostImage, op)
		op.GeoM.Translate(-5, 0)
		screen.DrawImage(s.image, op)

		if s.bufferSize() < 3 {
			s.xbuffer = append(s.xbuffer, s.x)
			s.ybuffer = append(s.ybuffer, s.y)
		} else {
			s.xbuffer = append(s.xbuffer[1:], s.x)
			s.ybuffer = append(s.ybuffer[1:], s.y)
		}
	}
}

func (s *sprite) intersects(other *sprite) bool {
	w1, h1 := s.image.Size()
	w2, h2 := other.image.Size()

	if other.x < s.x+float64(w1) && s.x < other.x+float64(w2) && other.y < s.y+float64(h1) {
		return s.y < other.y+float64(h2)
	}
	return false
}
