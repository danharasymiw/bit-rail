package types

import (
	"github.com/google/uuid"
)

// Block represents a section of track
// Blocks are defined by the area between signals
// If one of these tiles in a section is removed/modified, a block needs to be recalculated
// Only one train is allowed to be inside of a block at a time
type Block struct {
	ID         uuid.UUID
	OccupiedBy Occupier
}

func NewBlock() *Block {
	return &Block{
		ID: uuid.New(),
	}
}

// Occupier is an entity that occupies a block
type Occupier interface {
	ID() string
}
