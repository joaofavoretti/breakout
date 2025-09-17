package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	color "image/color"
	"math"
	"strconv"
)

const (
	WindowWidth  = 768
	WindowHeight = 1024
)

var ()

const (
	PlayerPaddleHeight = 20
	PlayerPaddleWidth  = 100
	PlayerPaddleYPos   = WindowHeight - 100
	PlayerBaseSpeed    = 0.3

	BricksPerRow  = 14
	BricksPerCol  = 8
	BricksSpacing = 5
	BricksYOffset = 200
	BrickHeight   = 10

	BallBaseSpeed      = 0.4
	BallSpeedIncrement = 1.1

	BallSize = 10
)

type ChangeStateConditions struct {
	UpperWallHit  bool
	OrangeContact bool
	RedContact    bool
	FourHits      bool
	TwelveHits    bool
}

func NewChangeStateConditions() ChangeStateConditions {
	return ChangeStateConditions{
		UpperWallHit:  false,
		OrangeContact: false,
		RedContact:    false,
		FourHits:      false,
		TwelveHits:    false,
	}
}

type GameState struct {
	Level            int32
	Collidables      []Collidable
	GameLost         bool
	GameWon          bool
	Paused           bool
	Player           PlayerPaddle
	Ball             Ball
	ChangeConditions ChangeStateConditions
	Bricks           []Brick
	Score            int32
	BrickHitCount    int32
}

func (s *GameState) isBricksEmpty() bool {
	return len(s.Bricks) == 0
}

func (s *GameState) IsGameOver() bool {
	return s.GameLost || s.GameWon
}

var state GameState

type Collidable interface {
	GetShape() rl.Rectangle
}

type IVector2 struct {
	X int32
	Y int32
}

type Brick struct {
	Color color.RGBA
	Pos   IVector2
}

func (b *Brick) Draw() {
	brickSize := (WindowWidth - (BricksPerRow+1)*BricksSpacing) / BricksPerRow
	x := b.Pos.X*int32(brickSize+BricksSpacing) + BricksSpacing
	y := b.Pos.Y*int32(BrickHeight+BricksSpacing) + BricksSpacing + BricksYOffset
	rl.DrawRectangle(x, y, int32(brickSize), BrickHeight, b.Color)
}

func (b *Brick) GetValue() int32 {
	return 2*int32((7-b.Pos.Y)/2) + 1
}

func (b *Brick) GetShape() rl.Rectangle {
	brickSize := (WindowWidth - (BricksPerRow+1)*BricksSpacing) / BricksPerRow
	x := b.Pos.X*int32(brickSize+BricksSpacing) + BricksSpacing
	y := b.Pos.Y*int32(BrickHeight+BricksSpacing) + BricksSpacing + BricksYOffset
	return rl.Rectangle{
		X:      float32(x),
		Y:      float32(y),
		Width:  float32(brickSize),
		Height: float32(BrickHeight),
	}
}

func NewBrick(x, y int32, color color.RGBA) *Brick {
	return &Brick{
		Pos:   IVector2{X: x, Y: y},
		Color: color,
	}
}

func CreateLevelBricks() []Brick {
	bricks := make([]Brick, 0)

	for i := range BricksPerRow {
		for j := range BricksPerCol {
			var color = rl.Red

			if j >= 6 {
				color = rl.Yellow
			} else if j >= 4 {
				color = rl.Green
			} else if j >= 2 {
				color = rl.Orange
			}

			brick := NewBrick(int32(i), int32(j), color)
			bricks = append(bricks, *brick)
		}
	}

	return bricks
}

type PlayerPaddle struct {
	Width      float32
	X          float32
	Speed      float32
	SpeedScale int32
}

func NewPlayerPaddle(x float32) PlayerPaddle {
	return PlayerPaddle{
		Width:      PlayerPaddleWidth,
		X:          x,
		Speed:      PlayerBaseSpeed * float32(2),
		SpeedScale: 2,
	}
}

func (p *PlayerPaddle) Draw() {
	px := int32(p.X*float32(WindowWidth) - p.Width/2)
	rl.DrawRectangle(px, PlayerPaddleYPos, int32(p.Width), PlayerPaddleHeight, rl.RayWhite)
}

func (p *PlayerPaddle) Update(timeDelta float32) {

	// Horizontal Move
	keyToDelta := map[int32]float32{
		rl.KeyA: -p.Speed,
		rl.KeyD: p.Speed,
	}

	for key, delta := range keyToDelta {
		if rl.IsKeyDown(key) {
			p.X += delta * timeDelta
		}
	}

	p.X = max(0, min(1, p.X))

	// Speed Scale Update
	keyToScale := map[int32]int32{
		rl.KeyW: 1,
		rl.KeyS: -1,
	}

	for key, scale := range keyToScale {
		if rl.IsKeyPressed(key) {
			p.SpeedScale += scale
			p.SpeedScale = max(1, min(5, p.SpeedScale))
			p.Speed = PlayerBaseSpeed * float32(p.SpeedScale)
		}
	}

}

func (p *PlayerPaddle) GetShape() rl.Rectangle {
	px := int32(p.X*float32(WindowWidth) - p.Width/2)
	return rl.Rectangle{
		X:      float32(px),
		Y:      float32(PlayerPaddleYPos),
		Width:  float32(p.Width),
		Height: float32(PlayerPaddleHeight),
	}
}

type Ball struct {
	Pos IVector2
	Vel rl.Vector2
}

func NewBall() Ball {
	return Ball{
		Pos: IVector2{X: WindowWidth / 2, Y: WindowHeight / 2},
		Vel: rl.Vector2{
			X: BallBaseSpeed,
			Y: BallBaseSpeed,
		},
	}
}

func (b *Ball) GetShape() rl.Rectangle {
	return rl.Rectangle{
		X:      float32(b.Pos.X),
		Y:      float32(b.Pos.Y),
		Width:  BallSize,
		Height: BallSize,
	}
}

func (b *Ball) CollidesWith(c Collidable) bool {
	return rl.CheckCollisionRecs(b.GetShape(), c.GetShape())
}

type CollisionAxis int

const (
	Vertical CollisionAxis = iota
	Horizontal
)

func GetCollisionAxis(b Collidable, c Collidable) CollisionAxis {
	bShape := b.GetShape()
	cShape := c.GetShape()
	bCenter := rl.Vector2{
		X: bShape.X + bShape.Width/2,
		Y: bShape.Y + bShape.Height/2,
	}
	cCenter := rl.Vector2{
		X: cShape.X + cShape.Width/2,
		Y: cShape.Y + cShape.Height/2,
	}
	absDiff := rl.Vector2{
		X: float32(math.Abs(float64(bCenter.X - cCenter.X))),
		Y: float32(math.Abs(float64(bCenter.Y - cCenter.Y))),
	}
	halfWidths := (bShape.Width + cShape.Width) / 2
	halfHeights := (bShape.Height + cShape.Height) / 2
	overlapX := halfWidths - absDiff.X
	overlapY := halfHeights - absDiff.Y
	if overlapX < overlapY {
		return Horizontal
	} else {
		return Vertical
	}
}

func (b *Ball) ReflectBrick(c Collidable) {
	axis := GetCollisionAxis(b, c)
	switch axis {
	case Vertical:
		b.Vel.Y = -b.Vel.Y
	case Horizontal:
		b.Vel.X = -b.Vel.X
	}
}

func (b *Ball) ReflectPaddle(c Collidable) {
	paddle := c.(*PlayerPaddle)
	paddleCenterX := paddle.X * float32(WindowWidth)
	ballCenterX := float32(b.Pos.X) + BallSize/2
	relativeIntersectX := (ballCenterX - paddleCenterX) / (state.Player.Width / 2)
	bounceAngle := relativeIntersectX * (5 * math.Pi / 12) // Max bounce angle of 75 degrees
	speed := float32(math.Sqrt(float64(b.Vel.X*b.Vel.X + b.Vel.Y*b.Vel.Y)))
	b.Vel.X = speed * float32(math.Sin(float64(bounceAngle)))
	b.Vel.Y = -speed * float32(math.Cos(float64(bounceAngle)))
}

func (b *Ball) Draw() {
	rl.DrawRectangle(b.Pos.X, b.Pos.Y, BallSize, BallSize, rl.RayWhite)
}

func (b *Ball) Update(timeDelta float32) {
	b.Pos.X += int32(b.Vel.X * timeDelta * float32(WindowWidth))
	b.Pos.Y += int32(b.Vel.Y * timeDelta * float32(WindowHeight))

	if b.Pos.X <= 0 || b.Pos.X+BallSize >= WindowWidth {
		b.Vel.X = -b.Vel.X
	}

	if b.Pos.Y <= 0 {
		b.Vel.Y = -b.Vel.Y
	}

	if b.Pos.Y+BallSize >= WindowHeight {
		state.GameLost = true
		return
	}

	for c := range state.Collidables {
		collidable := state.Collidables[c]
		if collidable != b && b.CollidesWith(collidable) {

			if brick, ok := collidable.(*Brick); ok {
				b.ReflectBrick(collidable)

				rl.PlaySound(BrickHitSound)

				state.Score += brick.GetValue()
				state.BrickHitCount += 1

				if brick.Color == rl.Red && !state.ChangeConditions.RedContact {
					state.ChangeConditions.RedContact = true
					state.Ball.Vel.X *= BallSpeedIncrement
					state.Ball.Vel.Y *= BallSpeedIncrement
				} else if brick.Color == rl.Orange && !state.ChangeConditions.OrangeContact {
					state.ChangeConditions.OrangeContact = true
					state.Ball.Vel.X *= BallSpeedIncrement
					state.Ball.Vel.Y *= BallSpeedIncrement
				}

				// Remove brick from game
				for i := range state.Bricks {
					if &state.Bricks[i] == brick {
						state.Bricks = append(state.Bricks[:i], state.Bricks[i+1:]...)
						break
					}
				}

				// Rebuild collidables list
				state.Collidables = make([]Collidable, 0)
				state.Collidables = append(state.Collidables, &state.Player)
				state.Collidables = append(state.Collidables, &state.Ball)
				for i := range state.Bricks {
					state.Collidables = append(state.Collidables, &state.Bricks[i])
				}
			} else if _, ok := collidable.(*PlayerPaddle); ok {
				b.ReflectPaddle(collidable)
				rl.PlaySound(PaddleHitSound)
			}

			break
		}
	}

	if state.BrickHitCount >= 4 && !state.ChangeConditions.FourHits {
		state.ChangeConditions.FourHits = true
		state.Ball.Vel.X *= BallSpeedIncrement
		state.Ball.Vel.Y *= BallSpeedIncrement
	}

	if state.BrickHitCount >= 12 && !state.ChangeConditions.TwelveHits {
		state.ChangeConditions.TwelveHits = true
		state.Ball.Vel.X *= BallSpeedIncrement
		state.Ball.Vel.Y *= BallSpeedIncrement
	}

	if b.Pos.Y <= 0 && !state.ChangeConditions.UpperWallHit {
		state.ChangeConditions.UpperWallHit = true
		state.Player.Width /= 2
	}
}

func Setup() {
	state.Level = 1
	state.Player = NewPlayerPaddle(0.5)
	state.Score = 0
	state.GameLost = false
	state.GameWon = false
	state.Paused = true
	state.Bricks = CreateLevelBricks()
	state.ChangeConditions = NewChangeStateConditions()
	state.Ball = NewBall()

	state.Collidables = make([]Collidable, 0)
	state.Collidables = append(state.Collidables, &state.Player)
	state.Collidables = append(state.Collidables, &state.Ball)
	for i := range state.Bricks {
		state.Collidables = append(state.Collidables, &state.Bricks[i])
	}

	PaddleHitSound = rl.LoadSound("assets/paddle_hit.wav")
	BrickHitSound = rl.LoadSound("assets/brick_hit.wav")
}

func Update(timeDelta float32) {
	if state.isBricksEmpty() && state.Level <= 2 {
		state.Level += 1

		state.Bricks = CreateLevelBricks()

		state.Collidables = make([]Collidable, 0)
		state.Collidables = append(state.Collidables, &state.Player)
		state.Collidables = append(state.Collidables, &state.Ball)
		for i := range state.Bricks {
			state.Collidables = append(state.Collidables, &state.Bricks[i])
		}
	}

	if state.Level > 2 {
		state.GameWon = true
		return
	}

	if state.IsGameOver() {
		if rl.IsKeyPressed(rl.KeyR) {
			Setup()
		}
		return
	}

	if state.Paused {
		if rl.IsKeyPressed(rl.KeySpace) {
			state.Paused = false
		}
		return
	}

	state.Player.Update(timeDelta)

	state.Ball.Update(timeDelta)
}

func Draw() {
	if state.GameWon {
		rl.DrawText("Game Won! Press R to Restart", WindowWidth/2-180, WindowHeight/2, 20, rl.RayWhite)
		rl.DrawText("Final Score: "+strconv.Itoa(int(state.Score)), WindowWidth/2-100, WindowHeight/2+40, 20, rl.RayWhite)
		return
	}

	if state.GameLost {
		rl.DrawText("Game Lost! Press R to Restart", WindowWidth/2-180, WindowHeight/2, 20, rl.RayWhite)
		rl.DrawText("Final Score: "+strconv.Itoa(int(state.Score)), WindowWidth/2-100, WindowHeight/2+40, 20, rl.RayWhite)
		return
	}

	if state.Paused {
		rl.DrawText("Paused! Press Space to Resume", WindowWidth/2-150, WindowHeight/2+40, 20, rl.RayWhite)
	}

	rl.DrawText(strconv.Itoa(int(state.Score)), 20, 20, 40, rl.RayWhite)

	state.Player.Draw()

	state.Ball.Draw()

	for _, brick := range state.Bricks {
		brick.Draw()
	}
}

var PaddleHitSound rl.Sound
var BrickHitSound rl.Sound

func main() {
	rl.InitWindow(WindowWidth, WindowHeight, "Breakout")
	defer rl.CloseWindow()

	rl.InitAudioDevice()
	defer rl.CloseAudioDevice()

	rl.SetTargetFPS(144)

	Setup()

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()

		rl.ClearBackground(rl.Black)

		timeDelta := rl.GetFrameTime()

		Update(timeDelta)

		Draw()

		rl.EndDrawing()
	}
}
