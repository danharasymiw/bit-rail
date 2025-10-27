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

//go:generate go run ../cmd/generator/main.go -type=ChatMessage -output=../message/chat_message_gen.go
type ChatMessage struct {
	Author  string
	Message string
}

//go:generate go run ../cmd/generator/main.go -type=ChunksMessage -output=../message/chunks_message_gen.go
type ChunksMessage struct {
	Chunks []*world.Chunk
}

//go:generate go run ../cmd/generator/main.go -type=GetChunksMessage -output=../message/get_chunks_message_gen.go
type GetChunksMessage struct {
	Positions []*world.Pos
}

//go:generate go run ../cmd/generator/main.go -type=LoginMessage -output=../message/login_message_gen.go
type LoginMessage struct {
	Username string
}

//go:generate go run ../cmd/generator/main.go -type=InitialLoadMessage -output=../message/initial_load_message_gen.go
type InitialLoadMessage struct {
	Width, Height uint16
	CameraPos     world.Pos
	Chunks        []*world.Chunk
	Trains        []*trains.Train
	Tracks        []*TrackUpdate
}

//go:generate go run ../cmd/generator/main.go -type=WorldUpdateMessage -output=../message/world_update_message.go
type WorldUpdateMessage struct {
	TilesUpdated  []*TileUpdate
	TracksUpdated []*TrackUpdate
	// TODO: One day don't send entire batch of trains
	Trains []*trains.Train
}

//go:generate go run ../cmd/generator/main.go -type=TileUpdate -output=../message/tile_update_gen.go
type TileUpdate struct {
	Pos  world.Pos
	Tile types.Tile
}

//go:generate go run ../cmd/generator/main.go -type=TrackUpdate -output=../message/track_update_gen.go
type TrackUpdate struct {
	Pos   world.Pos
	Track types.Track
}
