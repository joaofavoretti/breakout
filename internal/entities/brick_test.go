package entities

import (
	"testing"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func TestBrickGetValue(t *testing.T) {
	tests := []struct {
		name     string
		row      int32
		expected int32
	}{
		{"Top row (0)", 0, 7},
		{"Second row (1)", 1, 7},
		{"Third row (2)", 2, 5},
		{"Fourth row (3)", 3, 5},
		{"Fifth row (4)", 4, 3},
		{"Sixth row (5)", 5, 3},
		{"Seventh row (6)", 6, 1},
		{"Bottom row (7)", 7, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			brick := NewBrick(0, tt.row, rl.Red)
			if got := brick.GetValue(); got != tt.expected {
				t.Errorf("GetValue() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestBrickColorChecks(t *testing.T) {
	redBrick := NewBrick(0, 0, rl.Red)
	orangeBrick := NewBrick(0, 0, rl.Orange)
	greenBrick := NewBrick(0, 0, rl.Green)

	if !redBrick.IsRed() {
		t.Error("Red brick should return true for IsRed()")
	}
	if redBrick.IsOrange() {
		t.Error("Red brick should return false for IsOrange()")
	}

	if orangeBrick.IsRed() {
		t.Error("Orange brick should return false for IsRed()")
	}
	if !orangeBrick.IsOrange() {
		t.Error("Orange brick should return true for IsOrange()")
	}

	if greenBrick.IsRed() || greenBrick.IsOrange() {
		t.Error("Green brick should return false for both IsRed() and IsOrange()")
	}
}

func TestCreateLevelBricks(t *testing.T) {
	bricks := CreateLevelBricks()

	expectedCount := BricksPerRow * BricksPerCol
	if len(bricks) != expectedCount {
		t.Errorf("Expected %d bricks, got %d", expectedCount, len(bricks))
	}

	// Check that we have the right colors in the right rows
	colorCounts := make(map[string]int)
	for _, brick := range bricks {
		switch brick.color {
		case rl.Red:
			colorCounts["red"]++
		case rl.Orange:
			colorCounts["orange"]++
		case rl.Green:
			colorCounts["green"]++
		case rl.Yellow:
			colorCounts["yellow"]++
		}
	}

	// Each color should appear in exactly 2 rows * BricksPerRow bricks
	expectedPerColor := 2 * BricksPerRow
	for color, count := range colorCounts {
		if count != expectedPerColor {
			t.Errorf("Expected %d %s bricks, got %d", expectedPerColor, color, count)
		}
	}
}