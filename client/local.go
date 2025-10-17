package client

import (
	"time"

	"github.com/danharasymiw/trains/client/renderer"
	"github.com/danharasymiw/trains/trains"
	"github.com/danharasymiw/trains/world"
	"github.com/gdamore/tcell"
)

type localClient struct {
	w       *world.World
	r       renderer.Renderer
	quitCh  chan bool
	running bool

	camX, camY int
}

func NewLocal(w *world.World) (*localClient, chan bool) {
	quitCh := make(chan bool)
	return &localClient{
		w:       w,
		quitCh:  quitCh,
		running: false,
	}, quitCh
}

func (c *localClient) Run() {
	screen, err := tcell.NewScreen()
	if err != nil {
		panic(err)
	}
	if err := screen.Init(); err != nil {
		panic(err)
	}
	defer screen.Fini()

	c.r = renderer.NewSimpleRenderer(screen, c.w)
	c.running = true

	// Buffered event channel to receive user input
	events := make(chan tcell.Event, 32)

	// Poll events in background
	go func() {
		for {
			ev := screen.PollEvent()
			if ev == nil {
				return
			}
			events <- ev
		}
	}()

	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()

	for c.running {
		select {
		case ev := <-events:
			switch tev := ev.(type) {
			case *tcell.EventKey:
				switch tev.Key() {
				case tcell.KeyUp:
					c.camY--
				case tcell.KeyDown:
					c.camY++
				case tcell.KeyLeft:
					c.camX--
				case tcell.KeyRight:
					c.camX++
				}
				if tev.Rune() == 'q' {
					c.running = false
				}
			case *tcell.EventResize:
				screen.Sync()
			}

		case <-ticker.C:
			// render game state periodically
			c.r.RenderRegion(c.camX, c.camY, 120, 90)
			c.r.RenderTrains([]*trains.Train{})
			c.r.Screen().Show()
		}
	}

	// Tell whoever launched us that we're done
	c.quitCh <- true
}
