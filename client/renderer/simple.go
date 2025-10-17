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

func (r *SimpleRenderer) Screen() tcell.Screen {
	return r.screen
}

func (r *SimpleRenderer) Draw() {
	r.drawLogs(30)
	r.screen.Show()
}

func (r *SimpleRenderer) RenderRegion(x, y, width, height int) {
	for y, row := range r.w.Tiles[y : y+height] {
		for x, t := range row[x : x+width] {
			ch, style := r.getTileChar(x, y, t)
			r.screen.SetContent(x, y, ch, nil, style)
		}
	}
}

func (r *SimpleRenderer) RenderTrains(trains []*trains.Train) {
	for _, t := range r.w.Trains {
		for _, c := range t.Cars {
			ch, col := r.getTrainCarChar(c)
			style := tcell.StyleDefault.Foreground(col)
			r.screen.SetContent(c.X, c.Y, ch, nil, style)
		}
	}
}

var (
	grassChars  = []rune(".,'`:")
	grassColors = []tcell.Color{
		tcell.ColorYellowGreen,
		tcell.ColorLightGreen,
		tcell.ColorLawnGreen,
	}

	treeChars  = []rune("TtYy")
	treeColors = []tcell.Color{
		tcell.ColorDarkGreen,
		tcell.ColorOliveDrab,
		tcell.ColorForestGreen,
	}

	waterChars  = []rune("~≈-`")
	waterColors = []tcell.Color{
		tcell.ColorBlue,
		tcell.ColorSteelBlue,
		tcell.ColorDeepSkyBlue,
	}

	mountainChars  = []rune("^M")
	mountainColors = []tcell.Color{
		tcell.ColorSlateGray,
		tcell.ColorDarkGray,
		tcell.ColorDimGray,
	}
)

func (r *SimpleRenderer) getTileChar(x, y int, t *types.Tile) (rune, tcell.Style) {
	var ch rune
	var fgCol tcell.Color
	switch t.Type {
	case types.TileGrass:
		ch = grassChars[(x^y)%len(grassChars)]
		fgCol = grassColors[(x^y)%len(grassColors)]
	case types.TileWater:
		ch = waterChars[(x^y)%len(waterChars)]
		fgCol = waterColors[(x^y)%len(waterColors)]
	case types.TileTree:
		ch = treeChars[(x^y)%len(treeChars)]
		fgCol = treeColors[(x^y)%len(treeColors)]
	case types.TileMountain:
		ch = mountainChars[(x^y)%len(mountainChars)]
		fgCol = mountainColors[(x^y)%len(mountainColors)]

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
		return trackChar, tcell.StyleDefault.Foreground(tcell.ColorGray)
	}
	return ch, tcell.StyleDefault.Foreground(fgCol)
}

func (r *SimpleRenderer) getTrainCarChar(c *trains.TrainCar) (rune, tcell.Color) {
	switch c.Type {
	case trains.CarTypeLocomotive:
		return '█', tcell.ColorRed
	case trains.CarTypeCargo:
		return '▓', tcell.ColorSilver
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
