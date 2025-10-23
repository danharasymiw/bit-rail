package engine

import (
	"time"

	"github.com/danharasymiw/bit-rail/message"
	"github.com/danharasymiw/bit-rail/trains"
	"github.com/danharasymiw/bit-rail/types"
	"github.com/danharasymiw/bit-rail/world"
	"github.com/sirupsen/logrus"
)

type Engine struct {
	w       *world.World
	tickDur time.Duration
	running bool
	nm      *networkManager
}

func New(w *world.World, tickDur time.Duration) *Engine {
	eng := &Engine{
		w:       w,
		tickDur: tickDur,
	}
	eng.nm = newNetworkManager()
	return eng
}

func (e *Engine) Run(quitCh <-chan struct{}, readyCh chan<- struct{}) {
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
	for _, t := range e.w.Trains {
		e.moveTrain(t)
	}
}

func (e *Engine) moveTrain(t *trains.Train) {
	// TODO investigate if this function makes more sense to turn/figure out direction then move
	// Currently we move, and then figure out out next direction
	if !t.IsMoving {
		return
	}

	car := t.Cars[0]
	moveDir := car.Direction
	if t.IsReversing {
		car = t.Cars[len(t.Cars)-1]
		moveDir = types.OppositeDir(car.Direction)
	}

	pos := world.Pos{X: car.X, Y: car.Y}
	dir := car.Direction
	nextPos := nextPos(pos, dir)
	nextTile := e.w.TileAt(nextPos)
	if nextTile.Type != types.TileTrack {
		return
	}

	if e.w.OccupiedAt(nextPos) {
		return
	}

	e.moveCars(t.Cars, moveDir, t.IsReversing)

	car = t.Cars[0]
	if t.IsReversing {
		car = t.Cars[len(t.Cars)-1]
		car.Direction = types.OppositeDir(car.Direction)
	}
	pos = world.Pos{X: car.X, Y: car.Y}
	dir = car.Direction
	track := e.w.Tracks[pos]
	if track == nil {
		return
	}

	incFrom := types.OppositeDir(dir)
	if track.Direction&incFrom == 0 {
		return
	}

	outgoing := track.Direction & ^incFrom

	if outgoing != 0 && (outgoing&(outgoing-1)) == 0 {
		car.Direction = outgoing & -outgoing
		return
	}

	if outgoing&dir != 0 {
		return
	}

	for d := types.DirNorth; d <= types.DirWest; d <<= 1 {
		if outgoing&types.Dir(d) != 0 {
			car.Direction = types.Dir(d)
			return
		}
	}
}

func (e *Engine) moveCars(cars []*trains.TrainCar, moveDir types.Dir, reverse bool) {
	start, end, step := 0, len(cars), 1
	if reverse {
		start, end, step = len(cars)-1, -1, -1
	}

	car := cars[start]

	newPos := nextPos(world.Pos{X: car.X, Y: car.Y}, moveDir)

	prevPos := world.Pos{X: car.X, Y: car.Y}
	prevDir := moveDir
	car.X, car.Y = newPos.X, newPos.Y
	e.w.SetOccupied(world.Pos{X: car.X, Y: car.Y})

	for i := start + step; i != end; i += step {
		car = cars[i]
		thisPrevPos := world.Pos{X: car.X, Y: car.Y}
		thisPrevDir := car.Direction

		car.X, car.Y, car.Direction = prevPos.X, prevPos.Y, prevDir

		prevPos, prevDir = thisPrevPos, thisPrevDir
	}
	e.w.UnsetOccupied(prevPos)
}

func nextPos(pos world.Pos, dir types.Dir) world.Pos {
	switch dir {
	case types.DirNorth:
		return world.Pos{X: pos.X, Y: pos.Y + 1}
	case types.DirSouth:
		return world.Pos{X: pos.X, Y: pos.Y - 1}
	case types.DirEast:
		return world.Pos{X: pos.X + 1, Y: pos.Y}
	case types.DirWest:
		return world.Pos{X: pos.X - 1, Y: pos.Y}
	default:
		return pos
	}
}

func (e *Engine) getChunksInRegion(worldPos world.Pos) []*world.Chunk {
	chunks := make([]*world.Chunk, 0)

	centerChunk := world.TileToChunkPos(worldPos)

	// Get 3x3 grid of chunks around the center
	for i := -3; i <= 3; i++ {
		for j := -3; j <= 3; j++ {
			chunkPos := world.Pos{X: centerChunk.X + i, Y: centerChunk.Y + j}

			if chunkPos.X < 0 || chunkPos.Y < 0 {
				continue
			}

			chunk := e.w.ChunkAt(chunkPos)
			chunks = append(chunks, chunk)
		}
	}
	return chunks
}

func (e *Engine) handlePlayerMessage(playerMsg playerMessage) {
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

func (e *Engine) handleChatMessage(playerMsg playerMessage) {
	entry := logrus.WithField("player", playerMsg.playerID).WithField("message", playerMsg.message.chatMessage.Message)
	e.nm.broadcastCh <- outgoingMessage{chatMessage: playerMsg.message.chatMessage}
	entry.Debug("Player sent chat message")
}

func (e *Engine) handleLoginMessage(playerMsg playerMessage) {
	entry := logrus.WithField("player", playerMsg.playerID).WithField("message", playerMsg.message.loginMessage.Username)

	camPos := world.Pos{X: e.w.Width / 2, Y: e.w.Height / 2}

	initialLoadMessage := message.InitialLoadMessage{
		Width:     e.w.Width,
		Height:    e.w.Height,
		CameraPos: world.Pos{X: camPos.X, Y: camPos.Y},
		Chunks:    e.getChunksInRegion(camPos),
		Trains:    e.w.Trains, // TODO: get trains in region
		Tracks:    e.w.Tracks, // TODO: get tracks in region
	}
	*playerMsg.responseCh <- outgoingMessage{initialLoadMessage: &initialLoadMessage}
	entry.Debug("Player sent initial load message")
}

func (e *Engine) handleGetChunksMessage(playerMsg playerMessage) {
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
	*playerMsg.responseCh <- outgoingMessage{chunksMessage: &message.ChunksMessage{Chunks: chunks}}
	entry.Debugf("Player requested chunks")
}
