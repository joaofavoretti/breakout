package audio

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

// Manager handles all audio operations
type Manager struct {
	paddleHitSound rl.Sound
	brickHitSound  rl.Sound
}

// New creates a new audio manager and loads sounds
func New() (*Manager, error) {
	paddleSound := rl.LoadSound("assets/paddle_hit.wav")
	brickSound := rl.LoadSound("assets/brick_hit.wav")

	return &Manager{
		paddleHitSound: paddleSound,
		brickHitSound:  brickSound,
	}, nil
}

// PlayPaddleHit plays the paddle hit sound
func (m *Manager) PlayPaddleHit() {
	rl.PlaySound(m.paddleHitSound)
}

// PlayBrickHit plays the brick hit sound
func (m *Manager) PlayBrickHit() {
	rl.PlaySound(m.brickHitSound)
}

// Cleanup unloads all sounds
func (m *Manager) Cleanup() {
	rl.UnloadSound(m.paddleHitSound)
	rl.UnloadSound(m.brickHitSound)
}