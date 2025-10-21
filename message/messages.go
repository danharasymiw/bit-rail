package message

import (
	"github.com/danharasymiw/bit-rail/trains"
	"github.com/danharasymiw/bit-rail/types"
)

type MessageType uint8

const (
	MessageTypeChat MessageType = iota
	MessageTypeChunks
	MessageTypeInitialLoad
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
	Chunks []Chunk
}

type Chunk struct {
	X, Y  int
	Size  int // size of width and height
	Tiles []*types.Tile
}

type LoginMessage struct {
	Username string
}

type InitialLoadMessage struct {
	Width, Height    int
	CameraX, CameraY int
	Chunks           []Chunk
	Trains           []*trains.Train
}
