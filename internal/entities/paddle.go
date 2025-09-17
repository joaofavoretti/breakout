package entities

import (
	"breakout/internal/types"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	PlayerPaddleHeight = 20
	PlayerPaddleWidth  = 100
	PlayerPaddleYPos   = WindowHeight - 100
	PlayerBaseSpeed    = 0.3
)

// PlayerPaddle represents the player's paddle
type PlayerPaddle struct {
	width      float32
	x          float32
	speed      float32
	speedScale int32
}

// NewPlayerPaddle creates a new player paddle
func NewPlayerPaddle(x float32) *PlayerPaddle {
	return &PlayerPaddle{
		width:      PlayerPaddleWidth,
		x:          x,
		speed:      PlayerBaseSpeed * 2,
		speedScale: 2,
	}
}

// X returns the paddle's X position (normalized 0-1)
func (p *PlayerPaddle) X() float32 {
	return p.x
}

// Width returns the paddle's width
func (p *PlayerPaddle) Width() float32 {
	return p.width
}

// GetBounds returns the collision bounds
func (p *PlayerPaddle) GetBounds() types.Rectangle {
	px := p.x*float32(WindowWidth) - p.width/2
	return types.Rectangle{
		X:      px,
		Y:      float32(PlayerPaddleYPos),
		Width:  p.width,
		Height: float32(PlayerPaddleHeight),
	}
}

// Draw renders the paddle
func (p *PlayerPaddle) Draw() {
	px := int32(p.x*float32(WindowWidth) - p.width/2)
	rl.DrawRectangle(px, PlayerPaddleYPos, int32(p.width), PlayerPaddleHeight, rl.RayWhite)
}

// Update handles paddle movement and speed changes
func (p *PlayerPaddle) Update(deltaTime float32) {
	p.handleMovement(deltaTime)
	p.handleSpeedChange()
}

// HalveWidth reduces the paddle width by half
func (p *PlayerPaddle) HalveWidth() {
	p.width /= 2
}

func (p *PlayerPaddle) handleMovement(deltaTime float32) {
	keyToDelta := map[int32]float32{
		rl.KeyA: -p.speed,
		rl.KeyD: p.speed,
	}

	for key, delta := range keyToDelta {
		if rl.IsKeyDown(key) {
			p.x += delta * deltaTime
		}
	}

	// Clamp position to screen bounds
	p.x = max(0, min(1, p.x))
}

func (p *PlayerPaddle) handleSpeedChange() {
	keyToScale := map[int32]int32{
		rl.KeyW: 1,
		rl.KeyS: -1,
	}

	for key, scale := range keyToScale {
		if rl.IsKeyPressed(key) {
			p.speedScale += scale
			p.speedScale = max(1, min(5, p.speedScale))
			p.speed = PlayerBaseSpeed * float32(p.speedScale)
		}
	}
}