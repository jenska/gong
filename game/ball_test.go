package game

import (
	"math"
	"testing"
)

func TestBallServeDirection(t *testing.T) {
	var b ball

	b.serve(LeftSide, -1)
	if b.xVelocity >= 0 || b.yVelocity >= 0 {
		t.Fatalf("left/up serve velocity = (%v, %v)", b.xVelocity, b.yVelocity)
	}

	b.serve(RightSide, 1)
	if b.xVelocity <= 0 || b.yVelocity <= 0 {
		t.Fatalf("right/down serve velocity = (%v, %v)", b.xVelocity, b.yVelocity)
	}
}

func TestBallBounceChangesDirectionAndIncreasesSpeed(t *testing.T) {
	b := ball{
		sprite:    sprite{x: 700, y: 280},
		xVelocity: ballVelocity,
		yVelocity: 2,
	}
	p := paddle{
		sprite: sprite{x: 730, y: 250},
		side:   RightSide,
	}

	b.bounceFrom(&p)

	if b.xVelocity >= 0 {
		t.Fatalf("xVelocity = %v, want leftward", b.xVelocity)
	}
	if got := math.Abs(b.xVelocity); got != ballVelocity+ballAcceleration {
		t.Fatalf("x speed = %v, want %v", got, ballVelocity+ballAcceleration)
	}
	if b.x != p.x-ballDiameter {
		t.Fatalf("x = %v, want %v", b.x, p.x-ballDiameter)
	}
}

func TestBallImpactControlsDeflection(t *testing.T) {
	top := ball{
		sprite:    sprite{y: 200},
		xVelocity: ballVelocity,
	}
	bottom := ball{
		sprite:    sprite{y: 310},
		xVelocity: ballVelocity,
	}
	p := paddle{
		sprite: sprite{x: 730, y: 250},
		side:   RightSide,
	}

	top.bounceFrom(&p)
	bottom.bounceFrom(&p)

	if top.yVelocity >= 0 {
		t.Fatalf("top impact yVelocity = %v, want upward", top.yVelocity)
	}
	if bottom.yVelocity <= 0 {
		t.Fatalf("bottom impact yVelocity = %v, want downward", bottom.yVelocity)
	}
}

func TestMovingPaddleAddsSpin(t *testing.T) {
	stillBall := ball{
		sprite:    sprite{y: 280},
		xVelocity: ballVelocity,
	}
	spinBall := stillBall
	stillPaddle := paddle{
		sprite: sprite{x: 730, y: 250},
		side:   RightSide,
	}
	movingPaddle := stillPaddle
	movingPaddle.yVelocity = 5

	stillBall.bounceFrom(&stillPaddle)
	spinBall.bounceFrom(&movingPaddle)

	if spinBall.yVelocity <= stillBall.yVelocity {
		t.Fatalf("spin yVelocity = %v, want greater than still paddle %v",
			spinBall.yVelocity, stillBall.yVelocity)
	}
}

func TestBallSpeedIsCapped(t *testing.T) {
	b := ball{
		sprite:    sprite{y: 250},
		xVelocity: ballMaxXVelocity,
		yVelocity: ballMaxYVelocity,
	}
	p := paddle{
		sprite:    sprite{x: 730, y: 250},
		side:      RightSide,
		yVelocity: 20,
	}

	b.bounceFrom(&p)

	if math.Abs(b.xVelocity) > ballMaxXVelocity {
		t.Fatalf("xVelocity = %v, exceeds cap %v", b.xVelocity, ballMaxXVelocity)
	}
	if math.Abs(b.yVelocity) > ballMaxYVelocity {
		t.Fatalf("yVelocity = %v, exceeds cap %v", b.yVelocity, ballMaxYVelocity)
	}
}
