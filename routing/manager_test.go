package routing

import (
	"testing"

	"github.com/danharasymiw/bit-rail/types"
	"github.com/danharasymiw/bit-rail/world"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestWorld() *world.World {
	return world.New(32, 32)
}

// buildJunctionWorld creates a 4-way junction at (2,2) with a 2-tile arm in
// each direction. The north arm ends at a track tile adjacent to a station;
// the other three arms end in a (semantically arbitrary) signal tile purely
// so each arm terminates in a real node and produces an edge.
//
//	           [North Station] (2,5)
//	           (2,4) <- adjacent to the station above
//	           (2,3)
//	(0,2)-(1,2)-(2,2)J-(3,2)-(4,2)
//	           (2,1)
//	           (2,0)
func buildJunctionWorld() (*world.World, types.Pos, types.Pos) {
	w := newTestWorld()

	junctionPos := types.Pos{X: 2, Y: 2}
	w.AddTrack(junctionPos, &types.Track{Direction: types.DirFourWay})

	stationPos := types.Pos{X: 2, Y: 4}
	w.AddTrack(types.Pos{X: 2, Y: 3}, &types.Track{Direction: types.DirNorthSouth})
	w.AddTrack(stationPos, &types.Track{Direction: types.DirSouth})
	w.AddStation(&types.Station{Pos: types.Pos{X: 2, Y: 5}, Width: 1, Height: 1, Name: "North Station"})

	w.AddTrack(types.Pos{X: 2, Y: 1}, &types.Track{Direction: types.DirNorthSouth})
	w.AddTrack(types.Pos{X: 2, Y: 0}, &types.Track{Direction: types.DirNorth, SignalDir: types.DirNorth})

	w.AddTrack(types.Pos{X: 3, Y: 2}, &types.Track{Direction: types.DirEastWest})
	w.AddTrack(types.Pos{X: 4, Y: 2}, &types.Track{Direction: types.DirWest, SignalDir: types.DirWest})

	w.AddTrack(types.Pos{X: 1, Y: 2}, &types.Track{Direction: types.DirEastWest})
	w.AddTrack(types.Pos{X: 0, Y: 2}, &types.Track{Direction: types.DirEast, SignalDir: types.DirEast})

	return w, junctionPos, stationPos
}

func TestManager_Rebuild_DiscoversNodes(t *testing.T) {
	w, junctionPos, stationPos := buildJunctionWorld()

	m := NewManager()
	m.Rebuild(w)

	require.Contains(t, m.Nodes, junctionPos)
	assert.Equal(t, NodeJunction, m.Nodes[junctionPos].Kind)

	require.Contains(t, m.Nodes, stationPos)
	assert.Equal(t, NodeStation, m.Nodes[stationPos].Kind)
	assert.Equal(t, "North Station", m.Nodes[stationPos].Name)

	// Plain pass-through track tiles (e.g. (2,3), the tile between the
	// junction and the station) are never nodes.
	assert.NotContains(t, m.Nodes, types.Pos{X: 2, Y: 3})
}

func TestManager_Rebuild_EdgesHaveCorrectCostAndDirection(t *testing.T) {
	w, junctionPos, stationPos := buildJunctionWorld()

	m := NewManager()
	m.Rebuild(w)

	edges := m.edges[junctionPos]
	require.Len(t, edges, 4, "junction should have one edge per arm")

	var north *edge
	for i := range edges {
		if edges[i].To == stationPos {
			north = &edges[i]
		}
	}
	require.NotNil(t, north, "junction should have an edge reaching the station")
	assert.Equal(t, types.Dir(types.DirNorth), north.Dir)
	assert.Equal(t, 2, north.Cost, "station is 2 tiles north of the junction")
}

func TestManager_NextHop_RoutesTowardStation(t *testing.T) {
	w, junctionPos, stationPos := buildJunctionWorld()

	m := NewManager()
	m.Rebuild(w)

	dir, ok := m.NextHop(junctionPos, stationPos)
	require.True(t, ok, "should find a route from the junction to the station")
	assert.Equal(t, types.Dir(types.DirNorth), dir)
}

func TestManager_NextHop_UnknownDestinationReturnsFalse(t *testing.T) {
	w, junctionPos, _ := buildJunctionWorld()

	m := NewManager()
	m.Rebuild(w)

	_, ok := m.NextHop(junctionPos, types.Pos{X: 99, Y: 99})
	assert.False(t, ok)
}

func TestManager_NextHop_SamePositionReturnsFalse(t *testing.T) {
	w, junctionPos, _ := buildJunctionWorld()

	m := NewManager()
	m.Rebuild(w)

	_, ok := m.NextHop(junctionPos, junctionPos)
	assert.False(t, ok)
}

// TestManager_Rebuild_LoopWithNoNodesHasNoEdges verifies a closed loop of
// plain track with no signals/junctions/stations produces no nodes at all
// (and thus the walk's cycle-detection never needs to answer a query).
func TestManager_Rebuild_LoopWithNoNodesHasNoEdges(t *testing.T) {
	w := newTestWorld()

	w.AddTrack(types.Pos{X: 0, Y: 0}, &types.Track{Direction: types.DirSouthEast})
	w.AddTrack(types.Pos{X: 1, Y: 0}, &types.Track{Direction: types.DirEastWest})
	w.AddTrack(types.Pos{X: 2, Y: 0}, &types.Track{Direction: types.DirSouthWest})
	w.AddTrack(types.Pos{X: 2, Y: 1}, &types.Track{Direction: types.DirNorthWest})
	w.AddTrack(types.Pos{X: 1, Y: 1}, &types.Track{Direction: types.DirEastWest})
	w.AddTrack(types.Pos{X: 0, Y: 1}, &types.Track{Direction: types.DirNorthEast})

	m := NewManager()
	m.Rebuild(w)

	assert.Empty(t, m.Nodes)
	assert.Empty(t, m.edges)
}
