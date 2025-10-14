package main

import (
	"time"

	"github.com/danharasymiw/trains/engine"
	"github.com/danharasymiw/trains/renderer"
	"github.com/danharasymiw/trains/trains"
	"github.com/danharasymiw/trains/types"
	"github.com/danharasymiw/trains/world"

	"github.com/gdamore/tcell"
)

func main() {
	w := world.New(50, 50)

	// Draw a simple world

	// First track
	// Simple horizontal tracks
	for x := 5; x <= 15; x++ {
		w.AddTrack(x, 10, types.DirEast|types.DirWest)
		w.AddTrack(x, 7, types.DirEast|types.DirWest)
	}

	// Vertical connections for the loop
	for y := 7; y <= 10; y++ {
		w.AddTrack(5, y, types.DirNorth|types.DirSouth)
		w.AddTrack(15, y, types.DirNorth|types.DirSouth)
	}

	// Corners
	w.AddTrack(5, 7, types.DirSouth|types.DirEast)
	w.AddTrack(15, 7, types.DirSouth|types.DirWest)
	w.AddTrack(5, 10, types.DirNorth|types.DirEast)
	w.AddTrack(15, 10, types.DirNorth|types.DirWest)

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
		w.AddTrack(x, 12, types.DirEast|types.DirWest)
		w.AddTrack(x, 9, types.DirEast|types.DirWest)
	}

	// Vertical connections for the loop
	for y := 9; y <= 12; y++ {
		w.AddTrack(10, y, types.DirNorth|types.DirSouth)
		w.AddTrack(17, y, types.DirNorth|types.DirSouth)
	}

	// Corners
	w.AddTrack(10, 9, types.DirSouth|types.DirEast)
	w.AddTrack(17, 9, types.DirSouth|types.DirWest)
	w.AddTrack(10, 12, types.DirNorth|types.DirEast)
	w.AddTrack(17, 12, types.DirNorth|types.DirWest)

	// Junctions
	w.AddTrack(10, 10, types.DirNorth|types.DirSouth|types.DirEast|types.DirWest)
	w.AddTrack(15, 9, types.DirNorth|types.DirSouth|types.DirEast|types.DirWest)

	w.Trains = append(w.Trains, &trains.Train{
		IsMoving: true,
		Cars: []*trains.TrainCar{
			{Type: trains.CarTypeLocomotive, X: 13, Y: 12, Direction: types.DirEast},
			{Type: trains.CarTypeCargo, X: 12, Y: 12, Direction: types.DirEast},
			{Type: trains.CarTypeCargo, X: 11, Y: 12, Direction: types.DirEast},
		},
	})

	screen, _ := tcell.NewScreen()
	screen.Init()
	defer screen.Fini()

	eng := engine.New(
		w,
		renderer.NewSimpleRenderer(screen, w),
		150*time.Millisecond,
	)
	eng.Run()
	println("bye.")
}
