package message

import (
	"io"

	"github.com/danharasymiw/bit-rail/trains"
	"github.com/danharasymiw/bit-rail/types"
	"github.com/danharasymiw/bit-rail/world"
)

type MessageType uint8

const (
	MessageTypeChat MessageType = iota
	MessageTypeChunks
	MessageTypeInitialLoad
	MessageTypeLogin
	MessageTypeGetChunks
	MessageTypeWorldUpdate
)

type OutgoingMessage struct {
	Type MessageType
	Data io.Writer
}

type IncomingMessage struct {
	Type MessageType
	Data io.Reader
}

type ChatMessage struct {
	Author  string
	Message string
}

type ChunksMessage struct {
	Chunks []*world.Chunk
}

type GetChunksMessage struct {
	Positions []*world.Pos
}

type LoginMessage struct {
	Username string
}

type InitialLoadMessage struct {
	Width, Height uint16
	CameraPos     world.Pos
	Chunks        []*world.Chunk
	Trains        []*trains.Train
	Tracks        []*TrackUpdate
}

type WorldUpdateMessage struct {
	TilesUpdated  []*TileUpdate
	TracksUpdated []*TrackUpdate
	// TODO: One day don't send entire batch of trains
	Trains []*trains.Train
}

type TileUpdate struct {
	Pos  world.Pos
	Tile types.Tile
}

type TrackUpdate struct {
	Pos   world.Pos
	Track types.Track
}
