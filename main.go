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
	w := world.New(30, 30)

	// Draw a simple world

	// First track
	// Simple horizontal tracks
	for x := 5; x <= 15; x++ {
		w.Tiles[10][x] = types.Tile{Type: types.TileTrack, Orientation: types.DirEast | types.DirWest}
		w.Tiles[7][x] = types.Tile{Type: types.TileTrack, Orientation: types.DirEast | types.DirWest}
	}

	// Vertical connections for the loop
	for y := 7; y <= 10; y++ {
		w.Tiles[y][5] = types.Tile{Type: types.TileTrack, Orientation: types.DirNorth | types.DirSouth}
		w.Tiles[y][15] = types.Tile{Type: types.TileTrack, Orientation: types.DirNorth | types.DirSouth}
	}

	// Corners
	w.Tiles[7][5] = types.Tile{Type: types.TileTrack, Orientation: types.DirSouth | types.DirEast}   // top-left
	w.Tiles[7][15] = types.Tile{Type: types.TileTrack, Orientation: types.DirSouth | types.DirWest}  // top-right
	w.Tiles[10][5] = types.Tile{Type: types.TileTrack, Orientation: types.DirNorth | types.DirEast}  // bottom-left
	w.Tiles[10][15] = types.Tile{Type: types.TileTrack, Orientation: types.DirNorth | types.DirWest} // bottom-right

	w.Trains = append(w.Trains, &trains.Train{
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
		w.Tiles[12][x] = types.Tile{Type: types.TileTrack, Orientation: types.DirEast | types.DirWest}
		w.Tiles[9][x] = types.Tile{Type: types.TileTrack, Orientation: types.DirEast | types.DirWest}
	}

	// Vertical connections for the loop
	for y := 9; y <= 12; y++ {
		w.Tiles[y][10] = types.Tile{Type: types.TileTrack, Orientation: types.DirNorth | types.DirSouth}
		w.Tiles[y][17] = types.Tile{Type: types.TileTrack, Orientation: types.DirNorth | types.DirSouth}
	}

	// Corners
	w.Tiles[9][10] = types.Tile{Type: types.TileTrack, Orientation: types.DirSouth | types.DirEast}  // top-left
	w.Tiles[9][17] = types.Tile{Type: types.TileTrack, Orientation: types.DirSouth | types.DirWest}  // top-right
	w.Tiles[12][10] = types.Tile{Type: types.TileTrack, Orientation: types.DirNorth | types.DirEast} // bottom-left
	w.Tiles[12][17] = types.Tile{Type: types.TileTrack, Orientation: types.DirNorth | types.DirWest} // bottom-right

	// Junctions
	w.Tiles[10][10] = types.Tile{Type: types.TileTrack, Orientation: types.DirNorth | types.DirSouth | types.DirEast | types.DirWest}
	w.Tiles[9][15] = types.Tile{Type: types.TileTrack, Orientation: types.DirNorth | types.DirSouth | types.DirEast | types.DirWest}

	w.Trains = append(w.Trains, &trains.Train{
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
