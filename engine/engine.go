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

	x, y, dir := car.X, car.Y, car.Direction
	nextX, nextY := nextPos(x, y, dir)
	nextTile := e.w.TileAt(nextX, nextY)
	if nextTile.Type != types.TileTrack {
		return
	}

	if e.w.OccupiedAt(nextX, nextY) {
		return
	}

	e.moveCars(t.Cars, moveDir, t.IsReversing)

	car = t.Cars[0]
	if t.IsReversing {
		car = t.Cars[len(t.Cars)-1]
		car.Direction = types.OppositeDir(car.Direction)
	}
	x, y, dir = car.X, car.Y, car.Direction
	track := e.w.Tracks[e.w.TileAt(x, y)]

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

	newX, newY := nextPos(car.X, car.Y, moveDir)

	prevX, prevY, prevDir := car.X, car.Y, moveDir
	car.X, car.Y = newX, newY
	e.w.SetOccupied(car.X, car.Y)

	for i := start + step; i != end; i += step {
		car = cars[i]
		thisPrevX, thisPrevY, thisPrevDir := car.X, car.Y, car.Direction

		car.X, car.Y, car.Direction = prevX, prevY, prevDir

		prevX, prevY, prevDir = thisPrevX, thisPrevY, thisPrevDir
	}
	e.w.UnsetOccupied(prevX, prevY)
}

func nextPos(x, y int, dir types.Dir) (int, int) {
	switch dir {
	case types.DirNorth:
		return x, y + 1
	case types.DirSouth:
		return x, y - 1
	case types.DirEast:
		return x + 1, y
	case types.DirWest:
		return x - 1, y
	default:
		return x, y
	}
}

func (e *Engine) getChunksInRegion(worldX, worldY int) []*world.Chunk {
	chunks := make([]*world.Chunk, 0)

	centerChunkCoords := world.TileToChunkCoords(worldX, worldY)

	// Get 3x3 grid of chunks around the center
	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			chunkX := centerChunkCoords.X + i
			chunkY := centerChunkCoords.Y + j

			if chunkX < 0 || chunkY < 0 {
				continue
			}

			chunk := e.w.ChunkAt(world.ChunkCoord{X: chunkX, Y: chunkY})
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
	entry.Infof("Player sent chat message")
}

func (e *Engine) handleLoginMessage(playerMsg playerMessage) {
	entry := logrus.WithField("player", playerMsg.playerID).WithField("message", playerMsg.message.loginMessage.Username)

	camX := e.w.Width / 2
	camY := e.w.Height / 2

	initialLoadMessage := message.InitialLoadMessage{
		Width:   e.w.Width,
		Height:  e.w.Height,
		CameraX: camX,
		CameraY: camY,
		Chunks:  e.getChunksInRegion(camX, camY),
		Trains:  e.getTrainsInRegion(camX, camY),
	}
	*playerMsg.responseCh <- outgoingMessage{initialLoadMessage: &initialLoadMessage}
	entry.Infof("Player sent initial load message")
}

func (e *Engine) handleGetChunksMessage(playerMsg playerMessage) {
	entry := logrus.WithField("player", playerMsg.playerID).WithField("message", playerMsg.message.getChunksMessage)

	chunks := make([]*world.Chunk, 0, len(playerMsg.message.getChunksMessage.Coords))
	for _, coord := range playerMsg.message.getChunksMessage.Coords {
		if coord.X < 0 || coord.Y < 0 || coord.X >= e.w.Width || coord.Y >= e.w.Height {
			continue
		}
		chunks = append(chunks, e.w.ChunkAt(coord))
	}
	*playerMsg.responseCh <- outgoingMessage{chunksMessage: &message.ChunksMessage{Chunks: chunks}}
	entry.Infof("Player requested chunks")
}
func (e *Engine) getTrainsInRegion(camX, camY int) []*trains.Train {
	return e.w.Trains
}

func (e *Engine) getInitialLoadForPlayer() message.InitialLoadMessage {
	camX := e.w.Width / 2
	camY := e.w.Height / 2

	return message.InitialLoadMessage{
		Width:   e.w.Width,
		Height:  e.w.Height,
		CameraX: camX,
		CameraY: camY,
		Chunks:  e.getChunksInRegion(camX, camY),
		Trains:  e.getTrainsInRegion(camX, camY),
	}
}
