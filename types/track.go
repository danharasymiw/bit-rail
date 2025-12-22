package types

import (
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
