package trains

import "github.com/danharasymiw/trains/types"

type Train struct {
	IsReversing  bool
	IsMoving     bool
	Speed        int
	Acceleration int

	Cars []*TrainCar
}

type CarType uint8

const (
	CarTypeLocomotive CarType = iota
	CarTypeCargo
	CarTypePassenger
)

type TrainCar struct {
	X, Y      int
	Direction types.Dir
	Type      CarType
}
