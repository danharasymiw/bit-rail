package engine

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/danharasymiw/bit-rail/message"
	"github.com/danharasymiw/bit-rail/types"
	"github.com/danharasymiw/bit-rail/world"
	"github.com/sirupsen/logrus"
)

const systemPlayerName = "SYSTEM"

type Engine struct {
	w       *world.World
	tickDur time.Duration
	running bool
	nm      *networkManager
	bm      *blockManager

	// TODO: do we really want this to be a message type?
	tilesUpdated  []*message.TileUpdate
	tracksUpdated []*message.TrackUpdate
}

func New(w *world.World, tickDur time.Duration, addr string) *Engine {
	eng := &Engine{
		w:       w,
		tickDur: tickDur,
		nm:      newNetworkManager(addr),
		bm:      newBlockManager(w),
	}

	return eng
}

func (e *Engine) Run(quitCh <-chan struct{}, readyCh chan<- struct{}) {
	// Calculate initial blocks
	// All tracks should already belong to a block, but make sure they do on startup
	// Also allows us to test block manager in test worlds
	e.bm.RebuildAll()

	go e.nm.startServer(readyCh)

	ticker := time.NewTicker(e.tickDur)
	defer ticker.Stop()

	e.running = true
	for e.running {
		select {
		case incoming := <-e.nm.incomingCh:
			e.handlePlayerMessage(incoming)
		case <-ticker.C:
			e.tick()
		case <-quitCh:
			e.running = false
		}
	}

	// Give goroutines time to clean up
	time.Sleep(100 * time.Millisecond)
}

func (e *Engine) tick() {
	e.bm.ProcessDirty()

	prevOccupancy := make(map[types.Pos]uuid.UUID)
	for pos, track := range e.w.Tracks {
		if track.HasSignal() && track.Block != nil {
			prevOccupancy[pos] = track.Block.OccupiedBy
		}
	}

	for _, t := range e.w.Trains {
		t.Tick(e.w)
	}

	for pos, prev := range prevOccupancy {
		track := e.w.Tracks[pos]
		if track.Block != nil && track.Block.OccupiedBy != prev {
			e.tracksUpdated = append(e.tracksUpdated, &message.TrackUpdate{Pos: pos, Track: *track})
		}
	}

	e.nm.broadcastCh <- &outgoingMessage{
		worldUpdateMessage: &message.WorldUpdateMessage{
			TilesUpdated:  e.tilesUpdated,
			TracksUpdated: e.tracksUpdated,
			Trains:        e.w.Trains,
		},
	}

	e.tilesUpdated = make([]*message.TileUpdate, 0)
	e.tracksUpdated = make([]*message.TrackUpdate, 0)
}

func (e *Engine) getChunksInRegion(worldPos types.Pos) []*world.Chunk {
	chunks := make([]*world.Chunk, 0)

	centerChunk := world.TileToChunkPos(worldPos)

	// Get 3x3 grid of chunks around the center
	for i := -3; i <= 3; i++ {
		for j := -3; j <= 3; j++ {
			x := centerChunk.X + i
			y := centerChunk.Y + j
			if x < 0 || y < 0 {
				continue
			}
			chunkPos := types.Pos{X: x, Y: y}

			chunk := e.w.ChunkAt(chunkPos)
			chunks = append(chunks, chunk)
		}
	}
	return chunks
}

func (e *Engine) handlePlayerMessage(playerMsg *playerMessage) {
	msg := playerMsg.message
	switch {
	case msg.chatMessage != nil:
		e.handleChatMessage(playerMsg)
	case msg.loginMessage != nil:
		e.handleLoginMessage(playerMsg)
	case msg.getChunksMessage != nil:
		e.handleGetChunksMessage(playerMsg)
	}
}

func (e *Engine) handleChatMessage(playerMsg *playerMessage) {
	entry := logrus.WithField("player", playerMsg.playerID).WithField("message", playerMsg.message.chatMessage.Message)
	e.nm.broadcastCh <- &outgoingMessage{chatMessage: playerMsg.message.chatMessage}
	entry.Debug("Player sent chat message")
}

func (e *Engine) handleLoginMessage(playerMsg *playerMessage) {
	entry := logrus.WithField("player", playerMsg.playerID).WithField("message", playerMsg.message.loginMessage.Username)
	entry.Debug("Processing login message")

	camPos := types.Pos{X: e.w.Width / 2, Y: e.w.Height / 2}

	initialLoadMessage := message.InitialLoadMessage{
		Width:     e.w.Width,
		Height:    e.w.Height,
		CameraPos: types.Pos{X: camPos.X, Y: camPos.Y},
		Chunks:    e.getChunksInRegion(camPos),
		Trains:    e.w.Trains,                                // TODO: get trains in region
		Tracks:    message.TrackMapToTrackUpdate(e.w.Tracks), // TODO: get tracks in region
	}
	entry.Debugf("Sending initial load message to response channel (chunks: %d, trains: %d, tracks: %d)",
		len(initialLoadMessage.Chunks), len(initialLoadMessage.Trains), len(initialLoadMessage.Tracks))
	playerMsg.responseCh <- &outgoingMessage{initialLoadMessage: &initialLoadMessage}
	e.nm.broadcastCh <- &outgoingMessage{
		chatMessage: &message.ChatMessage{
			Author:  systemPlayerName,
			Message: fmt.Sprintf("%s connected.", playerMsg.playerID),
		},
	}
}

func (e *Engine) handleGetChunksMessage(playerMsg *playerMessage) {
	entry := logrus.WithField("player", playerMsg.playerID).WithField("message", playerMsg.message.getChunksMessage)

	chunks := make([]*world.Chunk, 0, len(playerMsg.message.getChunksMessage.Positions))
	for _, pos := range playerMsg.message.getChunksMessage.Positions {
		chunkStartX := pos.X * world.ChunkSize
		chunkStartY := pos.Y * world.ChunkSize

		if pos.X < 0 || pos.Y < 0 || chunkStartX >= e.w.Width || chunkStartY >= e.w.Height {
			continue
		}
		chunks = append(chunks, e.w.ChunkAt(pos))
	}
	playerMsg.responseCh <- &outgoingMessage{chunksMessage: &message.ChunksMessage{Chunks: chunks}}
	entry.Debugf("Player requested chunks")
}

// BroadcastChatMessage sends a chat message to all connected players
func (e *Engine) BroadcastChatMessage(author, msg string) {
	e.nm.broadcastCh <- &outgoingMessage{
		chatMessage: &message.ChatMessage{
			Author:  author,
			Message: msg,
		},
	}
}

func (e *Engine) AddTrack(p types.Pos, t *types.Track) {
	e.w.AddTrack(p, t)
	e.bm.MarkDirty(p)
	e.tracksUpdated = append(e.tracksUpdated, &message.TrackUpdate{Pos: p, Track: *t})
}
