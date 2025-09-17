package main

import (
	"breakout/internal/game"
	"log"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	WindowWidth  = 768
	WindowHeight = 1024
	TargetFPS    = 144
)

func main() {
	// Initialize raylib
	rl.InitWindow(WindowWidth, WindowHeight, "Breakout")
	defer rl.CloseWindow()

	rl.InitAudioDevice()
	defer rl.CloseAudioDevice()

	rl.SetTargetFPS(TargetFPS)

	// Create and initialize game
	g, err := game.New()
	if err != nil {
		log.Fatalf("Failed to create game: %v", err)
	}
	defer g.Cleanup()

	g.Initialize()

	// Main game loop
	for !rl.WindowShouldClose() {
		deltaTime := rl.GetFrameTime()

		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)

		g.Update(deltaTime)
		g.Draw()

		rl.EndDrawing()
	}
}
