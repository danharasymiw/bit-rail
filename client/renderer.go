package client

import (
	"github.com/danharasymiw/bit-rail/trains"
	"github.com/danharasymiw/bit-rail/types"
	"github.com/danharasymiw/bit-rail/world"
	"github.com/gdamore/tcell"
)

type ChatMessage struct {
	Author  string
	Message string
}

type Renderer interface {
	Render(camX, camY int, chatMessages []ChatMessage)
	Screen() tcell.Screen
}

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

func (r *SimpleRenderer) Render(camX, camY int, chatMessages []ChatMessage) {
	termWidth, termHeight := r.screen.Size()

	infoPanelWidth := 35
	chatPanelHeight := 10
	worldWidth := termWidth - infoPanelWidth
	worldHeight := termHeight - chatPanelHeight

	r.renderRegion(camX, camY, worldWidth, worldHeight)
	r.renderTrains(camX, camY, worldWidth, worldHeight)
	r.renderInfoPanel(worldWidth, 0, infoPanelWidth, worldHeight)
	r.renderChatPanel(0, worldHeight, termWidth, chatPanelHeight, chatMessages)

	r.screen.Show()
}

func (r *SimpleRenderer) renderRegion(x, y, width, height int) {
	for relY, row := range r.w.Tiles[y : y+height] {
		for relX, t := range row[x : x+width] {
			worldX := x + relX
			worldY := y + relY
			ch, style := r.getTileChar(worldX, worldY, t)
			screenY := height - 1 - relY // Flip Y
			r.screen.SetContent(relX, screenY, ch, nil, style)
		}
	}
}

func (r *SimpleRenderer) renderTrains(x, y, width, height int) {
	for _, t := range r.w.Trains {
		// Assuming train limits of 100 - check the first car to see if its
		// even possible to be on screen
		if len(t.Cars) > 0 {
			c := t.Cars[0]
			if c.X < x-100 || c.X >= x+width+100 || c.Y < y-100 || c.Y >= y+height+100 {
				continue // Skip this train
			}
		}
		for _, c := range t.Cars {
			if c.X < x || c.X >= x+width || c.Y < y || c.Y >= y+height {
				continue // Skip this car
			}

			ch, col := r.getTrainCarChar(c)
			style := tcell.StyleDefault.Foreground(col)
			screenX := c.X - x
			screenY := height - 1 - (c.Y - y)

			r.screen.SetContent(screenX, screenY, ch, nil, style)
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

func (r *SimpleRenderer) renderInfoPanel(x, y, width, height int) {
	borderStyle := tcell.StyleDefault.Foreground(tcell.ColorWhite)

	// Draw border
	for i := 0; i < height; i++ {
		r.screen.SetContent(x, y+i, '│', nil, borderStyle)
	}

	// Title
	title := " Info "
	for i, ch := range title {
		r.screen.SetContent(x+1+i, y, ch, nil, borderStyle)
	}

	// Clear content area
	for py := y + 1; py < y+height; py++ {
		for px := x + 1; px < x+width; px++ {
			r.screen.SetContent(px, py, ' ', nil, tcell.StyleDefault)
		}
	}

	// TODO: Add actual info content here

}

func (r *SimpleRenderer) renderChatPanel(x, y, width, height int, chatMessages []ChatMessage) {
	borderStyle := tcell.StyleDefault.Foreground(tcell.ColorWhite)

	// Draw top border
	for i := 0; i < width; i++ {
		r.screen.SetContent(x+i, y, '─', nil, borderStyle)
	}

	// Title
	title := " Chat "
	for i, ch := range title {
		r.screen.SetContent(x+2+i, y, ch, nil, borderStyle)
	}

	// Clear content area
	for py := y + 1; py < y+height; py++ {
		for px := x; px < x+width; px++ {
			r.screen.SetContent(px, py, ' ', nil, tcell.StyleDefault)
		}
	}

	// Render chat messages (bottom-up, most recent at bottom)
	availableHeight := height - 1 // Subtract border
	startIdx := 0
	if len(chatMessages) > availableHeight {
		startIdx = len(chatMessages) - availableHeight
	}

	lineY := y + 1
	for i := startIdx; i < len(chatMessages) && lineY < y+height; i++ {
		msg := chatMessages[i]

		// Format: [Author] Message
		var displayText string
		if msg.Author != "" {
			displayText = "[" + msg.Author + "] " + msg.Message
		} else {
			displayText = msg.Message
		}

		// Truncate if too long
		if len(displayText) > width-2 {
			displayText = displayText[:width-2]
		}

		// Render the message
		msgStyle := tcell.StyleDefault.Foreground(tcell.ColorWhite)
		for col, ch := range displayText {
			r.screen.SetContent(x+1+col, lineY, ch, nil, msgStyle)
		}

		lineY++
	}
}
