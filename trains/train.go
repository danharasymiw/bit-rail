package trains

import (
	"github.com/danharasymiw/bit-rail/types"
	"github.com/google/uuid"
)

//go:generate go run ../cmd/generator/main.go -type=Train -output=../message/train_gen.go
type Train struct {
	ID           uuid.UUID
	IsReversing  bool
	IsMoving     bool
	Speed        int
	Acceleration int

	Cars []*TrainCar
}

//go:generate go run ../cmd/generator/main.go -type=CarType -output=../message/car_type_gen.go
type CarType uint8

const (
	CarTypeLocomotive CarType = iota
	CarTypeCargo
	CarTypePassenger
)

//go:generate go run ../cmd/generator/main.go -type=TrainCar -output=../message/train_car_gen.go
type TrainCar struct {
	X, Y      uint16
	Direction types.Dir
	Type      CarType
}
