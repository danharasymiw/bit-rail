package message

import (
	"bytes"
	"strings"
	"testing"

	"github.com/danharasymiw/bit-rail/trains"
	"github.com/danharasymiw/bit-rail/types"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestChatMessageRoundTrip(t *testing.T) {
	expected := &ChatMessage{
		Author:  "dan",
		Message: "hello world",
	}

	var buf bytes.Buffer
	require.NoError(t, WriteChatMessage(&buf, expected))

	got, err := ReadChatMessage(&buf)
	require.NoError(t, err)
	assert.Equal(t, expected.Author, got.Author)
	assert.Equal(t, expected.Message, got.Message)
}

func TestChatMessage_AuthorTooLong(t *testing.T) {
	msg := &ChatMessage{
		Author:  strings.Repeat("A", 256),
		Message: "hi",
	}
	var buf bytes.Buffer
	err := WriteChatMessage(&buf, msg)
	assert.Error(t, err)
}

func TestChatMessage_MessageTooLong(t *testing.T) {
	msg := &ChatMessage{
		Author:  "dan",
		Message: strings.Repeat("A", 256),
	}
	var buf bytes.Buffer
	err := WriteChatMessage(&buf, msg)
	assert.Error(t, err)
}

func TestLoginMessageRoundTrip(t *testing.T) {
	expected := &LoginMessage{Username: "testuser"}

	var buf bytes.Buffer
	require.NoError(t, WriteLoginMessage(&buf, expected))

	got, err := ReadLoginMessage(&buf)
	require.NoError(t, err)
	assert.Equal(t, expected.Username, got.Username)
}

func TestTrackMapToTrackUpdate(t *testing.T) {
	trackMap := map[types.Pos]*types.Track{
		{X: 1, Y: 2}: {Direction: types.DirNorthSouth},
		{X: 3, Y: 4}: {Direction: types.DirEastWest},
	}

	updates := TrackMapToTrackUpdate(trackMap)

	require.Len(t, updates, 2)

	// Build a map from pos -> direction to verify both entries, regardless of order.
	byPos := make(map[types.Pos]types.Dir, len(updates))
	for _, u := range updates {
		byPos[u.Pos] = u.Track.Direction
	}
	assert.Equal(t, types.Dir(types.DirNorthSouth), byPos[types.Pos{X: 1, Y: 2}])
	assert.Equal(t, types.Dir(types.DirEastWest), byPos[types.Pos{X: 3, Y: 4}])
}

func TestPosRoundTrip(t *testing.T) {
	expected := types.Pos{X: 100, Y: 200}

	var buf bytes.Buffer
	require.NoError(t, WritePos(&buf, &expected))

	got, err := ReadPos(&buf)
	require.NoError(t, err)
	assert.Equal(t, expected.X, got.X)
	assert.Equal(t, expected.Y, got.Y)
}

func TestTrackUpdateRoundTrip(t *testing.T) {
	expected := &TrackUpdate{
		Pos:   types.Pos{X: 5, Y: 10},
		Track: types.Track{Direction: types.DirEastWest},
	}

	var buf bytes.Buffer
	require.NoError(t, WriteTrackUpdate(&buf, expected))

	got, err := ReadTrackUpdate(&buf)
	require.NoError(t, err)
	assert.Equal(t, expected.Pos.X, got.Pos.X)
	assert.Equal(t, expected.Pos.Y, got.Pos.Y)
	assert.Equal(t, expected.Track.Direction, got.Track.Direction)
}

func TestTrainRoundTrip(t *testing.T) {
	expected := &trains.Train{
		ID:          uuid.New(),
		IsMoving:    true,
		IsReversing: false,
		Cars: []*trains.TrainCar{
			{X: 10, Y: 20, Direction: types.DirEast, Type: trains.CarTypeLocomotive},
			{X: 9, Y: 20, Direction: types.DirEast, Type: trains.CarTypeCargo},
		},
	}

	var buf bytes.Buffer
	require.NoError(t, WriteTrain(&buf, expected))

	got, err := ReadTrain(&buf)
	require.NoError(t, err)
	assert.Equal(t, expected.ID, got.ID)
	assert.Equal(t, expected.IsMoving, got.IsMoving)
	assert.Equal(t, expected.IsReversing, got.IsReversing)
	require.Len(t, got.Cars, 2)
	assert.Equal(t, expected.Cars[0].X, got.Cars[0].X)
	assert.Equal(t, expected.Cars[0].Y, got.Cars[0].Y)
	assert.Equal(t, expected.Cars[0].Type, got.Cars[0].Type)
	assert.Equal(t, expected.Cars[1].X, got.Cars[1].X)
	assert.Equal(t, expected.Cars[1].Y, got.Cars[1].Y)
	assert.Equal(t, expected.Cars[1].Type, got.Cars[1].Type)
}

func TestWorldUpdateMessageRoundTrip(t *testing.T) {
	expected := &WorldUpdateMessage{
		TilesUpdated: []*TileUpdate{
			{Pos: types.Pos{X: 1, Y: 2}, Tile: types.Tile{Type: types.TileGrass}},
		},
		TracksUpdated: []*TrackUpdate{
			{Pos: types.Pos{X: 3, Y: 4}, Track: types.Track{Direction: types.DirNorthSouth}},
		},
		Trains: []*trains.Train{
			{
				ID:       uuid.New(),
				IsMoving: true,
				Cars: []*trains.TrainCar{
					{X: 5, Y: 6, Direction: types.DirNorth, Type: trains.CarTypeLocomotive},
				},
			},
		},
	}

	var buf bytes.Buffer
	require.NoError(t, WriteWorldUpdateMessage(&buf, expected))

	got, err := ReadWorldUpdateMessage(&buf)
	require.NoError(t, err)
	require.Len(t, got.TilesUpdated, 1)
	assert.Equal(t, expected.TilesUpdated[0].Pos, got.TilesUpdated[0].Pos)
	assert.Equal(t, expected.TilesUpdated[0].Tile.Type, got.TilesUpdated[0].Tile.Type)
	require.Len(t, got.TracksUpdated, 1)
	assert.Equal(t, expected.TracksUpdated[0].Pos, got.TracksUpdated[0].Pos)
	assert.Equal(t, expected.TracksUpdated[0].Track.Direction, got.TracksUpdated[0].Track.Direction)
	require.Len(t, got.Trains, 1)
	assert.Equal(t, expected.Trains[0].ID, got.Trains[0].ID)
	assert.Equal(t, expected.Trains[0].IsMoving, got.Trains[0].IsMoving)
}

func TestGetChunksMessageRoundTrip(t *testing.T) {
	expected := &GetChunksMessage{
		Positions: []types.Pos{
			{X: 0, Y: 0},
			{X: 64, Y: 64},
			{X: 128, Y: 192},
		},
	}

	var buf bytes.Buffer
	require.NoError(t, WriteGetChunksMessage(&buf, expected))

	got, err := ReadGetChunksMessage(&buf)
	require.NoError(t, err)
	require.Len(t, got.Positions, 3)
	for i, pos := range expected.Positions {
		assert.Equal(t, pos.X, got.Positions[i].X)
		assert.Equal(t, pos.Y, got.Positions[i].Y)
	}
}
