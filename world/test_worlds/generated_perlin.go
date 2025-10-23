package test_worlds

import (
	"github.com/danharasymiw/bit-rail/trains"
	"github.com/danharasymiw/bit-rail/types"
	"github.com/danharasymiw/bit-rail/world"
)

func NewPerlinWorld(seed, smoothness int64) *world.World {
	w := world.New(500, 500)
	world.Generate(w, seed)

	w.AddTrack(world.Pos{X: 0, Y: 0}, &types.Track{Direction: types.DirNorth | types.DirSouth | types.DirEast | types.DirWest})

	for x := 35; x < 80; x++ {
		w.AddTrack(world.Pos{X: x, Y: 20}, &types.Track{Direction: types.DirEast | types.DirWest})
		w.AddTrack(world.Pos{X: x, Y: 55}, &types.Track{Direction: types.DirEast | types.DirWest})
	}
	for y := 20; y < 56; y++ {
		w.AddTrack(world.Pos{X: 35, Y: y}, &types.Track{Direction: types.DirNorth | types.DirSouth})
		w.AddTrack(world.Pos{X: 80, Y: y}, &types.Track{Direction: types.DirNorth | types.DirSouth})
	}
	w.AddTrack(world.Pos{X: 35, Y: 20}, &types.Track{Direction: types.DirNorth | types.DirEast})
	w.AddTrack(world.Pos{X: 80, Y: 20}, &types.Track{Direction: types.DirNorth | types.DirWest})
	w.AddTrack(world.Pos{X: 35, Y: 55}, &types.Track{Direction: types.DirSouth | types.DirEast})
	w.AddTrack(world.Pos{X: 80, Y: 55}, &types.Track{Direction: types.DirSouth | types.DirWest})

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
		w.AddTrack(world.Pos{X: x, Y: 40}, &types.Track{Direction: types.DirEast | types.DirWest})
		w.AddTrack(world.Pos{X: x, Y: 60}, &types.Track{Direction: types.DirEast | types.DirWest})
	}
	for y := 40; y < 61; y++ {
		w.AddTrack(world.Pos{X: 65, Y: y}, &types.Track{Direction: types.DirNorth | types.DirSouth})
		w.AddTrack(world.Pos{X: 90, Y: y}, &types.Track{Direction: types.DirNorth | types.DirSouth})
	}
	w.AddTrack(world.Pos{X: 65, Y: 40}, &types.Track{Direction: types.DirNorth | types.DirEast})
	w.AddTrack(world.Pos{X: 90, Y: 40}, &types.Track{Direction: types.DirNorth | types.DirWest})
	w.AddTrack(world.Pos{X: 65, Y: 60}, &types.Track{Direction: types.DirSouth | types.DirEast})
	w.AddTrack(world.Pos{X: 90, Y: 60}, &types.Track{Direction: types.DirSouth | types.DirWest})

	w.AddTrack(world.Pos{X: 80, Y: 40}, &types.Track{Direction: types.DirSouth | types.DirNorth | types.DirEast | types.DirWest})
	w.AddTrack(world.Pos{X: 65, Y: 55}, &types.Track{Direction: types.DirSouth | types.DirNorth | types.DirEast | types.DirWest})

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
