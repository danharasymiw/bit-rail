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

	// Destination is a track tile this train is routed toward — typically
	// one adjacent to a Station, in which case arriving at ANY track tile
	// adjacent to that same station counts as arrival (see arrivedAt),
	// since that's where a train pulls up to load/unload.
	// The zero value means "no destination" (direction choices at
	// junctions fall back to picking the first non-backtracking option).
	Destination types.Pos

	Cars []*TrainCar
}

// noDestination is the sentinel zero-value Destination meaning "not routed
// anywhere in particular".
var noDestination = types.Pos{}

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

	t.moveCars(w, moveDir)

	trackLeft := w.TrackAt(lastPrevPos)
	if trackLeft.HasSignal() && trackLeft.IsSignalDir(types.OppositeDir(last.Direction)) {
		trackLeft.ClearSignal()
	}

	// Arrived at destination: stop, and reverse if this is a dead end
	// (the only connecting direction is the one we just came from).
	if t.Destination != noDestination && arrivedAt(w, t.Destination, nextPos) {
		t.IsMoving = false
		if nextTrack.Direction == types.OppositeDir(moveDir) {
			reverse(t)
		}
		return
	}

	// Determine next direction; nextTrack is now our current tile.
	// If we can't continue forward, pick another valid direction that's not backtracking.
	if nextTrack.Direction&moveDir == 0 {
		newDir := types.DirNone
		if t.Destination != noDestination && nextTrack.IsJunction() {
			if dir, ok := w.NextHop(nextPos, t.Destination); ok {
				newDir = dir
			}
		}
		if newDir == types.DirNone {
			newDir = firstNonBacktrackDirection(nextTrack, moveDir)
		}
		frontCar.Direction = newDir
	}
}

// arrivedAt reports whether pos counts as having reached dest. An exact
// tile match always counts; if dest is adjacent to a station, any tile
// adjacent to that same station also counts, since a train can pull up
// alongside any side of a station to deliver goods.
func arrivedAt(w trainWorldView, dest, pos types.Pos) bool {
	if pos == dest {
		return true
	}
	destStation := w.StationAt(dest)
	return destStation != nil && w.StationAt(pos) == destStation
}

// firstNonBacktrackDirection returns the first direction nextTrack connects
// to other than the one moveDir just arrived from. Used as the fallback
// when there's no routing information (or nothing to route toward).
func firstNonBacktrackDirection(nextTrack *types.Track, moveDir types.Dir) types.Dir {
	for dir := types.Dir(types.DirNorth); dir <= types.DirWest; dir <<= 1 {
		if dir&nextTrack.Direction != 0 && dir != types.OppositeDir(moveDir) {
			return dir
		}
	}
	return types.DirNone
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

func (t *Train) moveCars(w trainWorldView, moveDir types.Dir) {
	oldPositions := make([]types.Pos, len(t.Cars))
	for i, car := range t.Cars {
		oldPositions[i] = types.Pos{X: car.X, Y: car.Y}
	}

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

	for _, pos := range oldPositions {
		w.UnsetOccupied(pos)
	}
	for _, car := range t.Cars {
		w.SetOccupied(types.Pos{X: car.X, Y: car.Y})
	}
}
