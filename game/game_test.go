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

	g.interrupt(LeftSide)

	if g.state != interrupt {
		t.Fatalf("state = %v, want interrupt", g.state)
	}
	if g.interruptTicks != interruptDurationTicks {
		t.Fatalf("interruptTicks = %d, want %d", g.interruptTicks, interruptDurationTicks)
	}
	if g.ball.x != windowWidth/2 || g.ball.y != windowHeight/2 {
		t.Fatalf("ball position = (%v, %v), want center", g.ball.x, g.ball.y)
	}
	if g.ball.xVelocity != -ballVelocity || g.ball.yVelocity == 0 {
		t.Fatalf("ball velocity = (%v, %v), want leftward serve",
			g.ball.xVelocity, g.ball.yVelocity)
	}
}

func TestNewAIControllerUsesSelectedLevel(t *testing.T) {
	tests := []struct {
		level aiLevel
		want  string
	}{
		{level: beginnerLevel, want: "BEGINNER AI"},
		{level: humanLikeLevel, want: "HUMAN AI"},
		{level: perfectLevel, want: "PERFECT AI"},
	}
	for _, tt := range tests {
		g := Gong{aiLevel: tt.level}
		if got := g.newAIController().Name(); got != tt.want {
			t.Fatalf("level %v created %q, want %q", tt.level, got, tt.want)
		}
	}
}

func TestStartMatchInstallsCustomControllers(t *testing.T) {
	left := NewBeginnerAI()
	right := NewPerfectAI()
	g := Gong{}
	g.reset()

	g.StartMatch(left, right)

	if g.leftPaddle.controller != left {
		t.Fatal("left controller was not installed")
	}
	if g.rightPaddle.controller != right {
		t.Fatal("right controller was not installed")
	}
	if g.state != play {
		t.Fatalf("state = %v, want play", g.state)
	}
}
