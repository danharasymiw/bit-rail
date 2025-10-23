package message

import (
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
)

type Message struct {
	Type MessageType
	Data []byte
}

type ChatMessage struct {
	Author  string
	Message string
}

type ChunksMessage struct {
	Chunks []*world.Chunk
}

type GetChunksMessage struct {
	Positions []world.Pos
}
type LoginMessage struct {
	Username string
}

type InitialLoadMessage struct {
	Width, Height int
	CameraPos     world.Pos
	Chunks        []*world.Chunk
	Trains        []*trains.Train
	Tracks        map[world.Pos]*types.Track
}
