package game

import (
	"breakout/internal/audio"
	"breakout/internal/entities"
	"breakout/internal/physics"
	"breakout/internal/renderer"
	"breakout/internal/types"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	WindowWidth  = 768
	WindowHeight = 1024
	MaxLevels    = 2
)

// Game represents the main game state and logic
type Game struct {
	state    *State
	renderer *renderer.Renderer
	audio    *audio.Manager
	physics  *physics.Engine
}

// State holds the current game state
type State struct {
	Level         int32
	Score         int32
	BrickHitCount int32
	GameLost      bool
	GameWon       bool
	Paused        bool

	Player *entities.PlayerPaddle
	Ball   *entities.Ball
	Bricks []*entities.Brick

	ChangeConditions *entities.ChangeStateConditions
}

// New creates a new game instance
func New() (*Game, error) {
	audioManager, err := audio.New()
	if err != nil {
		return nil, err
	}

	return &Game{
		state:    &State{},
		renderer: renderer.New(),
		audio:    audioManager,
		physics:  physics.New(),
	}, nil
}

// Initialize sets up the initial game state
func (g *Game) Initialize() {
	g.state.Level = 1
	g.state.Score = 0
	g.state.BrickHitCount = 0
	g.state.GameLost = false
	g.state.GameWon = false
	g.state.Paused = true

	g.state.Player = entities.NewPlayerPaddle(0.5)
	g.state.Ball = entities.NewBall()
	g.state.Bricks = entities.CreateLevelBricks()
	g.state.ChangeConditions = entities.NewChangeStateConditions()
}

// Update handles game logic updates
func (g *Game) Update(deltaTime float32) {
	if g.isLevelComplete() && g.state.Level <= MaxLevels {
		g.advanceLevel()
	}

	if g.state.Level > MaxLevels {
		g.state.GameWon = true
		return
	}

	if g.isGameOver() {
		if rl.IsKeyPressed(rl.KeyR) {
			g.Initialize()
		}
		return
	}

	if g.state.Paused {
		if rl.IsKeyPressed(rl.KeySpace) {
			g.state.Paused = false
		}
		return
	}

	g.state.Player.Update(deltaTime)
	g.updateBall(deltaTime)
}

// Draw renders the current game state
func (g *Game) Draw() {
	if g.state.GameWon {
		g.renderer.DrawGameWon(g.state.Score)
		return
	}

	if g.state.GameLost {
		g.renderer.DrawGameLost(g.state.Score)
		return
	}

	if g.state.Paused {
		g.renderer.DrawPaused()
	}

	g.renderer.DrawScore(g.state.Score)
	g.state.Player.Draw()
	g.state.Ball.Draw()

	for _, brick := range g.state.Bricks {
		brick.Draw()
	}
}

// Cleanup releases game resources
func (g *Game) Cleanup() {
	g.audio.Cleanup()
}

func (g *Game) isLevelComplete() bool {
	return len(g.state.Bricks) == 0
}

func (g *Game) isGameOver() bool {
	return g.state.GameLost || g.state.GameWon
}

func (g *Game) advanceLevel() {
	g.state.Level++
	g.state.Bricks = entities.CreateLevelBricks()
}

func (g *Game) updateBall(deltaTime float32) {
	oldPos := g.state.Ball.Position()

	g.state.Ball.Update(deltaTime)

	// Check wall collisions
	if g.state.Ball.Position().Y <= 0 && !g.state.ChangeConditions.UpperWallHit {
		g.state.ChangeConditions.UpperWallHit = true
		g.state.Player.HalveWidth()
	}

	// Check game over condition
	if g.state.Ball.Position().Y+entities.BallSize >= WindowHeight {
		g.state.GameLost = true
		return
	}

	// Check collisions with game objects
	g.handleCollisions(oldPos)

	// Check speed increase conditions
	g.checkSpeedIncreaseConditions()
}

func (g *Game) handleCollisions(oldBallPos types.Vector2) {
	// Check paddle collision
	if g.physics.CheckCollision(g.state.Ball, g.state.Player) {
		g.state.Ball.ReflectOffPaddle(g.state.Player)
		g.audio.PlayPaddleHit()
		return
	}

	// Check brick collisions
	for i, brick := range g.state.Bricks {
		if g.physics.CheckCollision(g.state.Ball, brick) {
			g.state.Ball.ReflectOffBrick(brick)
			g.audio.PlayBrickHit()

			g.state.Score += brick.GetValue()
			g.state.BrickHitCount++

			// Handle special brick effects
			g.handleBrickEffects(brick)

			// Remove brick
			g.state.Bricks = append(g.state.Bricks[:i], g.state.Bricks[i+1:]...)
			return
		}
	}
}

func (g *Game) handleBrickEffects(brick *entities.Brick) {
	if brick.IsRed() && !g.state.ChangeConditions.RedContact {
		g.state.ChangeConditions.RedContact = true
		g.state.Ball.IncreaseSpeed(entities.BallSpeedIncrement)
	} else if brick.IsOrange() && !g.state.ChangeConditions.OrangeContact {
		g.state.ChangeConditions.OrangeContact = true
		g.state.Ball.IncreaseSpeed(entities.BallSpeedIncrement)
	}
}

func (g *Game) checkSpeedIncreaseConditions() {
	if g.state.BrickHitCount >= 4 && !g.state.ChangeConditions.FourHits {
		g.state.ChangeConditions.FourHits = true
		g.state.Ball.IncreaseSpeed(entities.BallSpeedIncrement)
	}

	if g.state.BrickHitCount >= 12 && !g.state.ChangeConditions.TwelveHits {
		g.state.ChangeConditions.TwelveHits = true
		g.state.Ball.IncreaseSpeed(entities.BallSpeedIncrement)
	}
}