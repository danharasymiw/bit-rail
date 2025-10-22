package message

import (
	"github.com/danharasymiw/bit-rail/trains"
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
	Coords []world.ChunkCoord
}
type LoginMessage struct {
	Username string
}

type InitialLoadMessage struct {
	Width, Height    int
	CameraX, CameraY int
	Chunks           []*world.Chunk
	Trains           []*trains.Train
}
