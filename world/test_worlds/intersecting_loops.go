package test_worlds

import (
	"github.com/danharasymiw/bit-rail/trains"
	"github.com/danharasymiw/bit-rail/types"
	"github.com/danharasymiw/bit-rail/world"
)

func IntersectingLoopsTestWorld() *world.World {
	w := world.New(50, 50)

	// Draw a simple world

	// First track
	// Simple horizontal tracks
	for x := 5; x <= 15; x++ {
		w.AddTrack(world.Pos{X: x, Y: 10}, &types.Track{Direction: types.DirEast | types.DirWest})
		w.AddTrack(world.Pos{X: x, Y: 7}, &types.Track{Direction: types.DirEast | types.DirWest})
	}

	// Vertical connections for the loop
	for y := 7; y <= 10; y++ {
		w.AddTrack(world.Pos{X: 5, Y: y}, &types.Track{Direction: types.DirNorth | types.DirSouth})
		w.AddTrack(world.Pos{X: 15, Y: y}, &types.Track{Direction: types.DirNorth | types.DirSouth})
	}

	// Corners
	w.AddTrack(world.Pos{X: 5, Y: 7}, &types.Track{Direction: types.DirSouth | types.DirEast})
	w.AddTrack(world.Pos{X: 15, Y: 7}, &types.Track{Direction: types.DirSouth | types.DirWest})
	w.AddTrack(world.Pos{X: 5, Y: 10}, &types.Track{Direction: types.DirNorth | types.DirEast})
	w.AddTrack(world.Pos{X: 15, Y: 10}, &types.Track{Direction: types.DirNorth | types.DirWest})

	w.AddTrain(&trains.Train{
		IsMoving: true,
		Cars: []*trains.TrainCar{
			{Type: trains.CarTypeLocomotive, X: 9, Y: 10, Direction: types.DirWest},
			{Type: trains.CarTypeCargo, X: 10, Y: 10, Direction: types.DirWest},
			{Type: trains.CarTypeCargo, X: 11, Y: 10, Direction: types.DirWest},
			{Type: trains.CarTypeCargo, X: 12, Y: 10, Direction: types.DirWest},
			{Type: trains.CarTypeCargo, X: 13, Y: 10, Direction: types.DirWest},
			{Type: trains.CarTypeCargo, X: 14, Y: 10, Direction: types.DirWest},
		},
	})

	// Second track
	for x := 10; x <= 17; x++ {
		w.AddTrack(world.Pos{X: x, Y: 12}, &types.Track{Direction: types.DirEast | types.DirWest})
		w.AddTrack(world.Pos{X: x, Y: 9}, &types.Track{Direction: types.DirEast | types.DirWest})
	}

	// Vertical connections for the loop
	for y := 9; y <= 12; y++ {
		w.AddTrack(world.Pos{X: 10, Y: y}, &types.Track{Direction: types.DirNorth | types.DirSouth})
		w.AddTrack(world.Pos{X: 17, Y: y}, &types.Track{Direction: types.DirNorth | types.DirSouth})
	}

	// Corners
	w.AddTrack(world.Pos{X: 10, Y: 9}, &types.Track{Direction: types.DirSouth | types.DirEast})
	w.AddTrack(world.Pos{X: 17, Y: 9}, &types.Track{Direction: types.DirSouth | types.DirWest})
	w.AddTrack(world.Pos{X: 10, Y: 12}, &types.Track{Direction: types.DirNorth | types.DirEast})
	w.AddTrack(world.Pos{X: 17, Y: 12}, &types.Track{Direction: types.DirNorth | types.DirWest})

	// Junctions
	w.AddTrack(world.Pos{X: 10, Y: 10}, &types.Track{Direction: types.DirNorth | types.DirSouth | types.DirEast | types.DirWest})
	w.AddTrack(world.Pos{X: 15, Y: 9}, &types.Track{Direction: types.DirNorth | types.DirSouth | types.DirEast | types.DirWest})

	w.Trains = append(w.Trains, &trains.Train{
		IsMoving: true,
		Cars: []*trains.TrainCar{
			{Type: trains.CarTypeLocomotive, X: 13, Y: 12, Direction: types.DirEast},
			{Type: trains.CarTypeCargo, X: 12, Y: 12, Direction: types.DirEast},
			{Type: trains.CarTypeCargo, X: 11, Y: 12, Direction: types.DirEast},
		},
	})

	// water
	waterCoords := []struct{ x, y int }{
		{35, 35},
		{36, 35},
		{37, 35},
		{38, 35},
		{39, 35},
		{40, 35},
		{41, 35},
		{35, 36},
		{36, 36},
		{37, 36},
		{38, 36},
		{39, 36},
		{40, 36},
		{35, 37},
		{36, 37},
		{37, 37},
		{38, 37},
		{36, 38},
		{37, 38},
		{38, 38},
		{39, 38},
		{37, 39},
		{38, 39},
	}

	for _, coord := range waterCoords {
		w.Tiles[coord.y][coord.x] = &types.Tile{Type: types.TileWater}
	}

	return w
}
