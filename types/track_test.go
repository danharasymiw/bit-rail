package types

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHasSignal(t *testing.T) {
	t.Run("no signal when SignalDir is zero", func(t *testing.T) {
		track := Track{Direction: DirNorthSouth, SignalDir: DirNone, Block: NewBlock()}
		assert.False(t, track.HasSignal())
	})

	t.Run("has signal when SignalDir is nonzero", func(t *testing.T) {
		track := Track{Direction: DirNorthSouth, SignalDir: DirNorth, Block: NewBlock()}
		assert.True(t, track.HasSignal())
	})
}

func TestIsSignalDir(t *testing.T) {
	track := Track{Direction: DirNorthSouth, SignalDir: DirNorth, Block: NewBlock()}

	assert.True(t, track.IsSignalDir(DirNorth), "DirNorth matches signal direction")
	assert.False(t, track.IsSignalDir(DirSouth), "DirSouth does not match signal direction")
	assert.False(t, track.IsSignalDir(DirEast), "DirEast does not match signal direction")
	assert.False(t, track.IsSignalDir(DirWest), "DirWest does not match signal direction")
}

func TestIsSignalClear_Unoccupied(t *testing.T) {
	block := NewBlock()
	require.Equal(t, uuid.Nil, block.OccupiedBy, "new block should be unoccupied")

	track := Track{Direction: DirNorthSouth, SignalDir: DirNorth, Block: block}
	trainID := uuid.New()

	assert.True(t, track.IsSignalClear(trainID), "signal should be clear when block is unoccupied")
}

func TestIsSignalClear_OccupiedBySelf(t *testing.T) {
	trainID := uuid.New()
	block := NewBlock()
	block.OccupiedBy = trainID

	track := Track{Direction: DirNorthSouth, SignalDir: DirNorth, Block: block}

	assert.True(t, track.IsSignalClear(trainID), "signal should be clear when occupied by the same train")
}

func TestIsSignalClear_OccupiedByOther(t *testing.T) {
	otherTrain := uuid.New()
	block := NewBlock()
	block.OccupiedBy = otherTrain

	track := Track{Direction: DirNorthSouth, SignalDir: DirNorth, Block: block}

	myTrain := uuid.New()
	assert.False(t, track.IsSignalClear(myTrain), "signal should not be clear when occupied by a different train")
}

func TestSetSignal(t *testing.T) {
	block := NewBlock()
	track := Track{Direction: DirNorthSouth, SignalDir: DirNorth, Block: block}

	trainID := uuid.New()
	track.SetSignal(trainID)

	assert.Equal(t, trainID, block.OccupiedBy)
}

func TestClearSignal(t *testing.T) {
	trainID := uuid.New()
	block := NewBlock()
	block.OccupiedBy = trainID

	track := Track{Direction: DirNorthSouth, SignalDir: DirNorth, Block: block}
	track.ClearSignal()

	assert.Equal(t, uuid.Nil, block.OccupiedBy)
}
