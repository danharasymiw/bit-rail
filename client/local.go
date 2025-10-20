package client

import (
	"time"

	"github.com/danharasymiw/trains/world"
	"github.com/gdamore/tcell"
)

type localClient struct {
	w       *world.World
	r       Renderer
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

	c.r = NewSimpleRenderer(screen, c.w)
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
					_, height := c.r.Screen().Size()
					if c.camY < c.w.Height - height {
						c.camY++
					}
				case tcell.KeyDown:
					if c.camY > 0 {
						c.camY--
					}
				case tcell.KeyLeft:
					if c.camX > 0 {
						c.camX--
					}
				case tcell.KeyRight:
					width, _ := c.r.Screen().Size()
					if c.camX < c.w.Width - width {
						c.camX+=2
					}
				}
				if tev.Rune() == 'q' {
					c.running = false
				}
			case *tcell.EventResize:
				screen.Sync()
			}

		case <-ticker.C:
			c.r.Render(c.camX, c.camY)
		}
	}

	// Tell whoever launched us that we're done
	c.quitCh <- true
}
