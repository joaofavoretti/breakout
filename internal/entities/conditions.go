package entities

// ChangeStateConditions tracks various game state change conditions
type ChangeStateConditions struct {
	UpperWallHit  bool
	OrangeContact bool
	RedContact    bool
	FourHits      bool
	TwelveHits    bool
}

// NewChangeStateConditions creates a new set of change state conditions
func NewChangeStateConditions() *ChangeStateConditions {
	return &ChangeStateConditions{
		UpperWallHit:  false,
		OrangeContact: false,
		RedContact:    false,
		FourHits:      false,
		TwelveHits:    false,
	}
}