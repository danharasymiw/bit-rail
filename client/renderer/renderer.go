package renderer

import (
	"github.com/danharasymiw/trains/trains"
	"github.com/gdamore/tcell"
)

type Renderer interface {
	RenderRegion(x, y, width, height int)
	RenderTrains([]*trains.Train)
	Screen() tcell.Screen // Would rather not have this coupled with tcell... but its just to stop server
}
