package engine

import (
	"time"

	"github.com/danharasymiw/trains/log"
	"github.com/danharasymiw/trains/renderer"
	"github.com/danharasymiw/trains/trains"
	"github.com/danharasymiw/trains/types"
	"github.com/danharasymiw/trains/world"
	"github.com/gdamore/tcell"
)

type Engine struct {
	w       *world.World
	tickDur time.Duration
	Running bool
	r       renderer.Renderer
}

func New(w *world.World, r renderer.Renderer, tickDur time.Duration) *Engine {
	return &Engine{
		w:       w,
		r:       r,
		tickDur: tickDur,
	}
}

func (e *Engine) Run() {
	events := make(chan tcell.Event, 16)
	go func() {
		for {
			ev := e.r.Screen().PollEvent() // blocks here
			events <- ev
		}
	}()
	ticker := time.NewTicker(e.tickDur)
	defer ticker.Stop()

	e.Running = true
	for e.Running {
		<-ticker.C
		e.tick()
		e.r.Draw()

		// Process all queued input (non-blocking)
		for {
			select {
			case ev := <-events:
				if _, ok := ev.(*tcell.EventKey); ok {
					e.Running = false // quit on any key
				}
			default:
				break
			}
		}
	}
}

func (e *Engine) tick() {
	for _, t := range e.w.Trains {
		log.Log("TICK")
		e.moveTrain(t)
	}
}

func (e *Engine) moveTrain(t *trains.Train) {
	car := t.Cars[0]
	if t.IsReversing {
		car = t.Cars[len(t.Cars)-1]
	}

	x, y, dir := car.X, car.Y, car.Direction
	nextX, nextY := nextCarPos(x, y, dir)
	nextTile := e.w.TileAt(nextX, nextY)
	if nextTile.Type != types.TileTrack {
		return
	}

	e.moveCars(t.Cars, t.IsReversing)

	car = t.Cars[0]
	if t.IsReversing {
		car = t.Cars[len(t.Cars)-1]
		car.Direction = types.OppositeDir(car.Direction)
	}
	x, y, dir = car.X, car.Y, car.Direction
	tile := e.w.TileAt(x, y)

	incFrom := types.OppositeDir(dir)
	if tile.Orientation&incFrom == 0 {
		return
	}

	outgoing := tile.Orientation & ^incFrom

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

func (e *Engine) moveCars(cars []*trains.TrainCar, reverse bool) {
	start, end, step := 0, len(cars), 1
	if reverse {
		start, end, step = len(cars)-1, -1, -1
	}

	car := cars[start]

	newX, newY := nextCarPos(car.X, car.Y, car.Direction)
	if e.w.OccupiedAt(newX, newY) {
		return // blocked
	}

	prevX, prevY, prevDir := car.X, car.Y, car.Direction
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

func nextCarPos(x, y int, dir types.Dir) (int, int) {
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
