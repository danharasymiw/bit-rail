package test_worlds

import (
	"github.com/danharasymiw/bit-rail/trains"
	"github.com/danharasymiw/bit-rail/types"
	"github.com/danharasymiw/bit-rail/world"
	"github.com/google/uuid"
)

// RoutingDemoTestWorld builds a single switch (T-junction) with a station
// beside each branch, to manually verify that a train with a Destination
// set picks the correct branch instead of the arbitrary
// first-valid-direction fallback, and that arrival is detected by standing
// next to the station (not by being "on" it — stations aren't on track).
//
//	          ┌──────────┐
//	          │ Station A│ (4,29)-(15,33)
//	          └──────────┘
//	               (10,28) <- adjacent to Station A
//	               (10,27)
//	               (10,26)
//	(0,25)...(9,25)---(10,25)J
//	               (10,24)
//	               (10,23)
//	               (10,22) <- adjacent to Station B; train's Destination
//	          ┌──────────┐
//	          │ Station B│ (4,17)-(15,21)
//	          └──────────┘
//
// The train starts on the west arm heading East. At the junction it must
// choose North (toward Station A) or South (toward Station B); without
// routing the arbitrary fallback would pick North (first direction checked
// that isn't the backtrack West), so heading to Station B only works if
// NextHop is actually consulted. On arrival it stops and reverses back out.
func RoutingDemoTestWorld() *world.World {
	w := world.New(50, 50)

	for x := 0; x < 10; x++ {
		w.AddTrack(types.Pos{X: x, Y: 25}, &types.Track{Direction: types.DirEastWest})
	}

	junctionPos := types.Pos{X: 10, Y: 25}
	w.AddTrack(junctionPos, &types.Track{Direction: types.DirWest | types.DirNorth | types.DirSouth})

	w.AddTrack(types.Pos{X: 10, Y: 26}, &types.Track{Direction: types.DirNorthSouth})
	w.AddTrack(types.Pos{X: 10, Y: 27}, &types.Track{Direction: types.DirNorthSouth})
	stationATrack := types.Pos{X: 10, Y: 28}
	w.AddTrack(stationATrack, &types.Track{Direction: types.DirSouth})
	w.AddStation(&types.Station{Pos: types.Pos{X: 4, Y: 29}, Width: 12, Height: 5, Name: "Station A"})

	w.AddTrack(types.Pos{X: 10, Y: 24}, &types.Track{Direction: types.DirNorthSouth})
	w.AddTrack(types.Pos{X: 10, Y: 23}, &types.Track{Direction: types.DirNorthSouth})
	stationBTrack := types.Pos{X: 10, Y: 22}
	w.AddTrack(stationBTrack, &types.Track{Direction: types.DirNorth})
	w.AddStation(&types.Station{Pos: types.Pos{X: 4, Y: 17}, Width: 12, Height: 5, Name: "Station B"})

	trainCars := []*trains.TrainCar{
		{X: 5, Y: 25, Type: trains.CarTypeLocomotive, Direction: types.DirEast},
		{X: 4, Y: 25, Type: trains.CarTypeCargo, Direction: types.DirEast},
		{X: 3, Y: 25, Type: trains.CarTypeCargo, Direction: types.DirEast},
	}
	w.AddTrain(&trains.Train{
		ID:          uuid.New(),
		IsMoving:    true,
		Destination: stationBTrack,
		Cars:        trainCars,
	})

	return w
}
