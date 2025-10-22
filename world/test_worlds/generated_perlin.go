package test_worlds

import (
	"github.com/danharasymiw/bit-rail/trains"
	"github.com/danharasymiw/bit-rail/types"
	"github.com/danharasymiw/bit-rail/world"
)

func NewPerlinWorld(seed, smoothness int64) *world.World {
	w := world.New(500, 500)
	world.Generate(w, seed)

	w.AddTrack(0, 0, types.DirNorth|types.DirSouth|types.DirEast|types.DirWest)

	for x := 35; x < 80; x++ {
		w.AddTrack(x, 20, types.DirEast|types.DirWest)
		w.AddTrack(x, 55, types.DirEast|types.DirWest)
	}
	for y := 20; y < 56; y++ {
		w.AddTrack(35, y, types.DirNorth|types.DirSouth)
		w.AddTrack(80, y, types.DirNorth|types.DirSouth)
	}
	w.AddTrack(35, 20, types.DirNorth|types.DirEast)
	w.AddTrack(80, 20, types.DirNorth|types.DirWest)
	w.AddTrack(35, 55, types.DirSouth|types.DirEast)
	w.AddTrack(80, 55, types.DirSouth|types.DirWest)

	w.AddTrain(&trains.Train{
		IsMoving: true,
		Cars: []*trains.TrainCar{
			{X: 50, Y: 20, Type: trains.CarTypeLocomotive, Direction: types.DirWest},
			{X: 51, Y: 20, Type: trains.CarTypeCargo, Direction: types.DirWest},
			{X: 52, Y: 20, Type: trains.CarTypeCargo, Direction: types.DirWest},
			{X: 53, Y: 20, Type: trains.CarTypeCargo, Direction: types.DirWest},
			{X: 54, Y: 20, Type: trains.CarTypeCargo, Direction: types.DirWest},
			{X: 55, Y: 20, Type: trains.CarTypeCargo, Direction: types.DirWest},
			{X: 56, Y: 20, Type: trains.CarTypeCargo, Direction: types.DirWest},
			{X: 57, Y: 20, Type: trains.CarTypeCargo, Direction: types.DirWest},
			{X: 58, Y: 20, Type: trains.CarTypeCargo, Direction: types.DirWest},
			{X: 59, Y: 20, Type: trains.CarTypeCargo, Direction: types.DirWest},
			{X: 60, Y: 20, Type: trains.CarTypeCargo, Direction: types.DirWest},
			{X: 61, Y: 20, Type: trains.CarTypeCargo, Direction: types.DirWest},
			{X: 62, Y: 20, Type: trains.CarTypeCargo, Direction: types.DirWest},
			{X: 63, Y: 20, Type: trains.CarTypeCargo, Direction: types.DirWest},
			{X: 64, Y: 20, Type: trains.CarTypeCargo, Direction: types.DirWest},
			{X: 65, Y: 20, Type: trains.CarTypeCargo, Direction: types.DirWest},
			{X: 66, Y: 20, Type: trains.CarTypeCargo, Direction: types.DirWest},
			{X: 67, Y: 20, Type: trains.CarTypeCargo, Direction: types.DirWest},
			{X: 68, Y: 20, Type: trains.CarTypeCargo, Direction: types.DirWest},
			{X: 69, Y: 20, Type: trains.CarTypeCargo, Direction: types.DirWest},
			{X: 70, Y: 20, Type: trains.CarTypeCargo, Direction: types.DirWest},
			{X: 71, Y: 20, Type: trains.CarTypeCargo, Direction: types.DirWest},
			{X: 72, Y: 20, Type: trains.CarTypeCargo, Direction: types.DirWest},
			{X: 73, Y: 20, Type: trains.CarTypeCargo, Direction: types.DirWest},
			{X: 74, Y: 20, Type: trains.CarTypeCargo, Direction: types.DirWest},
			{X: 75, Y: 20, Type: trains.CarTypeCargo, Direction: types.DirWest},
			{X: 76, Y: 20, Type: trains.CarTypeCargo, Direction: types.DirWest},
			{X: 77, Y: 20, Type: trains.CarTypeCargo, Direction: types.DirWest},
			{X: 78, Y: 20, Type: trains.CarTypeCargo, Direction: types.DirWest},
			{X: 79, Y: 20, Type: trains.CarTypeCargo, Direction: types.DirWest},
		},
	})

	for x := 65; x < 90; x++ {
		w.AddTrack(x, 40, types.DirEast|types.DirWest)
		w.AddTrack(x, 60, types.DirEast|types.DirWest)
	}
	for y := 40; y < 61; y++ {
		w.AddTrack(65, y, types.DirNorth|types.DirSouth)
		w.AddTrack(90, y, types.DirNorth|types.DirSouth)
	}
	w.AddTrack(65, 40, types.DirNorth|types.DirEast)
	w.AddTrack(90, 40, types.DirNorth|types.DirWest)
	w.AddTrack(65, 60, types.DirSouth|types.DirEast)
	w.AddTrack(90, 60, types.DirSouth|types.DirWest)

	w.AddTrack(80, 40, types.DirSouth|types.DirNorth|types.DirEast|types.DirWest)
	w.AddTrack(65, 55, types.DirSouth|types.DirNorth|types.DirEast|types.DirWest)

	w.AddTrain(&trains.Train{
		IsMoving: true,
		Cars: []*trains.TrainCar{
			{X: 66, Y: 60, Type: trains.CarTypeLocomotive, Direction: types.DirWest},
			{X: 51, Y: 60, Type: trains.CarTypeCargo, Direction: types.DirWest},
			{X: 52, Y: 60, Type: trains.CarTypeCargo, Direction: types.DirWest},
			{X: 53, Y: 60, Type: trains.CarTypeCargo, Direction: types.DirWest},
			{X: 54, Y: 60, Type: trains.CarTypeCargo, Direction: types.DirWest},
			{X: 55, Y: 60, Type: trains.CarTypeCargo, Direction: types.DirWest},
			{X: 56, Y: 60, Type: trains.CarTypeCargo, Direction: types.DirWest},
			{X: 57, Y: 60, Type: trains.CarTypeCargo, Direction: types.DirWest},
			{X: 58, Y: 60, Type: trains.CarTypeCargo, Direction: types.DirWest},
			{X: 59, Y: 60, Type: trains.CarTypeCargo, Direction: types.DirWest},
		},
	})
	return w
}
