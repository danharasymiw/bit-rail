package test_worlds

import (
	"github.com/danharasymiw/bit-rail/trains"
	"github.com/danharasymiw/bit-rail/types"
	"github.com/danharasymiw/bit-rail/world"
	"github.com/google/uuid"
)

func IntersectingLoopsTestWorld() *world.World {
	w := world.New(500, 500)

	for x := 300; x < 351; x++ {
		w.AddTrack(types.Pos{X: x, Y: 290}, &types.Track{Direction: types.DirEastWest})
		w.AddTrack(types.Pos{X: x, Y: 270}, &types.Track{Direction: types.DirEastWest})
	}
	for y := 270; y < 291; y++ {
		w.AddTrack(types.Pos{X: 300, Y: y}, &types.Track{Direction: types.DirNorthSouth})
		w.AddTrack(types.Pos{X: 350, Y: y}, &types.Track{Direction: types.DirNorthSouth})
	}
	w.AddTrack(types.Pos{X: 300, Y: 290}, &types.Track{Direction: types.DirSouthEast})
	w.AddTrack(types.Pos{X: 300, Y: 270}, &types.Track{Direction: types.DirNorthEast})
	w.AddTrack(types.Pos{X: 350, Y: 290}, &types.Track{Direction: types.DirSouthWest})
	w.AddTrack(types.Pos{X: 350, Y: 270}, &types.Track{Direction: types.DirNorthWest})

	trainCars := make([]*trains.TrainCar, 0)
	for i := range 30 {
		trainCars = append(trainCars,
			&trains.TrainCar{
				X: 340 - i, Y: 290,
				Type:      trains.CarTypeCargo,
				Direction: types.DirEast,
			},
		)
	}
	trainCars[0].Type = trains.CarTypeLocomotive
	w.AddTrain(&trains.Train{
		ID:       uuid.New(),
		IsMoving: true,
		Cars:     trainCars,
	})

	for x := 325; x < 361; x++ {
		w.AddTrack(types.Pos{X: x, Y: 275}, &types.Track{Direction: types.DirEastWest})
		w.AddTrack(types.Pos{X: x, Y: 265}, &types.Track{Direction: types.DirEastWest})
	}
	for y := 265; y < 276; y++ {
		w.AddTrack(types.Pos{X: 325, Y: y}, &types.Track{Direction: types.DirNorthSouth})
		w.AddTrack(types.Pos{X: 360, Y: y}, &types.Track{Direction: types.DirNorthSouth})
	}
	w.AddTrack(types.Pos{X: 325, Y: 275}, &types.Track{Direction: types.DirSouthEast})
	w.AddTrack(types.Pos{X: 325, Y: 265}, &types.Track{Direction: types.DirNorthEast})
	w.AddTrack(types.Pos{X: 360, Y: 275}, &types.Track{Direction: types.DirSouthWest})
	w.AddTrack(types.Pos{X: 360, Y: 265}, &types.Track{Direction: types.DirNorthWest})
	w.AddTrack(types.Pos{X: 325, Y: 270}, &types.Track{Direction: types.DirFourWay})
	w.AddTrack(types.Pos{X: 350, Y: 275}, &types.Track{Direction: types.DirFourWay})

	trainCars2 := make([]*trains.TrainCar, 0)
	for i := range 20 {
		trainCars2 = append(trainCars2,
			&trains.TrainCar{
				X: 360 - i, Y: 265,
				Type:      trains.CarTypeCargo,
				Direction: types.DirNorth,
			},
		)
	}
	trainCars2[0].Type = trains.CarTypeLocomotive
	w.AddTrain(&trains.Train{
		ID:       uuid.New(),
		IsMoving: true,
		Cars:     trainCars2,
	})

	return w
}
