package game

import (
	"math"
	"math/rand/v2"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// Side identifies a player's side of the table.
type Side byte

const (
	LeftSide Side = iota
	RightSide
)

// GameView is the read-only game state supplied to a Controller each tick.
// Coordinates use the game's logical 800x600 pixel space.
type GameView struct {
	Side           Side
	Width, Height  float64
	PaddleX        float64
	PaddleY        float64
	PaddleHeight   float64
	PaddleVelocity float64
	BallX          float64
	BallY          float64
	BallVelocityX  float64
	BallVelocityY  float64
}

// Control describes how a controller wants its paddle to move this tick.
// Velocities are measured in logical pixels per tick. Acceleration and Braking
// must be non-negative.
type Control struct {
	TargetVelocity float64
	Acceleration   float64
	Braking        float64
}

// Controller controls one paddle. Implementations may retain state between
// calls, but must not mutate the supplied GameView.
type Controller interface {
	Name() string
	Control(GameView) Control
}

type keyboardController struct {
	upKey, downKey ebiten.Key
}

// NewKeyboardController creates a controller using the supplied movement keys.
func NewKeyboardController(upKey, downKey ebiten.Key) Controller {
	return &keyboardController{upKey: upKey, downKey: downKey}
}

func (*keyboardController) Name() string {
	return "PLAYER"
}

func (c *keyboardController) Control(_ GameView) Control {
	var direction float64
	if inpututil.KeyPressDuration(c.upKey) > 0 {
		direction--
	}
	if inpututil.KeyPressDuration(c.downKey) > 0 {
		direction++
	}
	return Control{
		TargetVelocity: direction * 7,
		Acceleration:   paddleAcceleration,
		Braking:        1.5,
	}
}

type beginnerAI struct{}

// NewBeginnerAI creates an AI that follows the ball without predicting bounces.
func NewBeginnerAI() Controller {
	return &beginnerAI{}
}

func (*beginnerAI) Name() string {
	return "BEGINNER AI"
}

func (*beginnerAI) Control(view GameView) Control {
	target := view.BallY + ballRadius - view.PaddleHeight/2
	delta := target - view.PaddleY
	if !ballApproaching(view) {
		delta = view.Height/2 - view.PaddleHeight/2 - view.PaddleY
	}
	if math.Abs(delta) < 18 {
		delta = 0
	}
	return Control{
		TargetVelocity: min(max(delta*0.08, -4.2), 4.2),
		Acceleration:   0.35,
		Braking:        0.55,
	}
}

type humanLikeAI struct {
	targetY     float64
	cooldown    int
	tracking    bool
	initialized bool
	random      *rand.Rand
}

const (
	humanAIMaxSpeed          = 5.8
	humanAIAcceleration      = 0.45
	humanAIBraking           = 0.75
	humanAIDeadZone          = 7.0
	humanAIAimError          = 14.0
	humanAIDistanceAimError  = 22.0
	humanAIIdleAimError      = 12.0
	humanAIReactionMin       = 5
	humanAIReactionJitter    = 5
	humanAIDistanceDelay     = 7
	humanAIIdleDelay         = 24
	humanAIMistakeChance     = 0.08
	humanAIMistakeError      = 55.0
	humanAIVelocityError     = 0.08
	humanAIObservationYError = 5.0
)

// NewHumanLikeAI creates an imperfect AI with reaction time, observation
// errors, gradual movement, and occasional mistakes.
func NewHumanLikeAI() Controller {
	return newHumanLikeAI(rand.New(rand.NewPCG(rand.Uint64(), rand.Uint64())))
}

func newHumanLikeAI(random *rand.Rand) *humanLikeAI {
	return &humanLikeAI{random: random}
}

func (*humanLikeAI) Name() string {
	return "HUMAN AI"
}

func (c *humanLikeAI) Control(view GameView) Control {
	if !c.initialized {
		c.targetY = view.PaddleY
		c.initialized = true
	}
	approaching := ballApproaching(view)
	if approaching && !c.tracking {
		c.cooldown = c.reactionDelay(view)
	}
	c.tracking = approaching

	if c.cooldown > 0 {
		c.cooldown--
	}
	if c.cooldown == 0 {
		c.targetY = c.nextTarget(view)
		if approaching {
			c.cooldown = c.reactionDelay(view)
		} else {
			c.cooldown = humanAIIdleDelay + c.random.IntN(humanAIReactionJitter+1)
		}
	}

	delta := c.targetY - view.PaddleY
	targetVelocity := min(max(delta*0.12, -humanAIMaxSpeed), humanAIMaxSpeed)
	if math.Abs(delta) < humanAIDeadZone {
		targetVelocity = 0
	}
	return Control{
		TargetVelocity: targetVelocity,
		Acceleration:   humanAIAcceleration,
		Braking:        humanAIBraking,
	}
}

func (c *humanLikeAI) nextTarget(view GameView) float64 {
	if ballApproaching(view) {
		observed := view
		observed.BallY += c.randomError(humanAIObservationYError)
		observed.BallVelocityY *= 1 + c.randomError(humanAIVelocityError)

		distance := math.Abs(view.PaddleX - view.BallX)
		errorRange := humanAIAimError + humanAIDistanceAimError*distance/view.Width
		aimError := c.randomError(errorRange)
		if c.random.Float64() < humanAIMistakeChance {
			aimError += c.randomError(humanAIMistakeError)
		}
		target := predictBallY(observed) + ballRadius - view.PaddleHeight/2 + aimError
		return min(max(target, 0), view.Height-view.PaddleHeight)
	}
	return view.Height/2 - view.PaddleHeight/2 + c.randomError(humanAIIdleAimError)
}

func (c *humanLikeAI) reactionDelay(view GameView) int {
	distance := math.Abs(view.PaddleX - view.BallX)
	distanceDelay := int(humanAIDistanceDelay * min(distance/view.Width, 1))
	return humanAIReactionMin + distanceDelay + c.random.IntN(humanAIReactionJitter+1)
}

func (c *humanLikeAI) randomError(maximum float64) float64 {
	return (c.random.Float64()*2 - 1) * maximum
}

type perfectAI struct{}

// NewPerfectAI creates an AI with exact ball prediction and fast movement.
func NewPerfectAI() Controller {
	return &perfectAI{}
}

func (*perfectAI) Name() string {
	return "PERFECT AI"
}

func (*perfectAI) Control(view GameView) Control {
	target := view.Height/2 - view.PaddleHeight/2
	if ballApproaching(view) {
		target = predictBallY(view) + ballRadius - view.PaddleHeight/2
	}
	target = min(max(target, 0), view.Height-view.PaddleHeight)
	delta := target - view.PaddleY
	if math.Abs(delta) < 2 {
		delta = 0
	}
	return Control{
		TargetVelocity: min(max(delta*0.25, -7.5), 7.5),
		Acceleration:   1.2,
		Braking:        1.8,
	}
}

func ballApproaching(view GameView) bool {
	if view.Side == LeftSide {
		return view.BallVelocityX < 0
	}
	return view.BallVelocityX > 0
}

func predictBallY(view GameView) float64 {
	targetX := view.PaddleX + paddleWidth/2
	if view.Side == LeftSide {
		targetX += ballRadius
	} else {
		targetX -= ballRadius
	}
	if view.BallVelocityX == 0 {
		return view.BallY
	}
	timeToReach := (targetX - view.BallX) / view.BallVelocityX
	if timeToReach <= 0 {
		return view.BallY
	}

	predictedY := view.BallY + view.BallVelocityY*timeToReach
	minY := 0.0
	maxY := view.Height - ballDiameter
	for predictedY < minY || predictedY > maxY {
		if predictedY < minY {
			predictedY = minY + (minY - predictedY)
		} else {
			predictedY = maxY - (predictedY - maxY)
		}
	}
	return predictedY
}
