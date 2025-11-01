package world

import (
	"github.com/danharasymiw/bit-rail/trains"
	"github.com/danharasymiw/bit-rail/types"
)

const ChunkSize = 64

type World struct {
	Width, Height int
	Tiles         [][]*types.Tile
	Tracks        map[Pos]*types.Track
	Trains        []*trains.Train
	Occupied      map[Pos]bool
}

func New(width, height int) *World {
	w := &World{
		Width:    width,
		Height:   height,
		Tiles:    make([][]*types.Tile, height),
		Tracks:   make(map[Pos]*types.Track),
		Trains:   make([]*trains.Train, 0),
		Occupied: make(map[Pos]bool),
	}

	for y := range w.Tiles {
		w.Tiles[y] = make([]*types.Tile, width)
		for x := range w.Tiles[y] {
			w.Tiles[y][x] = &types.Tile{Type: types.TileGrass}
		}
	}

	return w
}

// TileAt exists incase we decide to switch to a 1D array for the world
func (w *World) TileAt(pos Pos) *types.Tile {
	return w.Tiles[pos.Y][pos.X]
}

// TODO: Maybe we move this to a different package. Feels bad
// having custom marshal logic in this package
type Pos struct {
	X int `binary:"uint16"`
	Y int `binary:"uint16"`
}

type Chunk struct {
	Pos   Pos
	Tiles []*types.Tile
}

func (w *World) ChunkAt(chunkPos Pos) *Chunk {
	tiles := make([]*types.Tile, 0, ChunkSize*ChunkSize)

	for y := chunkPos.Y * ChunkSize; y < (chunkPos.Y+1)*ChunkSize; y++ {
		for x := chunkPos.X * ChunkSize; x < (chunkPos.X+1)*ChunkSize; x++ {
			// Bounds check to prevent index out of range
			if x < w.Width && y < w.Height {
				tile := w.Tiles[y][x]
				tiles = append(tiles, tile)

			} else {
				// For out-of-bounds tiles, create a default grass tile
				tiles = append(tiles, &types.Tile{Type: types.TileGrass})
			}
		}
	}
	return &Chunk{
		Pos:   chunkPos,
		Tiles: tiles,
	}
}

func TileToChunkPos(pos Pos) Pos {
	return Pos{X: pos.X / ChunkSize, Y: pos.Y / ChunkSize}
}

func ChunkToTilePos(chunkPos Pos) Pos {
	return Pos{X: chunkPos.X * ChunkSize, Y: chunkPos.Y * ChunkSize}
}

func (w *World) OccupiedAt(pos Pos) bool {
	return w.Occupied[pos]
}

func (w *World) SetOccupied(pos Pos) {
	w.Occupied[pos] = true
}

func (w *World) UnsetOccupied(pos Pos) {
	w.Occupied[pos] = false
}

func (w *World) AddTrack(pos Pos, track *types.Track) {
	tile := &types.Tile{Type: types.TileTrack}
	w.Tiles[pos.Y][pos.X] = tile

	w.Tracks[pos] = track
}

func (w *World) AddTrain(t *trains.Train) {
	w.Trains = append(w.Trains, t)
	for _, c := range t.Cars {
		w.SetOccupied(Pos{X: c.X, Y: c.Y})
	}
}
