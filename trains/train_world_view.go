package trains

import "github.com/danharasymiw/bit-rail/types"

type trainWorldView interface {
	TileAt(pos types.Pos) *types.Tile
	TrackAt(pos types.Pos) *types.Track
	OccupiedAt(pos types.Pos) bool
	SetOccupied(pos types.Pos)
	UnsetOccupied(pos types.Pos)
}
