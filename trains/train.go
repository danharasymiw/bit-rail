package trains

import (
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/danharasymiw/bit-rail/types"
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
	if !t.IsMoving {
		return
	}

	frontCar := firstCar(t)
	moveDir := frontCar.Direction

	currPos := types.Pos{X: frontCar.X, Y: frontCar.Y}
	// TODO: maybe set direction before moving
	// Awkward that our current track may not work with our direction.

	nextPos := types.NextPos(currPos, moveDir)
	nextTrack := w.TrackAt(nextPos)
	if nextTrack == nil {
		logrus.Debug("Next track is nil!")
		return // nowhere to go
	}

	// Make sure next track connects back (TODO otherwise blow up or whatever trains do)
	if types.OppositeDir(moveDir)&nextTrack.Direction == 0 {
		logrus.Debug("BOOM - train can't move tracks don't align")
		logrus.Debug(nextPos)
		t.IsMoving = false
		return
	}

	if nextTrack.HasSignal() && nextTrack.IsSignalDir(frontCar.Direction) {
		if !nextTrack.IsSignalClear(t.ID) {
			return
		}

		nextTrack.SetSignal(t.ID)
	}

	// Move (remember tail car so we can clear signals after tail passes)
	last := lastCar(t)
	lastPrevPos := types.Pos{X: last.X, Y: last.Y}

	t.moveCars(moveDir)

	trackLeft := w.TrackAt(lastPrevPos)
	if trackLeft.HasSignal() && trackLeft.IsSignalDir(types.OppositeDir(last.Direction)) {
		trackLeft.ClearSignal()
	}

	// Determine next direction; nextTrack is now our current tile.
	// If we can't continue forward, pick another valid direction that's not backtracking.
	if nextTrack.Direction&moveDir == 0 {
		for dir := types.Dir(types.DirNorth); dir <= types.DirWest; dir <<= 1 {
			if dir&nextTrack.Direction != 0 && dir != types.OppositeDir(moveDir) {
				frontCar.Direction = dir
				break
			}
		}
	}
}

func firstCar(t *Train) *TrainCar {
	if !t.IsReversing {
		return t.Cars[0]
	}
	return t.Cars[len(t.Cars)-1]
}

func lastCar(t *Train) *TrainCar {
	if t.IsReversing {
		return t.Cars[0]
	}
	return t.Cars[len(t.Cars)-1]
}

func reverse(t *Train) {
	t.IsReversing = !t.IsReversing
	for _, c := range t.Cars {
		c.Direction = types.OppositeDir(c.Direction)
	}
}

func (t *Train) moveCars(moveDir types.Dir) {
	start, end, step := 0, len(t.Cars), 1
	if t.IsReversing {
		start, end, step = len(t.Cars)-1, -1, -1
	}

	car := t.Cars[start]
	prevPos := types.Pos{X: car.X, Y: car.Y}
	prevDir := moveDir
	newPos := types.NextPos(prevPos, moveDir)

	car.X, car.Y = newPos.X, newPos.Y

	for i := start + step; i != end; i += step {
		car = t.Cars[i]
		thisPrevPos := types.Pos{X: car.X, Y: car.Y}
		thisPrevDir := car.Direction

		car.X, car.Y, car.Direction = prevPos.X, prevPos.Y, prevDir

		prevPos, prevDir = thisPrevPos, thisPrevDir
	}
}
