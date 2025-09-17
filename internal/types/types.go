package types

import rl "github.com/gen2brain/raylib-go/raylib"

// Vector2 represents a 2D vector with integer coordinates
type Vector2 struct {
	X, Y int32
}

// ToRaylib converts to raylib Vector2
func (v Vector2) ToRaylib() rl.Vector2 {
	return rl.Vector2{X: float32(v.X), Y: float32(v.Y)}
}

// Rectangle represents a rectangle shape
type Rectangle struct {
	X, Y          float32
	Width, Height float32
}

// ToRaylib converts to raylib Rectangle
func (r Rectangle) ToRaylib() rl.Rectangle {
	return rl.Rectangle{
		X:      r.X,
		Y:      r.Y,
		Width:  r.Width,
		Height: r.Height,
	}
}

// Collidable represents any object that can participate in collision detection
type Collidable interface {
	GetBounds() Rectangle
}

// Drawable represents any object that can be drawn
type Drawable interface {
	Draw()
}

// Updatable represents any object that can be updated
type Updatable interface {
	Update(deltaTime float32)
}