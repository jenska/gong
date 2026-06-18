package game

import "testing"

func TestHUDOnlyInvalidatesWhenDisplayedStateChanges(t *testing.T) {
	h := hud{}
	g := Gong{state: play, score1: 1, score2: 2}

	h.update(&g)
	if !h.dirty {
		t.Fatal("first update did not invalidate HUD")
	}

	h.dirty = false
	h.update(&g)
	if h.dirty {
		t.Fatal("unchanged game state invalidated HUD")
	}

	g.score1++
	h.update(&g)
	if !h.dirty {
		t.Fatal("score change did not invalidate HUD")
	}
}

func TestSpriteTrailUsesFixedRingBuffer(t *testing.T) {
	s := sprite{visible: true}
	for i := range trailLength + 2 {
		s.x = float64(i)
		s.y = float64(i * 10)
		s.recordPosition()
	}

	if s.trailCount != trailLength {
		t.Fatalf("trailCount = %d, want %d", s.trailCount, trailLength)
	}
	if s.trailNext != 2 {
		t.Fatalf("trailNext = %d, want 2", s.trailNext)
	}

	wantX := [trailLength]float64{3, 4, 2}
	if s.trailX != wantX {
		t.Fatalf("trailX = %v, want %v", s.trailX, wantX)
	}

	s.visible = false
	s.recordPosition()
	if s.trailCount != 0 || s.trailNext != 0 {
		t.Fatalf("hidden sprite retained trail: count=%d next=%d", s.trailCount, s.trailNext)
	}
}

func BenchmarkHUDUpdateUnchanged(b *testing.B) {
	h := hud{}
	g := Gong{state: play, score1: 4, score2: 7, isComputer1: true}
	h.update(&g)

	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		h.update(&g)
	}
}

func BenchmarkSpriteRecordPosition(b *testing.B) {
	s := sprite{visible: true}

	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		s.x++
		s.recordPosition()
	}
}
