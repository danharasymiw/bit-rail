package trains

import "github.com/danharasymiw/bit-rail/types"

type trainWorldView interface {
	TileAt(pos types.Pos) *types.Tile
	TrackAt(pos types.Pos) *types.Track
	OccupiedAt(pos types.Pos) bool
	SetOccupied(pos types.Pos)
	UnsetOccupied(pos types.Pos)
	// NextHop returns the direction to take from a junction at pos to make
	// progress toward dest.
	NextHop(pos, dest types.Pos) (types.Dir, bool)
	// StationAt returns the station (if any) that pos is adjacent to.
	StationAt(pos types.Pos) *types.Station
}
