package world

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/danharasymiw/bit-rail/trains"
	"github.com/danharasymiw/bit-rail/types"
)

func TestNew(t *testing.T) {
	w := New(128, 64)

	assert.Equal(t, 128, w.Width)
	assert.Equal(t, 64, w.Height)

	// All tiles should be TileGrass
	for y := 0; y < w.Height; y++ {
		for x := 0; x < w.Width; x++ {
			tile := w.Tiles[y][x]
			require.NotNil(t, tile, "tile at (%d,%d) should not be nil", x, y)
			assert.Equal(t, types.TileGrass, tile.Type, "tile at (%d,%d) should be TileGrass", x, y)
		}
	}

	assert.Empty(t, w.Tracks, "Tracks should be empty")
	assert.Empty(t, w.Occupied, "Occupied should be empty")
}

func TestTileAt(t *testing.T) {
	w := New(10, 10)
	pos := types.Pos{X: 3, Y: 5}

	tile := w.TileAt(pos)
	require.NotNil(t, tile)
	assert.Equal(t, types.TileGrass, tile.Type)
}

func TestTrackAt_Missing(t *testing.T) {
	w := New(10, 10)
	pos := types.Pos{X: 2, Y: 2}

	track := w.TrackAt(pos)
	assert.Nil(t, track)
}

func TestTrackAt_Present(t *testing.T) {
	w := New(10, 10)
	pos := types.Pos{X: 4, Y: 4}
	track := &types.Track{Direction: types.DirNorthSouth}

	w.AddTrack(pos, track)

	result := w.TrackAt(pos)
	require.NotNil(t, result)
	assert.Equal(t, track, result)
	assert.Equal(t, types.Dir(types.DirNorthSouth), result.Direction)
}

func TestSetAndUnsetOccupied(t *testing.T) {
	w := New(10, 10)
	pos := types.Pos{X: 1, Y: 1}

	assert.False(t, w.OccupiedAt(pos), "position should not be occupied initially")

	w.SetOccupied(pos)
	assert.True(t, w.OccupiedAt(pos), "position should be occupied after SetOccupied")

	w.UnsetOccupied(pos)
	assert.False(t, w.OccupiedAt(pos), "position should not be occupied after UnsetOccupied")
}

func TestAddTrack(t *testing.T) {
	w := New(10, 10)
	pos := types.Pos{X: 5, Y: 3}
	track := &types.Track{Direction: types.DirEastWest}

	w.AddTrack(pos, track)

	// Tile type should become TileTrack
	tile := w.TileAt(pos)
	require.NotNil(t, tile)
	assert.Equal(t, types.TileTrack, tile.Type)

	// TrackAt should return the track
	result := w.TrackAt(pos)
	require.NotNil(t, result)
	assert.Equal(t, track, result)
}

func TestAddTrain(t *testing.T) {
	w := New(20, 20)

	car1 := &trains.TrainCar{X: 5, Y: 5, Direction: types.DirEast, Type: trains.CarTypeLocomotive}
	car2 := &trains.TrainCar{X: 4, Y: 5, Direction: types.DirEast, Type: trains.CarTypeCargo}

	train := &trains.Train{
		ID:       uuid.New(),
		IsMoving: false,
		Cars:     []*trains.TrainCar{car1, car2},
	}

	w.AddTrain(train)

	// Train should appear in w.Trains
	require.Len(t, w.Trains, 1)
	assert.Equal(t, train, w.Trains[0])

	// All car positions should be in w.Occupied
	assert.True(t, w.OccupiedAt(types.Pos{X: 5, Y: 5}), "car1 position should be occupied")
	assert.True(t, w.OccupiedAt(types.Pos{X: 4, Y: 5}), "car2 position should be occupied")
}

func TestChunkAt(t *testing.T) {
	w := New(ChunkSize*4, ChunkSize*4)
	chunkPos := types.Pos{X: 0, Y: 0}

	chunk := w.ChunkAt(chunkPos)
	require.NotNil(t, chunk)
	assert.Len(t, chunk.Tiles, ChunkSize*ChunkSize)

	// Tiles in chunk should match world tiles
	for idx, tile := range chunk.Tiles {
		x := idx % ChunkSize
		y := idx / ChunkSize
		worldTile := w.Tiles[y][x]
		assert.Equal(t, worldTile, tile, "chunk tile at (%d,%d) should match world tile", x, y)
	}
}

func TestTileToChunkPos(t *testing.T) {
	tests := []struct {
		tilePos  types.Pos
		expected types.Pos
	}{
		{types.Pos{X: 0, Y: 0}, types.Pos{X: 0, Y: 0}},
		{types.Pos{X: 64, Y: 128}, types.Pos{X: 1, Y: 2}},
		{types.Pos{X: 63, Y: 63}, types.Pos{X: 0, Y: 0}},
		{types.Pos{X: 127, Y: 64}, types.Pos{X: 1, Y: 1}},
	}

	for _, tt := range tests {
		result := TileToChunkPos(tt.tilePos)
		assert.Equal(t, tt.expected, result, "TileToChunkPos(%v)", tt.tilePos)
	}
}

func TestChunkToTilePos(t *testing.T) {
	tests := []struct {
		chunkPos types.Pos
		expected types.Pos
	}{
		{types.Pos{X: 0, Y: 0}, types.Pos{X: 0, Y: 0}},
		{types.Pos{X: 1, Y: 2}, types.Pos{X: 64, Y: 128}},
		{types.Pos{X: 2, Y: 3}, types.Pos{X: 128, Y: 192}},
	}

	for _, tt := range tests {
		result := ChunkToTilePos(tt.chunkPos)
		assert.Equal(t, tt.expected, result, "ChunkToTilePos(%v)", tt.chunkPos)
	}
}
