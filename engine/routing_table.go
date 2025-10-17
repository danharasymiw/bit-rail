package engine

import "github.com/danharasymiw/trains/types"

type routingEntry struct {
	blockID types.BlockID
	dir     types.Dir
	cost    int
	dirty   bool // do we need this? maybe entry should just be removed.
}

type routingTable struct {
	entries map[types.NodeID]routingEntry
}
