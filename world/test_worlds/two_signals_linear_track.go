package test_worlds

import (
	"github.com/danharasymiw/bit-rail/trains"
	"github.com/danharasymiw/bit-rail/types"
	"github.com/danharasymiw/bit-rail/world"
	"github.com/google/uuid"
)

func NewTwoSignalsLinearTrack() *world.World {
	w := world.New(500, 500)
	for x := 300; x <= 350; x++ {
		w.AddTrack(types.Pos{X: x, Y: 275}, &types.Track{Direction: types.DirEastWest})
	}

	w.TrackAt(types.Pos{X: 315, Y: 275}).SignalDir = types.DirEast
	w.TrackAt(types.Pos{X: 330, Y: 275}).SignalDir = types.DirWest

	w.AddTrain(&trains.Train{
		ID:       uuid.New(),
		IsMoving: true,
		Cars:     []*trains.TrainCar{{Type: trains.CarTypeLocomotive, X: 300, Y: 275, Direction: types.DirEast}},
	})
	w.AddTrain(&trains.Train{
		ID:       uuid.New(),
		IsMoving: true,
		Cars:     []*trains.TrainCar{{Type: trains.CarTypeLocomotive, X: 350, Y: 275, Direction: types.DirWest}},
	})

	return w
}
