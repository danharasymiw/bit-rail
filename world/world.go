package world

import (
	"github.com/danharasymiw/trains/trains"
	"github.com/danharasymiw/trains/types"
)

type World struct {
	Width, Height int
	Tiles         [][]*types.Tile
	Tracks        map[*types.Tile]types.Track
	Trains        []*trains.Train
	Occupied      map[int]bool
}

func New(width, height int) *World {
	w := &World{
		Width:    width,
		Height:   height,
		Tiles:    make([][]*types.Tile, 30),
		Tracks:   make(map[*types.Tile]types.Track),
		Trains:   make([]*trains.Train, 0),
		Occupied: make(map[int]bool),
	}

	for y := range w.Tiles {
		w.Tiles[y] = make([]*types.Tile, 30)
		for x := range w.Tiles[y] {
			w.Tiles[y][x] = &types.Tile{Type: types.TileGrass}
		}
	}

	return w
}

// TileAt exists incase we decide to switch to a 1D array for the world
func (w *World) TileAt(x, y int) *types.Tile {
	return w.Tiles[y][x]
}

func (w *World) OccupiedAt(x, y int) bool {
	return w.Occupied[w.occupiedIndex(x, y)]
}

func (w *World) SetOccupied(x, y int) {
	w.Occupied[w.occupiedIndex(x, y)] = true
}

func (w *World) UnsetOccupied(x, y int) {
	w.Occupied[w.occupiedIndex(x, y)] = false
}

func (w *World) occupiedIndex(x, y int) int {
	return y*w.Width + x
}

func (w *World) AddTrack(x, y int, dir types.Dir) *types.Track {
	tile := &types.Tile{Type: types.TileTrack}
	w.Tiles[y][x] = tile

	track := &types.Track{
		Tile:      tile,
		Direction: dir,
	}
	w.Tracks[tile] = *track

	return track
}

func (w *World) AddTrain(t *trains.Train) {
	w.Trains = append(w.Trains, t)
	for _, c := range t.Cars {
		w.SetOccupied(c.X, c.Y)
	}
}
