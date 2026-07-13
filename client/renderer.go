package client

import (
	"github.com/gdamore/tcell"
	"github.com/google/uuid"

	"github.com/danharasymiw/bit-rail/trains"
	"github.com/danharasymiw/bit-rail/types"
	"github.com/danharasymiw/bit-rail/world"
)

type ChatMessage struct {
	Author  string
	Message string
}

type Renderer interface {
	Render(camPos types.Pos, chatMessages []ChatMessage)
	Screen() tcell.Screen
}

type SimpleRenderer struct {
	screen    tcell.Screen
	w         *world.World
	debugMode bool
}

func NewSimpleRenderer(screen tcell.Screen, w *world.World, debugMode bool) *SimpleRenderer {
	return &SimpleRenderer{
		screen:    screen,
		w:         w,
		debugMode: debugMode,
	}
}

func (r *SimpleRenderer) Screen() tcell.Screen {
	return r.screen
}

func (r *SimpleRenderer) Render(camPos types.Pos, chatMessages []ChatMessage) {
	termWidth, termHeight := r.screen.Size()

	infoPanelWidth := 35
	chatPanelHeight := 10
	worldWidth := termWidth - infoPanelWidth
	worldHeight := termHeight - chatPanelHeight

	r.renderRegion(camPos, worldWidth, worldHeight)
	r.renderStations(camPos, worldWidth, worldHeight)
	r.renderTrains(camPos, worldWidth, worldHeight)
	r.renderInfoPanel(worldWidth, 0, infoPanelWidth, worldHeight)
	r.renderChatPanel(0, worldHeight, termWidth, chatPanelHeight, chatMessages)

	r.screen.Show()
}

func (r *SimpleRenderer) renderRegion(pos types.Pos, width, height int) {
	for relY := range height {
		worldY := pos.Y + relY
		if worldY >= len(r.w.Tiles) {
			// Fill remaining rows with spaces if we run out of tiles
			for relX := range width {
				screenY := height - 1 - relY
				r.screen.SetContent(relX, screenY, ' ', nil, tcell.StyleDefault)
			}
			continue
		}
		row := r.w.Tiles[worldY]
		for relX := range width {
			worldX := pos.X + relX
			if worldX >= len(row) {
				// Fill with space if we run out of tiles in this row
				screenY := height - 1 - relY
				r.screen.SetContent(relX, screenY, ' ', nil, tcell.StyleDefault)
				continue
			}
			tile := row[worldX]
			if tile == nil {
				// Fill with space if tile is nil
				screenY := height - 1 - relY
				r.screen.SetContent(relX, screenY, ' ', nil, tcell.StyleDefault)
				continue
			}
			worldPos := types.Pos{X: worldX, Y: worldY}
			ch, style := r.getTileChar(worldPos, tile)
			screenY := height - 1 - relY // Flip Y
			r.screen.SetContent(relX, screenY, ch, nil, style)
		}
	}
}

func (r *SimpleRenderer) renderTrains(pos types.Pos, width, height int) {
	posX, posY := pos.X, pos.Y
	for _, t := range r.w.Trains {
		// Assuming train limits of 100 - check the first car to see if its
		// even possible to be on screen
		if len(t.Cars) > 0 {
			c := t.Cars[0]
			if c.X < posX-100 || c.X >= posX+width+100 || c.Y < posY-100 || c.Y >= posY+height+100 {
				continue // Skip this train
			}
		}

		for _, c := range t.Cars {
			if c.X < posX || c.X >= posX+width || c.Y < posY || c.Y >= posY+height {
				continue // Skip this car
			}

			ch, col := r.getTrainCarChar(c)
			style := tcell.StyleDefault.Foreground(col)
			screenX := c.X - pos.X
			screenY := height - 1 - (c.Y - posY)

			r.screen.SetContent(screenX, screenY, ch, nil, style)
		}
	}
}

// renderStations draws each station's ASCII footprint on top of the
// terrain/track underneath it. Stations aren't on the tracks themselves —
// track tiles bordering a station's footprint are what count as "at" it
// (see types.Station.Adjoins).
func (r *SimpleRenderer) renderStations(camPos types.Pos, width, height int) {
	for _, s := range r.w.Stations {
		art := stationArt(s)
		for row, line := range art {
			// art[0] is the top of the building (the roof); increasing
			// world Y is north/"up" on screen (see renderRegion's Y-flip),
			// so the roof belongs at the highest world Y in the footprint.
			worldY := s.Pos.Y + (s.Height - 1 - row)
			screenY := height - 1 - (worldY - camPos.Y)
			if screenY < 0 || screenY >= height {
				continue
			}
			for col, cell := range line {
				if cell == nil {
					continue // transparent: leave the terrain/track showing
				}
				worldX := s.Pos.X + col
				screenX := worldX - camPos.X
				if screenX < 0 || screenX >= width {
					continue
				}
				style := tcell.StyleDefault.Foreground(cell.fg).Background(cell.bg)
				r.screen.SetContent(screenX, screenY, cell.ch, nil, style)
			}
		}
	}
}

type stationCell struct {
	ch     rune
	fg, bg tcell.Color
}

var (
	stationFillBG   = tcell.ColorTan
	stationRoofBG   = tcell.ColorFireBrick
	stationBorderFG = tcell.ColorBlack
)

func filled(ch rune, fg, bg tcell.Color) *stationCell {
	return &stationCell{ch: ch, fg: fg, bg: bg}
}

func painted(bg tcell.Color) *stationCell {
	return &stationCell{ch: ' ', fg: bg, bg: bg}
}

// stationArt draws a station as a plain solid rectangle: a ┌─┐/└─┘ border
// (matching the box-drawing style used elsewhere for panels), a roof-colored
// band on the top interior row, and a single flat fill color below that.
// Deliberately minimal — just enough to read as a distinct structure next
// to the track, without competing for attention.
func stationArt(s *types.Station) [][]*stationCell {
	art := make([][]*stationCell, s.Height)
	for row := range art {
		art[row] = make([]*stationCell, s.Width)
		for col := range art[row] {
			if ch, ok := stationBorderRune(s, row, col); ok {
				art[row][col] = filled(ch, stationBorderFG, stationFillBG)
			} else if row == 1 && s.Height > 2 {
				art[row][col] = painted(stationRoofBG)
			} else {
				art[row][col] = painted(stationFillBG)
			}
		}
	}
	return art
}

func stationBorderRune(s *types.Station, row, col int) (rune, bool) {
	switch {
	case row == 0 && col == 0:
		return '┌', true
	case row == 0 && col == s.Width-1:
		return '┐', true
	case row == s.Height-1 && col == 0:
		return '└', true
	case row == s.Height-1 && col == s.Width-1:
		return '┘', true
	case row == 0 || row == s.Height-1:
		return '─', true
	case col == 0 || col == s.Width-1:
		return '│', true
	}
	return 0, false
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

func (r *SimpleRenderer) getTileChar(pos types.Pos, t *types.Tile) (rune, tcell.Style) {
	var ch rune
	var fgCol tcell.Color
	switch t.Type {
	case types.TileGrass:
		ch = grassChars[(pos.X^pos.Y)%len(grassChars)]
		fgCol = grassColors[(pos.X^pos.Y)%len(grassColors)]
	case types.TileWater:
		ch = waterChars[(pos.X^pos.Y)%len(waterChars)]
		fgCol = waterColors[(pos.X^pos.Y)%len(waterColors)]
	case types.TileTree:
		ch = treeChars[(pos.X^pos.Y)%len(treeChars)]
		fgCol = treeColors[(pos.X^pos.Y)%len(treeColors)]
	case types.TileMountain:
		ch = mountainChars[(pos.X^pos.Y)%len(mountainChars)]
		fgCol = mountainColors[(pos.X^pos.Y)%len(mountainColors)]

	case types.TileTrack:
		track := r.w.Tracks[pos]

		// If track has signals, render signal symbols instead of track
		if track.HasSignal() {
			// Render signal for the first direction found (prioritize north, south, east, west)
			var signalChar rune
			var signalColor tcell.Color

			if track.SignalDir&types.DirNorth != 0 {
				signalChar = '^'
			} else if track.SignalDir&types.DirSouth != 0 {
				signalChar = 'v'
			} else if track.SignalDir&types.DirEast != 0 {
				signalChar = '>'
			} else if track.SignalDir&types.DirWest != 0 {
				signalChar = '<'
			}

			if track.Block.OccupiedBy == uuid.Nil {
				signalColor = tcell.ColorGreen
			} else {
				signalColor = tcell.ColorRed
			}

			return signalChar, tcell.StyleDefault.Foreground(signalColor)
		}

		// No signal, render track character
		var trackChar rune
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

		// In debug mode, color tracks by their block ID
		if r.debugMode {
			if track.Block != nil && track.Block.ID != uuid.Nil {
				// Block has an ID - use color derived from UUID
				blockColor := colorFromUUID(&track.Block.ID)
				return trackChar, tcell.StyleDefault.Foreground(blockColor)
			} else if track.Block == nil {
				// Track has no block assigned - draw as white
				return trackChar, tcell.StyleDefault.Foreground(tcell.ColorWhite)
			}
		}
		// Default: gray (or when not in debug mode)
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
	for i := range width {
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

// colorFromUUID generates a consistent color from a UUID by using its bytes
// to create RGB values. This ensures each block gets a unique, stable color.
func colorFromUUID(id *uuid.UUID) tcell.Color {
	// UUID is 16 bytes, we'll use them to generate RGB
	// Use first 3 bytes for R, G, B, ensuring minimum brightness for visibility
	bytes := id[:]

	// Extract RGB values, ensuring they're bright enough to be visible
	// We'll map 0-255 to a range that's visible but not too dark
	// Using a range of 80-255 to ensure good visibility
	r := int(bytes[0])%176 + 80 // 0-255 -> 80-255
	g := int(bytes[1])%176 + 80 // 0-255 -> 80-255
	b := int(bytes[2])%176 + 80 // 0-255 -> 80-255

	return tcell.NewRGBColor(int32(r), int32(g), int32(b))
}
