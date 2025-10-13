package renderer

import "github.com/gdamore/tcell"

type Renderer interface {
	Draw()
	Screen() tcell.Screen // Would rather not have this coupled with tcell... but its just to stop server
}
