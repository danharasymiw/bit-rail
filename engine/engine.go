package engine

import (
	"time"

	"github.com/danharasymiw/bit-rail/client"
	"github.com/danharasymiw/bit-rail/message"
	"github.com/danharasymiw/bit-rail/trains"
	"github.com/danharasymiw/bit-rail/types"
	"github.com/danharasymiw/bit-rail/world"
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
	eng.nm = newNetworkManager(eng.getInitialLoadForPlayer)
	return eng
}

func (e *Engine) Run() {
	go e.nm.startServer()

	ticker := time.NewTicker(e.tickDur)
	defer ticker.Stop()

	// TODO one day add headless mode
	localClient, quitCh := client.New()
	go func() {
		time.Sleep(1 * time.Second)
		err := localClient.Run()
		if err != nil {
			panic(err)
		}
	}()

	e.running = true
	for e.running {
		select {
		case <-ticker.C:
			e.tick()
		case <-quitCh:
			e.running = false
		}
	}
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
		return x, y - 1 // Think I'll have to reverse the Y to normal for actual client
	case types.DirSouth:
		return x, y + 1
	case types.DirEast:
		return x + 1, y
	case types.DirWest:
		return x - 1, y
	default:
		return x, y
	}
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

func (e *Engine) getChunksInRegion(chunkX, chunkY int) []message.Chunk {
	chunks := make([]message.Chunk, 0)
	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			// TODO: chunk size const somewhere
			worldX := (chunkX + i) * 32
			worldY := (chunkY + j) * 32
			tiles := make([]*types.Tile, 0, 32*32)
			for y := worldY; y < worldY+32; y++ {
				for x := worldX; x < worldX+32; x++ {
					tiles = append(tiles, e.w.Tiles[y][x])
				}
			}
			chunks = append(chunks, message.Chunk{
				X:     chunkX + i,
				Y:     chunkY + j,
				Size:  32,
				Tiles: tiles,
			})
		}
	}
	return chunks
}

func (e *Engine) getTrainsInRegion(camX, camY int) []*trains.Train {
	return e.w.Trains
}
