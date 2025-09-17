package physics

import (
	"breakout/internal/types"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Engine handles collision detection
type Engine struct{}

// New creates a new physics engine
func New() *Engine {
	return &Engine{}
}

// CheckCollision checks if two collidable objects are colliding
func (e *Engine) CheckCollision(a, b types.Collidable) bool {
	boundsA := a.GetBounds()
	boundsB := b.GetBounds()
	
	return rl.CheckCollisionRecs(boundsA.ToRaylib(), boundsB.ToRaylib())
}