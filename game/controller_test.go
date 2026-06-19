package game

import (
	"math/rand/v2"
	"testing"
)

func TestBallApproaching(t *testing.T) {
	tests := []struct {
		name     string
		side     Side
		velocity float64
		want     bool
	}{
		{name: "left approaching", side: LeftSide, velocity: -1, want: true},
		{name: "left receding", side: LeftSide, velocity: 1, want: false},
		{name: "right approaching", side: RightSide, velocity: 1, want: true},
		{name: "right receding", side: RightSide, velocity: -1, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			view := GameView{Side: tt.side, BallVelocityX: tt.velocity}
			if got := ballApproaching(view); got != tt.want {
				t.Fatalf("ballApproaching() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPredictBallY(t *testing.T) {
	tests := []struct {
		name string
		view GameView
		want float64
	}{
		{
			name: "direct path",
			view: testView(400, 100, 5, 5),
			want: 430,
		},
		{
			name: "reflects from bottom",
			view: testView(400, 500, 5, 5),
			want: 330,
		},
		{
			name: "stationary",
			view: testView(400, 123, 0, 0),
			want: 123,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := predictBallY(tt.view); got != tt.want {
				t.Fatalf("predictBallY() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHumanLikeAIAcceleratesInsteadOfSnappingToSpeed(t *testing.T) {
	controller := newHumanLikeAI(testRandom())
	controller.targetY = 500
	controller.cooldown = 10
	controller.tracking = true
	controller.initialized = true
	view := testView(windowWidth/2, 200, ballVelocity, 0)
	view.PaddleY = 100

	control := controller.Control(view)
	p := paddle{yVelocity: 0}
	p.applyControl(control)

	if p.yVelocity != humanAIAcceleration {
		t.Fatalf("yVelocity = %v, want gradual acceleration of %v",
			p.yVelocity, humanAIAcceleration)
	}
}

func TestHumanLikeAIWaitsAtCurrentPositionDuringInitialReaction(t *testing.T) {
	controller := newHumanLikeAI(testRandom())
	view := testView(100, 500, ballVelocity, ballVelocity)
	view.PaddleY = 240

	control := controller.Control(view)

	if control.TargetVelocity != 0 {
		t.Fatalf("TargetVelocity = %v, want 0 during initial reaction",
			control.TargetVelocity)
	}
}

func TestHumanLikeAIBrakesBeforeChangingDirection(t *testing.T) {
	controller := newHumanLikeAI(testRandom())
	controller.targetY = 0
	controller.cooldown = 10
	controller.tracking = true
	controller.initialized = true
	view := testView(windowWidth/2, 200, ballVelocity, 0)
	view.PaddleY = 300
	view.PaddleVelocity = 2

	control := controller.Control(view)
	p := paddle{yVelocity: view.PaddleVelocity}
	p.applyControl(control)

	want := 2 - humanAIBraking
	if p.yVelocity != want {
		t.Fatalf("yVelocity = %v, want braking to %v", p.yVelocity, want)
	}
}

func TestHumanLikeAIReactsMoreSlowlyToDistantBall(t *testing.T) {
	near := newHumanLikeAI(testRandom())
	far := newHumanLikeAI(testRandom())
	nearView := testView(windowWidth-paddleShift-paddleWidth-50, 200, ballVelocity, 0)
	farView := testView(50, 200, ballVelocity, 0)

	nearDelay := near.reactionDelay(nearView)
	farDelay := far.reactionDelay(farView)
	if farDelay <= nearDelay {
		t.Fatalf("far reaction delay = %d, want greater than near delay %d",
			farDelay, nearDelay)
	}
}

func TestHumanLikeAITargetStaysWithinReachableArea(t *testing.T) {
	controller := newHumanLikeAI(testRandom())
	view := testView(100, 10, ballVelocity, -ballVelocity)

	for range 100 {
		target := controller.nextTarget(view)
		if target < 0 || target > windowHeight-paddleHeight {
			t.Fatalf("target = %v, want within [0, %d]",
				target, windowHeight-paddleHeight)
		}
	}
}

func TestAIImplementationsHaveDistinctBehavior(t *testing.T) {
	view := testView(100, 500, ballVelocity, ballVelocity)
	view.PaddleY = 200

	beginner := NewBeginnerAI().Control(view)
	perfect := NewPerfectAI().Control(view)

	if beginner.TargetVelocity == perfect.TargetVelocity {
		t.Fatalf("beginner and perfect target velocity both equal %v",
			beginner.TargetVelocity)
	}
	if perfect.Acceleration <= beginner.Acceleration {
		t.Fatalf("perfect acceleration = %v, want greater than beginner %v",
			perfect.Acceleration, beginner.Acceleration)
	}
}

func TestControllerNames(t *testing.T) {
	tests := []struct {
		controller Controller
		want       string
	}{
		{controller: NewBeginnerAI(), want: "BEGINNER AI"},
		{controller: NewHumanLikeAI(), want: "HUMAN AI"},
		{controller: NewPerfectAI(), want: "PERFECT AI"},
	}
	for _, tt := range tests {
		if got := tt.controller.Name(); got != tt.want {
			t.Fatalf("Name() = %q, want %q", got, tt.want)
		}
	}
}

func testView(ballX, ballY, velocityX, velocityY float64) GameView {
	return GameView{
		Side:          RightSide,
		Width:         windowWidth,
		Height:        windowHeight,
		PaddleX:       windowWidth - paddleShift - paddleWidth,
		PaddleY:       windowHeight/2 - paddleHeight/2,
		PaddleHeight:  paddleHeight,
		BallX:         ballX,
		BallY:         ballY,
		BallVelocityX: velocityX,
		BallVelocityY: velocityY,
	}
}

func testRandom() *rand.Rand {
	return rand.New(rand.NewPCG(1, 2))
}
