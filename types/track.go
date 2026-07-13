package types

import (
	"math/bits"

	"github.com/google/uuid"
)

type Track struct {
	Direction Dir
	SignalDir Dir
	Block     *Block
}

func (t Track) HasSignal() bool {
	return t.SignalDir != 0
}

// IsJunction returns true if the track connects more than two directions,
// meaning a train arriving at it has more than one way to continue.
// Note: unlike a signal, a junction does NOT split a block (see
// IsSignalBoundary) — it's only a decision point for the routing graph (see
// the routing package). A diamond crossing is commonly protected as a
// single shared block guarded by signals on all its approaches, so two
// perpendicular trains see each other's occupancy; splitting the block at
// the junction itself would silently defeat that.
func (t Track) IsJunction() bool {
	return bits.OnesCount8(uint8(t.Direction)) > 2
}

// IsSignalBoundary reports whether a block-BFS walk crossing from curr to
// neighbour in dir should stop there, because one of the tiles has a signal
// facing the direction of travel. Blocks are contiguous sections of track
// between signals.
func IsSignalBoundary(curr, neighbour *Track, dir Dir) bool {
	if neighbour.HasSignal() && neighbour.IsSignalDir(dir) {
		return true
	}

	if curr.HasSignal() && curr.IsSignalDir(OppositeDir(dir)) {
		return true
	}

	return false
}

func (t Track) IsSignalClear(id uuid.UUID) bool {
	return t.Block.OccupiedBy == uuid.Nil || t.Block.OccupiedBy == id
}

func (t Track) IsSignalDir(dir Dir) bool {
	return t.SignalDir&dir != 0
}

func (t *Track) SetSignal(id uuid.UUID) {
	t.Block.OccupiedBy = id
}

func (t *Track) ClearSignal() {
	t.Block.OccupiedBy = uuid.Nil
}
