package client

import (
	"os/user"
	"time"

	"github.com/danharasymiw/bit-rail/message"
	"github.com/danharasymiw/bit-rail/world"
	"github.com/gdamore/tcell"
	"github.com/gorilla/websocket"
)

type ChunkCoord struct {
	X, Y int
}
type Client struct {
	w            *world.World
	chunksLoaded map[ChunkCoord]struct{}
	chatMessages []ChatMessage
	username     string

	running bool
	ws      *websocket.Conn

	camX, camY int
	r          Renderer

	quitCh chan struct{}
}

func New() (*Client, chan struct{}) {
	quitCh := make(chan struct{})
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}
	return &Client{
		quitCh:   quitCh,
		running:  false,
		username: usr.Username,
	}, quitCh
}

func (c *Client) Run() error {
	screen, err := tcell.NewScreen()
	if err != nil {
		return err
	}
	if err := screen.Init(); err != nil {
		return err
	}
	defer screen.Fini()

	ws, _, err := websocket.DefaultDialer.Dial("ws://localhost:2977/ws", nil)
	if err != nil {
		return err
	}
	c.ws = ws

	if err := c.ws.WriteJSON(message.LoginMessage{
		Username: c.username,
	}); err != nil {
		return err
	}

	var initialLoadMessage message.InitialLoadMessage
	if err := c.ws.ReadJSON(&initialLoadMessage); err != nil {
		return err
	}
	c.w = world.New(initialLoadMessage.Width, initialLoadMessage.Height)
	c.camX = initialLoadMessage.CameraX
	c.camY = initialLoadMessage.CameraY

	c.chunksLoaded = make(map[ChunkCoord]struct{})
	for _, chunk := range initialLoadMessage.Chunks {
		c.chunksLoaded[ChunkCoord{X: chunk.X, Y: chunk.Y}] = struct{}{}
		for i, tile := range chunk.Tiles {
			worldY := chunk.Y*chunk.Size + i/chunk.Size
			worldX := chunk.X*chunk.Size + i%chunk.Size
			if worldY < c.w.Height && worldX < c.w.Width {
				c.w.Tiles[worldY][worldX] = tile
			}
		}
	}
	for _, train := range initialLoadMessage.Trains {
		c.w.AddTrain(train)
	}

	c.r = NewSimpleRenderer(screen, c.w)

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

	c.running = true

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
					if c.camY < c.w.Height-height {
						c.camY++
					}
				case tcell.KeyDown:
					if c.camY > 0 {
						c.camY--
					}
				case tcell.KeyLeft:
					if c.camX > 0 {
						c.camX -= 2
					}
				case tcell.KeyRight:
					width, _ := c.r.Screen().Size()
					if c.camX < c.w.Width-width {
						c.camX += 2
					}
				}
				if tev.Rune() == 'q' {
					c.running = false
				}
			case *tcell.EventResize:
				screen.Sync()
			}

		case <-ticker.C:
			c.r.Render(c.camX, c.camY, c.chatMessages)
		}
	}

	// Tell whoever launched us that we're done
	close(c.quitCh)
	return nil
}
