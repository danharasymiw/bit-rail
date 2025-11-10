package types

import "strings"

type Pos struct {
	X int `binary:"uint16"`
	Y int `binary:"uint16"`
}

type Dir uint8

const (
	DirNone  Dir = 0
	DirNorth     = 1 << 0 // 0001
	DirEast      = 1 << 1 // 0010
	DirSouth     = 1 << 2 // 0100
	DirWest      = 1 << 3 // 1000

	// Really only used for manually writing tracks for testing
	DirNorthSouth = DirNorth | DirSouth
	DirNorthWest  = DirNorth | DirWest
	DirNorthEast  = DirNorth | DirEast
	DirSouthWest  = DirSouth | DirWest
	DirSouthEast  = DirSouth | DirEast
	DirEastWest   = DirEast | DirWest
	DirFourWay    = DirNorthSouth | DirEastWest
)

func (d Dir) String() string {
	dirs := make([]string, 0, 4)
	if d&DirNorth != 0 {
		dirs = append(dirs, "North")
	}
	if d&DirEast != 0 {
		dirs = append(dirs, "East")
	}
	if d&DirSouth != 0 {
		dirs = append(dirs, "South")
	}
	if d&DirWest != 0 {
		dirs = append(dirs, "West")
	}
	return strings.Join(dirs, ", ")
}

func OppositeDir(d Dir) Dir {
	switch d {
	case DirNorth:
		return DirSouth
	case DirSouth:
		return DirNorth
	case DirEast:
		return DirWest
	case DirWest:
		return DirEast
	default:
		return DirNone
	}
}
