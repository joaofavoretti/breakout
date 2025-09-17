package renderer

import (
	"strconv"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	WindowWidth  = 768
	WindowHeight = 1024
)

// Renderer handles all rendering operations
type Renderer struct{}

// New creates a new renderer
func New() *Renderer {
	return &Renderer{}
}

// DrawScore renders the current score
func (r *Renderer) DrawScore(score int32) {
	rl.DrawText(strconv.Itoa(int(score)), 20, 20, 40, rl.RayWhite)
}

// DrawGameWon renders the game won screen
func (r *Renderer) DrawGameWon(score int32) {
	r.drawCenteredText("Game Won! Press R to Restart", WindowHeight/2, 20)
	r.drawCenteredText("Final Score: "+strconv.Itoa(int(score)), WindowHeight/2+40, 20)
}

// DrawGameLost renders the game lost screen
func (r *Renderer) DrawGameLost(score int32) {
	r.drawCenteredText("Game Lost! Press R to Restart", WindowHeight/2, 20)
	r.drawCenteredText("Final Score: "+strconv.Itoa(int(score)), WindowHeight/2+40, 20)
}

// DrawPaused renders the paused screen overlay
func (r *Renderer) DrawPaused() {
	r.drawCenteredText("Paused! Press Space to Resume", WindowHeight/2+40, 20)
}

func (r *Renderer) drawCenteredText(text string, y int32, fontSize int32) {
	textWidth := rl.MeasureText(text, fontSize)
	x := (WindowWidth - textWidth) / 2
	rl.DrawText(text, x, y, fontSize, rl.RayWhite)
}