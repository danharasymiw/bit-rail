package types

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNextPos(t *testing.T) {
	origin := Pos{X: 5, Y: 5}

	assert.Equal(t, Pos{X: 5, Y: 6}, NextPos(origin, DirNorth), "North should increment Y")
	assert.Equal(t, Pos{X: 5, Y: 4}, NextPos(origin, DirSouth), "South should decrement Y")
	assert.Equal(t, Pos{X: 6, Y: 5}, NextPos(origin, DirEast), "East should increment X")
	assert.Equal(t, Pos{X: 4, Y: 5}, NextPos(origin, DirWest), "West should decrement X")
	assert.Equal(t, origin, NextPos(origin, DirNone), "Invalid/None dir should return unchanged pos")
	assert.Equal(t, origin, NextPos(origin, DirNorthSouth), "Combined dir should return unchanged pos")
}

func TestOppositeDir(t *testing.T) {
	assert.Equal(t, Dir(DirSouth), OppositeDir(DirNorth), "opposite of North is South")
	assert.Equal(t, Dir(DirNorth), OppositeDir(DirSouth), "opposite of South is North")
	assert.Equal(t, Dir(DirWest), OppositeDir(DirEast), "opposite of East is West")
	assert.Equal(t, Dir(DirEast), OppositeDir(DirWest), "opposite of West is East")
	assert.Equal(t, Dir(DirNone), OppositeDir(DirNone), "opposite of None is None")
}

func TestDirString(t *testing.T) {
	assert.Equal(t, "North", Dir(DirNorth).String())
	assert.Equal(t, "East", Dir(DirEast).String())
	assert.Equal(t, "South", Dir(DirSouth).String())
	assert.Equal(t, "West", Dir(DirWest).String())
	assert.Equal(t, "", Dir(DirNone).String())

	combined := Dir(DirNorthSouth).String()
	assert.True(t, strings.Contains(combined, "North"), "DirNorthSouth string should contain North")
	assert.True(t, strings.Contains(combined, "South"), "DirNorthSouth string should contain South")
}
