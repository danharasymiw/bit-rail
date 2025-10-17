package types

import (
	"github.com/google/uuid"
)

type BlockID uuid.UUID

// Block represents a section of track
// Blocks are defined by the area between signals
// If one of these tiles in a section is removed/modified, a block needs to be recalculated
// Only one train is allowed to be inside of a block at a time
type Block struct {
	ID         BlockID
	OccupiedBy Occupier
}

func NewBlock() *Block {
	return &Block{
		ID: BlockID(uuid.New()),
	}
}

// Occupier is an entity that occupies a block
type Occupier interface {
	ID() string
}
