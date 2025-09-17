package entities

import (
	"breakout/internal/types"
	"image/color"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	BricksPerRow  = 14
	BricksPerCol  = 8
	BricksSpacing = 5
	BricksYOffset = 200
	BrickHeight   = 10
)

// Brick represents a destructible brick
type Brick struct {
	color color.RGBA
	pos   types.Vector2
}

// NewBrick creates a new brick at the specified grid position
func NewBrick(x, y int32, color color.RGBA) *Brick {
	return &Brick{
		pos:   types.Vector2{X: x, Y: y},
		color: color,
	}
}

// GetBounds returns the collision bounds
func (b *Brick) GetBounds() types.Rectangle {
	brickSize := (WindowWidth - (BricksPerRow+1)*BricksSpacing) / BricksPerRow
	x := float32(b.pos.X*int32(brickSize+BricksSpacing) + BricksSpacing)
	y := float32(b.pos.Y*int32(BrickHeight+BricksSpacing) + BricksSpacing + BricksYOffset)
	
	return types.Rectangle{
		X:      x,
		Y:      y,
		Width:  float32(brickSize),
		Height: float32(BrickHeight),
	}
}

// Draw renders the brick
func (b *Brick) Draw() {
	bounds := b.GetBounds()
	rl.DrawRectangle(
		int32(bounds.X),
		int32(bounds.Y),
		int32(bounds.Width),
		int32(bounds.Height),
		b.color,
	)
}

// GetValue returns the point value of the brick based on its row
func (b *Brick) GetValue() int32 {
	return 2*int32((7-b.pos.Y)/2) + 1
}

// IsRed returns true if the brick is red
func (b *Brick) IsRed() bool {
	return b.color == rl.Red
}

// IsOrange returns true if the brick is orange
func (b *Brick) IsOrange() bool {
	return b.color == rl.Orange
}

// CreateLevelBricks creates all bricks for a level
func CreateLevelBricks() []*Brick {
	bricks := make([]*Brick, 0, BricksPerRow*BricksPerCol)

	for i := 0; i < BricksPerRow; i++ {
		for j := 0; j < BricksPerCol; j++ {
			color := getBrickColor(j)
			brick := NewBrick(int32(i), int32(j), color)
			bricks = append(bricks, brick)
		}
	}

	return bricks
}

func getBrickColor(row int) color.RGBA {
	switch {
	case row >= 6:
		return rl.Yellow
	case row >= 4:
		return rl.Green
	case row >= 2:
		return rl.Orange
	default:
		return rl.Red
	}
}