package config

// Config holds all game configuration values
type Config struct {
	Window WindowConfig `json:"window"`
	Game   GameConfig   `json:"game"`
	Audio  AudioConfig  `json:"audio"`
}

// WindowConfig holds window-related settings
type WindowConfig struct {
	Width     int32 `json:"width"`
	Height    int32 `json:"height"`
	Title     string `json:"title"`
	TargetFPS int32 `json:"target_fps"`
}

// GameConfig holds game-related settings
type GameConfig struct {
	MaxLevels         int32   `json:"max_levels"`
	BallBaseSpeed     float32 `json:"ball_base_speed"`
	BallSpeedIncrement float32 `json:"ball_speed_increment"`
	PaddleBaseSpeed   float32 `json:"paddle_base_speed"`
	BricksPerRow      int32   `json:"bricks_per_row"`
	BricksPerCol      int32   `json:"bricks_per_col"`
}

// AudioConfig holds audio-related settings
type AudioConfig struct {
	Enabled           bool   `json:"enabled"`
	PaddleHitSoundPath string `json:"paddle_hit_sound_path"`
	BrickHitSoundPath  string `json:"brick_hit_sound_path"`
}

// Default returns the default configuration
func Default() Config {
	return Config{
		Window: WindowConfig{
			Width:     768,
			Height:    1024,
			Title:     "Breakout",
			TargetFPS: 144,
		},
		Game: GameConfig{
			MaxLevels:         2,
			BallBaseSpeed:     0.4,
			BallSpeedIncrement: 1.1,
			PaddleBaseSpeed:   0.3,
			BricksPerRow:      14,
			BricksPerCol:      8,
		},
		Audio: AudioConfig{
			Enabled:           true,
			PaddleHitSoundPath: "assets/paddle_hit.wav",
			BrickHitSoundPath:  "assets/brick_hit.wav",
		},
	}
}