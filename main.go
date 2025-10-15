package main

import (
	"time"

	"github.com/danharasymiw/trains/engine"
	"github.com/danharasymiw/trains/renderer"
	"github.com/danharasymiw/trains/world/test_worlds"

	"github.com/gdamore/tcell"
)

func main() {
	w := test_worlds.NewPerlinWorld(123, 123)
	screen, _ := tcell.NewScreen()
	screen.Init()
	defer screen.Fini()

	eng := engine.New(
		w,
		renderer.NewSimpleRenderer(screen, w),
		150*time.Millisecond,
	)
	eng.Run()
	println("bye.")
}
