package types

// Station is a structure placed near track (not on it) representing a
// delivery destination. Pos is the top-left corner of its footprint.
type Station struct {
	Pos    Pos
	Width  int    `binary:"uint16"`
	Height int    `binary:"uint16"`
	Name   string
}

// Adjoins reports whether pos is a tile immediately bordering (not inside)
// the station's footprint — i.e. a track tile a train could pull up
// alongside to load/unload at this station.
func (s Station) Adjoins(pos Pos) bool {
	inside := pos.X >= s.Pos.X && pos.X < s.Pos.X+s.Width &&
		pos.Y >= s.Pos.Y && pos.Y < s.Pos.Y+s.Height
	if inside {
		return false
	}

	minX, maxX := s.Pos.X-1, s.Pos.X+s.Width
	minY, maxY := s.Pos.Y-1, s.Pos.Y+s.Height
	return pos.X >= minX && pos.X <= maxX && pos.Y >= minY && pos.Y <= maxY
}

// StationAt returns the station (if any) that pos is adjacent to.
func StationAt(pos Pos, stations []*Station) *Station {
	for _, s := range stations {
		if s.Adjoins(pos) {
			return s
		}
	}
	return nil
}
