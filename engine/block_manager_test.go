package engine

import (
	"testing"

	"github.com/danharasymiw/bit-rail/types"
	"github.com/danharasymiw/bit-rail/world"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestIsSignalBoundary tests the unexported isSignalBoundary function.
func TestIsSignalBoundary(t *testing.T) {
	t.Run("no signals on either track returns false", func(t *testing.T) {
		curr := &types.Track{Direction: types.DirEastWest}
		neigh := &types.Track{Direction: types.DirEastWest}
		assert.False(t, isSignalBoundary(curr, neigh, types.DirEast))
	})

	t.Run("neighbour has signal facing back toward curr returns true", func(t *testing.T) {
		// Travelling East (dir=DirEast). Neighbour has SignalDir=DirEast,
		// meaning the signal faces East — i.e. it checks trains coming from the West,
		// which is the direction we are travelling from.
		// isSignalBoundary: neighbour.HasSignal() && neighbour.IsSignalDir(dir=DirEast) → true
		curr := &types.Track{Direction: types.DirEastWest}
		neigh := &types.Track{Direction: types.DirEastWest, SignalDir: types.DirEast}
		assert.True(t, isSignalBoundary(curr, neigh, types.DirEast))
	})

	t.Run("curr has signal facing toward neighbour returns true", func(t *testing.T) {
		// Travelling East (dir=DirEast). curr has SignalDir=DirWest (OppositeDir(DirEast)).
		// isSignalBoundary: curr.HasSignal() && curr.IsSignalDir(OppositeDir(DirEast)=DirWest) → true
		curr := &types.Track{Direction: types.DirEastWest, SignalDir: types.DirWest}
		neigh := &types.Track{Direction: types.DirEastWest}
		assert.True(t, isSignalBoundary(curr, neigh, types.DirEast))
	})

	t.Run("signal on neighbour facing a different direction returns false", func(t *testing.T) {
		// Travelling East (dir=DirEast). Neighbour has SignalDir=DirNorth (unrelated direction).
		// neighbour.IsSignalDir(DirEast) → DirNorth & DirEast = 0 → false
		// curr has no signal → false
		curr := &types.Track{Direction: types.DirEastWest}
		neigh := &types.Track{Direction: types.DirEastWest, SignalDir: types.DirNorth}
		assert.False(t, isSignalBoundary(curr, neigh, types.DirEast))
	})

	t.Run("junction alone is not a boundary", func(t *testing.T) {
		// Junctions do NOT split blocks — see TestRebuildAll_JunctionDoesNotSplitBlock.
		curr := &types.Track{Direction: types.DirEastWest}
		neigh := &types.Track{Direction: types.DirFourWay}
		assert.False(t, isSignalBoundary(curr, neigh, types.DirEast))
	})
}

// newTestWorld creates a small world suitable for testing (avoids Perlin generation).
func newTestWorld() *world.World {
	return world.New(16, 16)
}

// TestRebuildAll_LinearTrack verifies that three connected E-W tracks all get
// assigned the same non-nil Block after RebuildAll.
func TestRebuildAll_LinearTrack(t *testing.T) {
	w := newTestWorld()

	// Place three EastWest tracks in a row: (0,0) → (1,0) → (2,0)
	pos0 := types.Pos{X: 0, Y: 0}
	pos1 := types.Pos{X: 1, Y: 0}
	pos2 := types.Pos{X: 2, Y: 0}

	w.AddTrack(pos0, &types.Track{Direction: types.DirEastWest})
	w.AddTrack(pos1, &types.Track{Direction: types.DirEastWest})
	w.AddTrack(pos2, &types.Track{Direction: types.DirEastWest})

	bm := newBlockManager(w)
	bm.RebuildAll()

	tr0 := w.Tracks[pos0]
	tr1 := w.Tracks[pos1]
	tr2 := w.Tracks[pos2]

	require.NotNil(t, tr0.Block, "track at (0,0) should have a block")
	require.NotNil(t, tr1.Block, "track at (1,0) should have a block")
	require.NotNil(t, tr2.Block, "track at (2,0) should have a block")

	assert.Equal(t, tr0.Block.ID, tr1.Block.ID, "tracks (0,0) and (1,0) should share the same block")
	assert.Equal(t, tr1.Block.ID, tr2.Block.ID, "tracks (1,0) and (2,0) should share the same block")
	assert.Same(t, tr0.Block, tr1.Block, "tracks (0,0) and (1,0) should share the same block pointer")
	assert.Same(t, tr1.Block, tr2.Block, "tracks (1,0) and (2,0) should share the same block pointer")
}

// TestRebuildAll_SignalSplits verifies that a signal at (1,0) with SignalDir=DirEast
// creates a boundary so that track (0,0) is in a different block from (1,0) and (2,0).
//
// Signal logic (isSignalBoundary):
//   When BFS travels East from (0,0) to (1,0):
//     neighbour(1,0).HasSignal() && neighbour(1,0).IsSignalDir(DirEast) → true → boundary
//   So (0,0) and (1,0) are in separate blocks.
//
//   When BFS travels East from (1,0) to (2,0):
//     neighbour(2,0).HasSignal() = false
//     curr(1,0).IsSignalDir(OppositeDir(DirEast)=DirWest) = SignalDir(DirEast) & DirWest = 0 → false
//   So (1,0) and (2,0) are in the same block.
func TestRebuildAll_SignalSplits(t *testing.T) {
	w := newTestWorld()

	pos0 := types.Pos{X: 0, Y: 0}
	pos1 := types.Pos{X: 1, Y: 0}
	pos2 := types.Pos{X: 2, Y: 0}

	w.AddTrack(pos0, &types.Track{Direction: types.DirEastWest})
	// Signal at (1,0) faces East: trains coming from the West must stop here.
	w.AddTrack(pos1, &types.Track{Direction: types.DirEastWest, SignalDir: types.DirEast})
	w.AddTrack(pos2, &types.Track{Direction: types.DirEastWest})

	bm := newBlockManager(w)
	bm.RebuildAll()

	tr0 := w.Tracks[pos0]
	tr1 := w.Tracks[pos1]
	tr2 := w.Tracks[pos2]

	require.NotNil(t, tr0.Block, "track at (0,0) should have a block")
	require.NotNil(t, tr1.Block, "track at (1,0) should have a block")
	require.NotNil(t, tr2.Block, "track at (2,0) should have a block")

	// (0,0) must be isolated from (1,0) due to the signal boundary.
	assert.NotEqual(t, tr0.Block.ID, tr1.Block.ID,
		"signal at (1,0) facing East should split (0,0) into its own block")

	// (1,0) and (2,0) are on the unguarded side of the signal — same block.
	assert.Equal(t, tr1.Block.ID, tr2.Block.ID,
		"tracks (1,0) and (2,0) should share the same block (no boundary between them)")
}

// TestRebuildAll_JunctionDoesNotSplitBlock verifies that a 4-way junction
// tile is NOT a block boundary: all four arms plus the junction itself
// share one block, exactly as if the junction weren't there. This matters
// for diamond crossings guarded by signals on each approach — the whole
// crossing must be one shared block so perpendicular trains see each
// other's occupancy. (Junctions ARE decision points for the separate
// routing graph in the routing package — see types.Track.IsNode — but that
// is independent of block/collision topology.)
func TestRebuildAll_JunctionDoesNotSplitBlock(t *testing.T) {
	w := newTestWorld()

	junctionPos := types.Pos{X: 2, Y: 2}
	w.AddTrack(junctionPos, &types.Track{Direction: types.DirFourWay})

	north := types.Pos{X: 2, Y: 3}
	south := types.Pos{X: 2, Y: 1}
	east := types.Pos{X: 3, Y: 2}
	west := types.Pos{X: 1, Y: 2}
	w.AddTrack(north, &types.Track{Direction: types.DirSouth})
	w.AddTrack(south, &types.Track{Direction: types.DirNorth})
	w.AddTrack(east, &types.Track{Direction: types.DirWest})
	w.AddTrack(west, &types.Track{Direction: types.DirEast})

	bm := newBlockManager(w)
	bm.RebuildAll()

	junctionBlock := w.Tracks[junctionPos].Block
	require.NotNil(t, junctionBlock, "junction tile should have a block")

	for name, pos := range map[string]types.Pos{"north": north, "south": south, "east": east, "west": west} {
		tr := w.Tracks[pos]
		require.NotNil(t, tr.Block, "%s arm should have a block", name)
		assert.Same(t, junctionBlock, tr.Block, "%s arm should share the junction's block", name)
	}
}

// TestMarkDirtyProcessDirty verifies that a newly added track gets a block
// assigned after MarkDirty + ProcessDirty, without requiring RebuildAll.
func TestMarkDirtyProcessDirty(t *testing.T) {
	w := newTestWorld()

	// Seed the world with one track and build its block.
	pos0 := types.Pos{X: 5, Y: 5}
	w.AddTrack(pos0, &types.Track{Direction: types.DirEastWest})

	bm := newBlockManager(w)
	bm.RebuildAll()

	require.NotNil(t, w.Tracks[pos0].Block, "initial track should have a block after RebuildAll")

	// Add a disconnected track at a new position (not adjacent to pos0).
	pos1 := types.Pos{X: 10, Y: 10}
	w.AddTrack(pos1, &types.Track{Direction: types.DirNorthSouth})

	// The new track must have no block yet (RebuildAll was not called again).
	require.Nil(t, w.Tracks[pos1].Block, "new track should not have a block before ProcessDirty")

	bm.MarkDirty(pos1)
	bm.ProcessDirty()

	assert.NotNil(t, w.Tracks[pos1].Block, "new track should have a block after MarkDirty + ProcessDirty")
}
