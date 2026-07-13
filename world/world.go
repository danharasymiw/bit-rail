package world

import (
	"github.com/danharasymiw/bit-rail/trains"
	"github.com/danharasymiw/bit-rail/types"
)

const ChunkSize = 64

// Router answers "which direction should I go" queries for a train sitting
// at a routing node (junction) heading toward dest. Defined here (rather
// than depending on the routing package directly) to avoid an import cycle,
// since routing.Manager itself depends on *World.
type Router interface {
	NextHop(pos types.Pos, dest types.Pos) (types.Dir, bool)
}

type World struct {
	Width, Height int
	Tiles         [][]*types.Tile
	Tracks        map[types.Pos]*types.Track
	Trains        []*trains.Train
	Occupied      map[types.Pos]bool
	Router        Router
	Stations      []*types.Station
}

func New(width, height int) *World {
	w := &World{
		Width:    width,
		Height:   height,
		Tiles:    make([][]*types.Tile, height),
		Tracks:   make(map[types.Pos]*types.Track),
		Trains:   make([]*trains.Train, 0),
		Occupied: make(map[types.Pos]bool),
	}

	for y := range w.Tiles {
		w.Tiles[y] = make([]*types.Tile, width)
		for x := range w.Tiles[y] {
			w.Tiles[y][x] = &types.Tile{Type: types.TileGrass}
		}
	}

	return w
}

func (w *World) TileAt(pos types.Pos) *types.Tile {
	return w.Tiles[pos.Y][pos.X]
}

func (w *World) TrackAt(pos types.Pos) *types.Track {
	return w.Tracks[pos]
}

type Chunk struct {
	Pos   types.Pos
	Tiles []*types.Tile
}

func (w *World) ChunkAt(chunkPos types.Pos) *Chunk {
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

func TileToChunkPos(p types.Pos) types.Pos {
	return types.Pos{X: p.X / ChunkSize, Y: p.Y / ChunkSize}
}

func ChunkToTilePos(chunkPos types.Pos) types.Pos {
	return types.Pos{X: chunkPos.X * ChunkSize, Y: chunkPos.Y * ChunkSize}
}

func (w *World) OccupiedAt(pos types.Pos) bool {
	return w.Occupied[pos]
}

func (w *World) SetOccupied(pos types.Pos) {
	w.Occupied[pos] = true
}

func (w *World) UnsetOccupied(pos types.Pos) {
	w.Occupied[pos] = false
}

func (w *World) AddTrack(pos types.Pos, track *types.Track) {
	tile := &types.Tile{Type: types.TileTrack}
	w.Tiles[pos.Y][pos.X] = tile

	w.Tracks[pos] = track
}

// NextHop returns the direction a train sitting at pos should take to make
// progress toward dest. Delegates to the World's Router (set by the engine).
func (w *World) NextHop(pos, dest types.Pos) (types.Dir, bool) {
	return w.Router.NextHop(pos, dest)
}

func (w *World) AddTrain(t *trains.Train) {
	w.Trains = append(w.Trains, t)
	for _, c := range t.Cars {
		w.SetOccupied(types.Pos{X: c.X, Y: c.Y})
	}
}

func (w *World) AddStation(s *types.Station) {
	w.Stations = append(w.Stations, s)
}

// StationAt returns the station (if any) that pos is adjacent to.
func (w *World) StationAt(pos types.Pos) *types.Station {
	return types.StationAt(pos, w.Stations)
}
