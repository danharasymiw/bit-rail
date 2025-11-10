package trains

import (
	"github.com/danharasymiw/bit-rail/types"
	"github.com/google/uuid"
)

type Train struct {
	ID           uuid.UUID
	IsReversing  bool
	IsMoving     bool
	Speed        int
	Acceleration int

	Cars []*TrainCar
}

type CarType uint

const (
	CarTypeLocomotive CarType = iota
	CarTypeCargo
	CarTypePassenger
)

type TrainCar struct {
	X, Y      int `binary:"uint16"`
	Direction types.Dir
	Type      CarType
}

func (t *Train) Tick(w trainWorldView) {
}
