package trains

import (
	"github.com/danharasymiw/bit-rail/types"
	"github.com/google/uuid"
)

type Train struct {
	ID           uuid.UUID
	IsReversing  bool
	IsMoving     bool
	Speed        int
	Acceleration int

	Cars []*TrainCar
}

type CarType uint

const (
	CarTypeLocomotive CarType = iota
	CarTypeCargo
	CarTypePassenger
)

type TrainCar struct {
	X, Y      int `binary:"uint16"`
	Direction types.Dir
	Type      CarType
}

func (t *Train) Tick(w trainWorldView) {
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

	pos := types.Pos{X: car.X, Y: car.Y}
	dir := car.Direction
	nextPos := nextPos(pos, dir)
	nextTile := w.TileAt(nextPos)
	if nextTile.Type != types.TileTrack {
		return
	}

	if w.OccupiedAt(nextPos) {
		return
	}

	t.moveCars(moveDir, w)

	car = t.Cars[0]
	if t.IsReversing {
		car = t.Cars[len(t.Cars)-1]
		car.Direction = types.OppositeDir(car.Direction)
	}
	pos = types.Pos{X: car.X, Y: car.Y}
	dir = car.Direction
	track := w.TrackAt(pos)
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

func (t *Train) moveCars(moveDir types.Dir, w trainWorldView) {
	start, end, step := 0, len(t.Cars), 1
	if t.IsReversing {
		start, end, step = len(t.Cars)-1, -1, -1
	}

	car := t.Cars[start]

	newPos := nextPos(types.Pos{X: car.X, Y: car.Y}, moveDir)

	prevPos := types.Pos{X: car.X, Y: car.Y}
	prevDir := moveDir
	car.X, car.Y = newPos.X, newPos.Y
	w.SetOccupied(types.Pos{X: car.X, Y: car.Y})

	for i := start + step; i != end; i += step {
		car = t.Cars[i]
		thisPrevPos := types.Pos{X: car.X, Y: car.Y}
		thisPrevDir := car.Direction

		car.X, car.Y, car.Direction = prevPos.X, prevPos.Y, prevDir

		prevPos, prevDir = thisPrevPos, thisPrevDir
	}
	w.UnsetOccupied(prevPos)
}

// TODO: This is duplicated in block_manager.go
func nextPos(pos types.Pos, dir types.Dir) types.Pos {
	switch dir {
	case types.DirNorth:
		return types.Pos{X: pos.X, Y: pos.Y + 1}
	case types.DirSouth:
		return types.Pos{X: pos.X, Y: pos.Y - 1}
	case types.DirEast:
		return types.Pos{X: pos.X + 1, Y: pos.Y}
	case types.DirWest:
		return types.Pos{X: pos.X - 1, Y: pos.Y}
	default:
		return pos
	}
}
