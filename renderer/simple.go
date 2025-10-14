package renderer

import (
	"github.com/danharasymiw/trains/log"
	"github.com/danharasymiw/trains/trains"
	"github.com/danharasymiw/trains/types"
	"github.com/danharasymiw/trains/world"
	"github.com/gdamore/tcell"
)

type SimpleRenderer struct {
	screen tcell.Screen
	w      *world.World
}

func NewSimpleRenderer(screen tcell.Screen, w *world.World) *SimpleRenderer {
	return &SimpleRenderer{
		screen: screen,
		w:      w,
	}
}

func (r *SimpleRenderer) Draw() {
	r.drawWorld()
	r.drawLogs(30)
	r.screen.Show()
}

func (r *SimpleRenderer) drawWorld() {
	for y, row := range r.w.Tiles {
		for x, t := range row {
			ch := r.getTileChar(t)
			r.screen.SetContent(x, y, ch, nil, tcell.StyleDefault)
		}
	}

	for _, t := range r.w.Trains {
		for _, c := range t.Cars {
			var ch rune
			switch c.Type {
			case trains.CarTypeLocomotive:
				ch = 'L'
			case trains.CarTypeCargo:
				ch = 'C'
			case trains.CarTypePassenger:
				ch = 'C'
			}
			r.screen.SetContent(c.X, c.Y, ch, nil, tcell.StyleDefault)
		}
	}
}

func (r *SimpleRenderer) getTileChar(t *types.Tile) rune {
	switch t.Type {
	case types.TileGrass:
		return ' '
	case types.TileWood:
		return '🌲'
	case types.TileTrack:
		track := r.w.Tracks[t]
		switch track.Direction {
		case types.DirNorth | types.DirSouth:
			return '║' // vertical
		case types.DirEast | types.DirWest:
			return '═' // horizontal
		case types.DirNorth | types.DirEast:
			return '╚' // curve NE
		case types.DirNorth | types.DirWest:
			return '╝' // curve NW
		case types.DirSouth | types.DirEast:
			return '╔' // curve SE
		case types.DirSouth | types.DirWest:
			return '╗' // curve SW
		case types.DirNorth | types.DirEast | types.DirWest:
			return '╩' // T junction pointing up
		case types.DirSouth | types.DirEast | types.DirWest:
			return '╦' // T junction pointing down
		case types.DirNorth | types.DirSouth | types.DirEast:
			return '╠' // T junction pointing left
		case types.DirNorth | types.DirSouth | types.DirWest:
			return '╣' // T junction pointing right
		case types.DirNorth | types.DirSouth | types.DirEast | types.DirWest:
			return '╬' // cross
		}
	}
	return ' ' // unrecognized/unsupported
}

func (r *SimpleRenderer) drawLogs(startY int) {
	logs := log.GetLogs()
	for i, msg := range logs {
		for x, c := range msg {
			r.screen.SetContent(x, startY+i, c, nil, tcell.StyleDefault)
		}
		for x := len(msg); x < 80; x++ {
			r.screen.SetContent(x, startY+i, ' ', nil, tcell.StyleDefault)
		}
	}
}

func (r *SimpleRenderer) Screen() tcell.Screen {
	return r.screen
}
