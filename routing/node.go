package routing

import "github.com/danharasymiw/bit-rail/types"

// NodeKind classifies why a track tile is a routing decision point.
type NodeKind uint8

const (
	NodeSignal NodeKind = iota
	NodeJunction
	NodeStation
)

// Node is a routing decision point: a signal, junction, or station.
// Blocks (contiguous runs of plain track) connect nodes together.
type Node struct {
	Pos  types.Pos
	Kind NodeKind
	// Name is the station name; empty for signal/junction nodes.
	Name string
}

// classify determines the node kind for a track tile that isn't adjacent to
// a station (station adjacency is checked separately — see Manager.Rebuild
// — since it depends on world-level Station data, not the Track itself).
func classify(t *types.Track) NodeKind {
	if t.IsJunction() {
		return NodeJunction
	}
	return NodeSignal
}
