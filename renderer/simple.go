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
			ch, col := r.getTileChar(x, y, t)
			style := tcell.StyleDefault.Foreground(col)
			r.screen.SetContent(x, y, ch, nil, style)
		}
	}

	for _, t := range r.w.Trains {
		for _, c := range t.Cars {
			ch, col := r.getTrainCarChar(c)
			style := tcell.StyleDefault.Foreground(col)
			r.screen.SetContent(c.X, c.Y, ch, nil, style)
		}
	}
}

var (
	grassChars  = []rune{'.', ',', '\'', '`', ':'}
	grassColors = []tcell.Color{
		tcell.ColorGreen,
		tcell.ColorDarkGreen,
		tcell.ColorLime,
	}

	waterChars  = []rune{'~', '≈', '-', '`', '"'}
	waterColors = []tcell.Color{
		tcell.ColorBlue,
		tcell.ColorBlue,
		tcell.ColorSteelBlue,
	}
)

func (r *SimpleRenderer) getTileChar(x, y int, t *types.Tile) (rune, tcell.Color) {
	switch t.Type {
	case types.TileGrass:
		return grassChars[(x^y)%len(grassChars)], grassColors[(x^y)%len(grassColors)]
	case types.TileWater:
		return waterChars[(x^y)%len(grassChars)], waterColors[(x^y)%len(waterColors)]
	case types.TileWood:
		return 'x', tcell.ColorBrown
	case types.TileTrack:
		var trackChar rune
		track := r.w.Tracks[t]
		switch track.Direction {
		case types.DirNorth | types.DirSouth:
			trackChar = '║' // vertical
		case types.DirEast | types.DirWest:
			trackChar = '═' // horizontal
		case types.DirNorth | types.DirEast:
			trackChar = '╚' // curve NE
		case types.DirNorth | types.DirWest:
			trackChar = '╝' // curve NW
		case types.DirSouth | types.DirEast:
			trackChar = '╔' // curve SE
		case types.DirSouth | types.DirWest:
			trackChar = '╗' // curve SW
		case types.DirNorth | types.DirEast | types.DirWest:
			trackChar = '╩' // T junction pointing up
		case types.DirSouth | types.DirEast | types.DirWest:
			trackChar = '╦' // T junction pointing down
		case types.DirNorth | types.DirSouth | types.DirEast:
			trackChar = '╠' // T junction pointing left
		case types.DirNorth | types.DirSouth | types.DirWest:
			trackChar = '╣' // T junction pointing right
		case types.DirNorth | types.DirSouth | types.DirEast | types.DirWest:
			trackChar = '╬' // cross
		default:
			trackChar = ' '
		}
		return trackChar, tcell.ColorGray
	}
	return ' ', tcell.ColorRed
}

func (r *SimpleRenderer) getTrainCarChar(c *trains.TrainCar) (rune, tcell.Color) {
	switch c.Type {
	case trains.CarTypeLocomotive:
		return 'H', tcell.ColorRed
	case trains.CarTypeCargo:
		return 'O', tcell.ColorOrange
	default:
		return 'X', tcell.ColorRed
	}
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
