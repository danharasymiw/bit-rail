package trains

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/danharasymiw/bit-rail/types"
)

// mockWorld implements trainWorldView for testing.
type mockWorld struct {
	tiles    map[types.Pos]*types.Tile
	tracks   map[types.Pos]*types.Track
	occupied map[types.Pos]bool
	nextHop  map[types.Pos]types.Dir // keyed by pos; presence means ok=true
	stations []*types.Station
}

func newMockWorld() *mockWorld {
	return &mockWorld{
		tiles:    make(map[types.Pos]*types.Tile),
		tracks:   make(map[types.Pos]*types.Track),
		occupied: make(map[types.Pos]bool),
		nextHop:  make(map[types.Pos]types.Dir),
	}
}

func (m *mockWorld) TileAt(pos types.Pos) *types.Tile {
	return m.tiles[pos]
}

func (m *mockWorld) TrackAt(pos types.Pos) *types.Track {
	return m.tracks[pos]
}

func (m *mockWorld) OccupiedAt(pos types.Pos) bool {
	return m.occupied[pos]
}

func (m *mockWorld) SetOccupied(pos types.Pos) {
	m.occupied[pos] = true
}

func (m *mockWorld) UnsetOccupied(pos types.Pos) {
	m.occupied[pos] = false
}

func (m *mockWorld) NextHop(pos, dest types.Pos) (types.Dir, bool) {
	dir, ok := m.nextHop[pos]
	return dir, ok
}

func (m *mockWorld) StationAt(pos types.Pos) *types.Station {
	return types.StationAt(pos, m.stations)
}

// --- firstCar / lastCar ---

func TestFirstCar_Forward(t *testing.T) {
	train := &Train{
		IsReversing: false,
		Cars: []*TrainCar{
			{X: 1, Y: 1, Direction: types.DirEast},
			{X: 2, Y: 1, Direction: types.DirEast},
			{X: 3, Y: 1, Direction: types.DirEast},
		},
	}
	fc := firstCar(train)
	assert.Equal(t, train.Cars[0], fc)
}

func TestFirstCar_Reversing(t *testing.T) {
	train := &Train{
		IsReversing: true,
		Cars: []*TrainCar{
			{X: 1, Y: 1, Direction: types.DirEast},
			{X: 2, Y: 1, Direction: types.DirEast},
			{X: 3, Y: 1, Direction: types.DirEast},
		},
	}
	fc := firstCar(train)
	assert.Equal(t, train.Cars[len(train.Cars)-1], fc)
}

func TestLastCar_Forward(t *testing.T) {
	train := &Train{
		IsReversing: false,
		Cars: []*TrainCar{
			{X: 1, Y: 1, Direction: types.DirEast},
			{X: 2, Y: 1, Direction: types.DirEast},
			{X: 3, Y: 1, Direction: types.DirEast},
		},
	}
	lc := lastCar(train)
	assert.Equal(t, train.Cars[len(train.Cars)-1], lc)
}

func TestLastCar_Reversing(t *testing.T) {
	train := &Train{
		IsReversing: true,
		Cars: []*TrainCar{
			{X: 1, Y: 1, Direction: types.DirEast},
			{X: 2, Y: 1, Direction: types.DirEast},
			{X: 3, Y: 1, Direction: types.DirEast},
		},
	}
	lc := lastCar(train)
	assert.Equal(t, train.Cars[0], lc)
}

// --- reverse ---

func TestReverse(t *testing.T) {
	train := &Train{
		IsReversing: false,
		Cars: []*TrainCar{
			{Direction: types.DirEast},
			{Direction: types.DirEast},
			{Direction: types.DirWest},
		},
	}
	originalDirs := []types.Dir{types.DirEast, types.DirEast, types.DirWest}

	reverse(train)

	assert.True(t, train.IsReversing, "IsReversing should be toggled to true")
	for i, car := range train.Cars {
		expected := types.OppositeDir(originalDirs[i])
		assert.Equal(t, expected, car.Direction, "car %d direction should be opposite", i)
	}

	// Toggle back
	reverse(train)
	assert.False(t, train.IsReversing, "IsReversing should toggle back to false")
}

// --- moveCars ---

func TestMoveCars_SingleCar(t *testing.T) {
	w := newMockWorld()

	startPos := types.Pos{X: 5, Y: 5}
	nextPos := types.Pos{X: 6, Y: 5}

	w.tracks[nextPos] = &types.Track{
		Direction: types.DirEastWest,
		Block:     types.NewBlock(),
	}
	// Mark start as occupied initially.
	w.occupied[startPos] = true

	train := &Train{
		IsReversing: false,
		Cars: []*TrainCar{
			{X: 5, Y: 5, Direction: types.DirEast},
		},
	}

	train.moveCars(w, types.DirEast)

	car := train.Cars[0]
	assert.Equal(t, 6, car.X, "car X should be 6 after moving east")
	assert.Equal(t, 5, car.Y, "car Y should remain 5")

	assert.False(t, w.occupied[startPos], "old position (5,5) should be unoccupied")
	assert.True(t, w.occupied[nextPos], "new position (6,5) should be occupied")
}

func TestMoveCars_TwoCars(t *testing.T) {
	w := newMockWorld()

	frontPos := types.Pos{X: 5, Y: 5}
	rearPos := types.Pos{X: 4, Y: 5}
	nextFrontPos := types.Pos{X: 6, Y: 5}

	w.tracks[nextFrontPos] = &types.Track{
		Direction: types.DirEastWest,
		Block:     types.NewBlock(),
	}
	w.occupied[frontPos] = true
	w.occupied[rearPos] = true

	train := &Train{
		IsReversing: false,
		Cars: []*TrainCar{
			{X: 5, Y: 5, Direction: types.DirEast}, // front
			{X: 4, Y: 5, Direction: types.DirEast}, // rear
		},
	}

	train.moveCars(w, types.DirEast)

	front := train.Cars[0]
	rear := train.Cars[1]

	// Front car moves to (6,5).
	assert.Equal(t, 6, front.X, "front car X should be 6")
	assert.Equal(t, 5, front.Y, "front car Y should be 5")

	// Rear car moves to where front was (5,5) and inherits the move direction.
	assert.Equal(t, 5, rear.X, "rear car X should be 5")
	assert.Equal(t, 5, rear.Y, "rear car Y should be 5")
	assert.Equal(t, types.Dir(types.DirEast), rear.Direction, "rear car Direction should be DirEast")

	// Occupancy: old positions cleared, new positions set.
	assert.False(t, w.occupied[rearPos], "old rear position (4,5) should be unoccupied")
	assert.True(t, w.occupied[frontPos], "middle position (5,5) should now be occupied by rear")
	assert.True(t, w.occupied[nextFrontPos], "new front position (6,5) should be occupied")
}

// --- Tick ---

func TestTick_NotMoving(t *testing.T) {
	w := newMockWorld()

	train := &Train{
		ID:          uuid.New(),
		IsMoving:    false,
		IsReversing: false,
		Cars: []*TrainCar{
			{X: 5, Y: 5, Direction: types.DirEast},
		},
	}

	train.Tick(w)

	assert.Equal(t, 5, train.Cars[0].X, "car X should not change when not moving")
	assert.Equal(t, 5, train.Cars[0].Y, "car Y should not change when not moving")
}

func TestTick_MovesNormally(t *testing.T) {
	w := newMockWorld()

	// Current track (needed for signal-clear check on lastPrevPos after move).
	currPos := types.Pos{X: 5, Y: 5}
	nextPos := types.Pos{X: 6, Y: 5}

	w.tracks[currPos] = &types.Track{
		Direction: types.DirEastWest,
		Block:     types.NewBlock(),
	}
	w.tracks[nextPos] = &types.Track{
		Direction: types.DirEastWest,
		Block:     types.NewBlock(),
	}
	w.occupied[currPos] = true

	train := &Train{
		ID:          uuid.New(),
		IsMoving:    true,
		IsReversing: false,
		Cars: []*TrainCar{
			{X: 5, Y: 5, Direction: types.DirEast},
		},
	}

	train.Tick(w)

	car := train.Cars[0]
	assert.Equal(t, 6, car.X, "car X should advance to 6 after tick")
	assert.Equal(t, 5, car.Y, "car Y should remain 5")
	assert.True(t, train.IsMoving, "train should still be moving")
}

func TestTick_StopsOnNoTrack(t *testing.T) {
	w := newMockWorld()

	// No track at nextPos (6,5) — nextTrack will be nil.
	currPos := types.Pos{X: 5, Y: 5}
	w.tracks[currPos] = &types.Track{
		Direction: types.DirEastWest,
		Block:     types.NewBlock(),
	}
	w.occupied[currPos] = true

	train := &Train{
		ID:          uuid.New(),
		IsMoving:    true,
		IsReversing: false,
		Cars: []*TrainCar{
			{X: 5, Y: 5, Direction: types.DirEast},
		},
	}

	train.Tick(w)

	// Tick returns early when nextTrack is nil; car should not have moved.
	car := train.Cars[0]
	assert.Equal(t, 5, car.X, "car should not move when there is no next track")
	assert.Equal(t, 5, car.Y, "car Y should remain 5")
}

func TestTick_StopsOnOccupied(t *testing.T) {
	// The engine uses block-level occupancy via signals to implement collision
	// avoidance. When the next track's block is occupied by another train and
	// the track has a signal facing the approach direction, Tick returns without
	// moving the car.
	w := newMockWorld()

	currPos := types.Pos{X: 5, Y: 5}
	nextPos := types.Pos{X: 6, Y: 5}

	block := types.NewBlock()
	otherTrainID := uuid.New()
	block.OccupiedBy = otherTrainID // block is held by another train

	w.tracks[currPos] = &types.Track{
		Direction: types.DirEastWest,
		Block:     types.NewBlock(),
	}
	// Signal faces East (the approach direction), so IsSignalDir(DirEast) is true.
	w.tracks[nextPos] = &types.Track{
		Direction: types.DirEastWest,
		SignalDir: types.DirEast,
		Block:     block,
	}
	w.occupied[currPos] = true

	trainID := uuid.New()
	train := &Train{
		ID:          trainID,
		IsMoving:    true,
		IsReversing: false,
		Cars: []*TrainCar{
			{X: 5, Y: 5, Direction: types.DirEast},
		},
	}

	train.Tick(w)

	// Tick returns early (without IsMoving=false) when the signal is not clear.
	car := train.Cars[0]
	require.Equal(t, 5, car.X, "car should not move when next block is occupied")
	require.Equal(t, 5, car.Y, "car Y should remain 5")
}

func TestTick_SignalBlocks(t *testing.T) {
	// Same as StopsOnOccupied but verifies that a signal whose block is held by
	// another train prevents movement, while IsMoving itself is not cleared
	// (Tick uses an early return, not IsMoving=false for this case).
	w := newMockWorld()

	currPos := types.Pos{X: 5, Y: 5}
	nextPos := types.Pos{X: 6, Y: 5}

	otherID := uuid.New()
	blockedBlock := &types.Block{
		ID:         uuid.New(),
		OccupiedBy: otherID,
	}

	w.tracks[currPos] = &types.Track{
		Direction: types.DirEastWest,
		Block:     types.NewBlock(),
	}
	w.tracks[nextPos] = &types.Track{
		Direction: types.DirEastWest,
		SignalDir: types.DirEast, // signal faces East — train approaches from West going East
		Block:     blockedBlock,
	}
	w.occupied[currPos] = true

	trainID := uuid.New()
	train := &Train{
		ID:          trainID,
		IsMoving:    true,
		IsReversing: false,
		Cars: []*TrainCar{
			{X: 5, Y: 5, Direction: types.DirEast},
		},
	}

	train.Tick(w)

	// Train should not have moved.
	assert.Equal(t, 5, train.Cars[0].X, "car X should not change when signal blocks")
	assert.Equal(t, 5, train.Cars[0].Y, "car Y should not change when signal blocks")
	// IsMoving is not set to false by the signal check — Tick just returns early.
	assert.True(t, train.IsMoving, "IsMoving should remain true when blocked by signal")
	// Block ownership should be unchanged.
	assert.Equal(t, otherID, blockedBlock.OccupiedBy, "block should still be owned by other train")
}

// --- Destination-aware junction routing + station arrival ---

func TestTick_JunctionUsesRoutingWhenDestinationSet(t *testing.T) {
	w := newMockWorld()

	currPos := types.Pos{X: 5, Y: 5}
	junctionPos := types.Pos{X: 6, Y: 5}
	dest := types.Pos{X: 100, Y: 100}

	w.tracks[currPos] = &types.Track{Direction: types.DirEastWest, Block: types.NewBlock()}
	// T-junction (North|South|West, no through-East): a train arriving from
	// the West moving East cannot continue straight, so it must branch
	// North or South. Without routing, firstNonBacktrackDirection would
	// pick North (first checked direction that isn't the backtrack West).
	w.tracks[junctionPos] = &types.Track{Direction: types.DirNorth | types.DirSouth | types.DirWest, Block: types.NewBlock()}
	w.occupied[currPos] = true
	// Routing says: from the junction, go South to reach dest.
	w.nextHop[junctionPos] = types.DirSouth

	train := &Train{
		ID:          uuid.New(),
		IsMoving:    true,
		Destination: dest,
		Cars: []*TrainCar{
			{X: 5, Y: 5, Direction: types.DirEast},
		},
	}

	train.Tick(w)

	assert.Equal(t, types.Dir(types.DirSouth), train.Cars[0].Direction,
		"should take the routed direction instead of the arbitrary first choice")
}

func TestTick_JunctionFallsBackWhenNoRoute(t *testing.T) {
	w := newMockWorld()

	currPos := types.Pos{X: 5, Y: 5}
	junctionPos := types.Pos{X: 6, Y: 5}
	dest := types.Pos{X: 100, Y: 100}

	w.tracks[currPos] = &types.Track{Direction: types.DirEastWest, Block: types.NewBlock()}
	w.tracks[junctionPos] = &types.Track{Direction: types.DirNorth | types.DirSouth | types.DirWest, Block: types.NewBlock()}
	w.occupied[currPos] = true
	// No entry in w.nextHop for junctionPos -> NextHop returns ok=false.

	train := &Train{
		ID:          uuid.New(),
		IsMoving:    true,
		Destination: dest,
		Cars: []*TrainCar{
			{X: 5, Y: 5, Direction: types.DirEast},
		},
	}

	train.Tick(w)

	assert.Equal(t, types.Dir(types.DirNorth), train.Cars[0].Direction,
		"should fall back to the first non-backtracking direction when routing has no answer")
}

func TestTick_ArrivesAtDestinationAndStops(t *testing.T) {
	w := newMockWorld()

	currPos := types.Pos{X: 5, Y: 5}
	stationPos := types.Pos{X: 6, Y: 5}

	w.tracks[currPos] = &types.Track{Direction: types.DirEastWest, Block: types.NewBlock()}
	// Through-station (connects both East and West) — not a dead end.
	w.tracks[stationPos] = &types.Track{Direction: types.DirEastWest, Block: types.NewBlock()}
	w.occupied[currPos] = true

	train := &Train{
		ID:          uuid.New(),
		IsMoving:    true,
		Destination: stationPos,
		Cars: []*TrainCar{
			{X: 5, Y: 5, Direction: types.DirEast},
		},
	}

	train.Tick(w)

	assert.Equal(t, 6, train.Cars[0].X, "train should have moved onto the station tile")
	assert.False(t, train.IsMoving, "train should stop upon arriving at its destination")
	assert.False(t, train.IsReversing, "through-station arrival should not reverse the train")
}

func TestTick_ArrivesAtDeadEndStationAndReverses(t *testing.T) {
	w := newMockWorld()

	currPos := types.Pos{X: 5, Y: 5}
	stationPos := types.Pos{X: 6, Y: 5}

	w.tracks[currPos] = &types.Track{Direction: types.DirEastWest, Block: types.NewBlock()}
	// Dead-end station: only connects back West (the direction we arrived from).
	w.tracks[stationPos] = &types.Track{Direction: types.DirWest, Block: types.NewBlock()}
	w.occupied[currPos] = true

	train := &Train{
		ID:          uuid.New(),
		IsMoving:    true,
		Destination: stationPos,
		Cars: []*TrainCar{
			{X: 5, Y: 5, Direction: types.DirEast},
		},
	}

	train.Tick(w)

	assert.False(t, train.IsMoving, "train should stop upon arriving at its destination")
	assert.True(t, train.IsReversing, "dead-end arrival should reverse the train")
	assert.Equal(t, types.Dir(types.DirWest), train.Cars[0].Direction, "car direction should flip on reverse")
}
