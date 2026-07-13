package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func testStation() Station {
	return Station{Pos: Pos{X: 10, Y: 10}, Width: 4, Height: 3, Name: "Depot"}
}

func TestStation_Adjoins(t *testing.T) {
	s := testStation() // footprint: x=10..13, y=10..12

	t.Run("tile inside the footprint is not adjoining", func(t *testing.T) {
		assert.False(t, s.Adjoins(Pos{X: 11, Y: 11}))
	})

	t.Run("tile directly bordering an edge is adjoining", func(t *testing.T) {
		assert.True(t, s.Adjoins(Pos{X: 10, Y: 9}))  // north of footprint
		assert.True(t, s.Adjoins(Pos{X: 10, Y: 13})) // south of footprint
		assert.True(t, s.Adjoins(Pos{X: 9, Y: 10}))  // west of footprint
		assert.True(t, s.Adjoins(Pos{X: 14, Y: 10})) // east of footprint
	})

	t.Run("tile touching only a corner is adjoining", func(t *testing.T) {
		assert.True(t, s.Adjoins(Pos{X: 9, Y: 9}))
		assert.True(t, s.Adjoins(Pos{X: 14, Y: 13}))
	})

	t.Run("tile two tiles away is not adjoining", func(t *testing.T) {
		assert.False(t, s.Adjoins(Pos{X: 8, Y: 10}))
		assert.False(t, s.Adjoins(Pos{X: 10, Y: 8}))
	})
}

func TestStationAt(t *testing.T) {
	a := testStation()
	b := Station{Pos: Pos{X: 100, Y: 100}, Width: 2, Height: 2, Name: "Yard"}
	stations := []*Station{&a, &b}

	t.Run("returns the station a position is adjacent to", func(t *testing.T) {
		found := StationAt(Pos{X: 9, Y: 10}, stations)
		assert.Equal(t, &a, found)
	})

	t.Run("returns nil when no station is adjacent", func(t *testing.T) {
		assert.Nil(t, StationAt(Pos{X: 0, Y: 0}, stations))
	})
}
