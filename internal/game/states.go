package game

// GameStateType represents different game states
type GameStateType int

const (
	StateMenu GameStateType = iota
	StatePlaying
	StatePaused
	StateGameOver
	StateGameWon
	StateSettings
)

// GameStateManager handles transitions between different game states
type GameStateManager struct {
	currentState GameStateType
	previousState GameStateType
}

// NewGameStateManager creates a new state manager
func NewGameStateManager() *GameStateManager {
	return &GameStateManager{
		currentState: StateMenu,
		previousState: StateMenu,
	}
}

// CurrentState returns the current game state
func (gsm *GameStateManager) CurrentState() GameStateType {
	return gsm.currentState
}

// PreviousState returns the previous game state
func (gsm *GameStateManager) PreviousState() GameStateType {
	return gsm.previousState
}

// TransitionTo changes to a new game state
func (gsm *GameStateManager) TransitionTo(newState GameStateType) {
	gsm.previousState = gsm.currentState
	gsm.currentState = newState
}

// IsPlaying returns true if the game is in playing state
func (gsm *GameStateManager) IsPlaying() bool {
	return gsm.currentState == StatePlaying
}

// IsPaused returns true if the game is paused
func (gsm *GameStateManager) IsPaused() bool {
	return gsm.currentState == StatePaused
}

// IsGameOver returns true if the game is over (lost or won)
func (gsm *GameStateManager) IsGameOver() bool {
	return gsm.currentState == StateGameOver || gsm.currentState == StateGameWon
}