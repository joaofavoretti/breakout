package entities

import (
	"breakout/internal/types"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	BallSize           = 10
	BallBaseSpeed      = 0.4
	BallSpeedIncrement = 1.1
	WindowWidth        = 768
	WindowHeight       = 1024
)

// Ball represents the game ball
type Ball struct {
	pos      types.Vector2
	velocity rl.Vector2
}

// NewBall creates a new ball at the center of the screen
func NewBall() *Ball {
	return &Ball{
		pos: types.Vector2{X: WindowWidth / 2, Y: WindowHeight / 2},
		velocity: rl.Vector2{
			X: BallBaseSpeed,
			Y: BallBaseSpeed,
		},
	}
}

// Position returns the current position
func (b *Ball) Position() types.Vector2 {
	return b.pos
}

// Velocity returns the current velocity
func (b *Ball) Velocity() rl.Vector2 {
	return b.velocity
}

// GetBounds returns the collision bounds
func (b *Ball) GetBounds() types.Rectangle {
	return types.Rectangle{
		X:      float32(b.pos.X),
		Y:      float32(b.pos.Y),
		Width:  BallSize,
		Height: BallSize,
	}
}

// Draw renders the ball
func (b *Ball) Draw() {
	rl.DrawRectangle(b.pos.X, b.pos.Y, BallSize, BallSize, rl.RayWhite)
}

// Update moves the ball and handles wall collisions
func (b *Ball) Update(deltaTime float32) {
	b.pos.X += int32(b.velocity.X * deltaTime * float32(WindowWidth))
	b.pos.Y += int32(b.velocity.Y * deltaTime * float32(WindowHeight))

	// Handle wall collisions
	if b.pos.X <= 0 || b.pos.X+BallSize >= WindowWidth {
		b.velocity.X = -b.velocity.X
	}

	if b.pos.Y <= 0 {
		b.velocity.Y = -b.velocity.Y
	}
}

// ReflectOffBrick reflects the ball off a brick
func (b *Ball) ReflectOffBrick(brick *Brick) {
	axis := b.getCollisionAxis(brick)
	switch axis {
	case CollisionAxisVertical:
		b.velocity.Y = -b.velocity.Y
	case CollisionAxisHorizontal:
		b.velocity.X = -b.velocity.X
	}
}

// ReflectOffPaddle reflects the ball off the paddle with angle variation
func (b *Ball) ReflectOffPaddle(paddle *PlayerPaddle) {
	paddleCenterX := paddle.X() * float32(WindowWidth)
	ballCenterX := float32(b.pos.X) + BallSize/2
	relativeIntersectX := (ballCenterX - paddleCenterX) / (paddle.Width() / 2)
	bounceAngle := relativeIntersectX * (5 * math.Pi / 12) // Max bounce angle of 75 degrees
	
	speed := float32(math.Sqrt(float64(b.velocity.X*b.velocity.X + b.velocity.Y*b.velocity.Y)))
	b.velocity.X = speed * float32(math.Sin(float64(bounceAngle)))
	b.velocity.Y = -speed * float32(math.Cos(float64(bounceAngle)))
}

// IncreaseSpeed multiplies the current speed by the given factor
func (b *Ball) IncreaseSpeed(factor float32) {
	b.velocity.X *= factor
	b.velocity.Y *= factor
}

type CollisionAxis int

const (
	CollisionAxisVertical CollisionAxis = iota
	CollisionAxisHorizontal
)

func (b *Ball) getCollisionAxis(other types.Collidable) CollisionAxis {
	bBounds := b.GetBounds()
	oBounds := other.GetBounds()
	
	bCenter := rl.Vector2{
		X: bBounds.X + bBounds.Width/2,
		Y: bBounds.Y + bBounds.Height/2,
	}
	oCenter := rl.Vector2{
		X: oBounds.X + oBounds.Width/2,
		Y: oBounds.Y + oBounds.Height/2,
	}
	
	absDiff := rl.Vector2{
		X: float32(math.Abs(float64(bCenter.X - oCenter.X))),
		Y: float32(math.Abs(float64(bCenter.Y - oCenter.Y))),
	}
	
	halfWidths := (bBounds.Width + oBounds.Width) / 2
	halfHeights := (bBounds.Height + oBounds.Height) / 2
	overlapX := halfWidths - absDiff.X
	overlapY := halfHeights - absDiff.Y
	
	if overlapX < overlapY {
		return CollisionAxisHorizontal
	}
	return CollisionAxisVertical
}