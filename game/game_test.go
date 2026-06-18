package game

import "testing"

func TestInterruptResetsBallWithoutBlocking(t *testing.T) {
	g := Gong{
		state: play,
		ball: &ball{
			sprite:    sprite{x: 10, y: 20},
			xVelocity: -3,
			yVelocity: 4,
		},
	}

	g.interrupt()

	if g.state != interrupt {
		t.Fatalf("state = %v, want interrupt", g.state)
	}
	if g.interruptTicks != interruptDurationTicks {
		t.Fatalf("interruptTicks = %d, want %d", g.interruptTicks, interruptDurationTicks)
	}
	if g.ball.x != windowWidth/2 || g.ball.y != windowHeight/2 {
		t.Fatalf("ball position = (%v, %v), want center", g.ball.x, g.ball.y)
	}
	if g.ball.xVelocity != ballVelocity || g.ball.yVelocity != ballVelocity {
		t.Fatalf("ball velocity = (%v, %v), want (%v, %v)",
			g.ball.xVelocity, g.ball.yVelocity, ballVelocity, ballVelocity)
	}
}
