package client

import (
	"os/user"
	"time"

	"github.com/danharasymiw/bit-rail/message"
	"github.com/danharasymiw/bit-rail/world"
	"github.com/gdamore/tcell"
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
	nm      *clientNetworkManager

	camX, camY, camSpeed int
	r                    Renderer

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
		camSpeed: 2,
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

	c.nm, err = newClientNetworkManager()
	if err != nil {
		return err
	}
	c.nm.start()

	c.nm.outgoing() <- outgoingMessage{
		loginMessage: &message.LoginMessage{
			Username: c.username,
		},
	}

	if err := c.waitForInitialLoad(); err != nil {
		return err
	}

	c.r = NewSimpleRenderer(screen, c.w)

	events := make(chan tcell.Event, 32)

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
					c.moveCamera(0, c.camSpeed)
				case tcell.KeyDown:
					c.moveCamera(0, -c.camSpeed)
				case tcell.KeyLeft:
					c.moveCamera(-c.camSpeed, 0)
				case tcell.KeyRight:
					c.moveCamera(c.camSpeed, 0)
				}
				if tev.Rune() == 'q' {
					c.running = false
				}
			case *tcell.EventResize:
				screen.Sync()
			}

		case incoming := <-c.nm.incoming():
			c.handleIncomingMessage(incoming)

		case <-ticker.C:
			c.r.Render(c.camX, c.camY, c.chatMessages)
		}
	}

	// Tell whoever launched us that we're done
	c.nm.close()
	close(c.quitCh)
	return nil
}

func (c *Client) waitForInitialLoad() error {
	for incoming := range c.nm.incoming() {
		if incoming.initialLoadMessage != nil {
			return c.handleInitialLoad(incoming.initialLoadMessage)
		}
	}
	return nil
}

func (c *Client) handleInitialLoad(msg *message.InitialLoadMessage) error {
	c.w = world.New(msg.Width, msg.Height)
	c.camX = msg.CameraX
	c.camY = msg.CameraY
	c.chunksLoaded = make(map[ChunkCoord]struct{})

	for _, chunk := range msg.Chunks {
		c.chunksLoaded[ChunkCoord{X: chunk.X, Y: chunk.Y}] = struct{}{}
		for i, tile := range chunk.Tiles {
			worldY := chunk.Y*world.ChunkSize + i/world.ChunkSize
			worldX := chunk.X*world.ChunkSize + i%world.ChunkSize
			if worldY < c.w.Height && worldX < c.w.Width {
				c.w.Tiles[worldY][worldX] = tile
			}
		}
	}

	for _, train := range msg.Trains {
		c.w.AddTrain(train)
	}

	return nil
}

func (c *Client) handleIncomingMessage(incoming incomingMessage) {
	switch {
	case incoming.chatMessage != nil:
		c.chatMessages = append(c.chatMessages, ChatMessage{
			Author:  incoming.chatMessage.Author,
			Message: incoming.chatMessage.Message,
		})

		// Keep only last N messages
		const maxChatMessages = 50
		if len(c.chatMessages) > maxChatMessages {
			c.chatMessages = c.chatMessages[len(c.chatMessages)-maxChatMessages:]
		}

	case incoming.chunksMessage != nil:
		for _, chunk := range incoming.chunksMessage.Chunks {
			c.chunksLoaded[ChunkCoord{X: chunk.X, Y: chunk.Y}] = struct{}{}
			for i, tile := range chunk.Tiles {
				worldY := chunk.Y*world.ChunkSize + i/world.ChunkSize
				worldX := chunk.X*world.ChunkSize + i%world.ChunkSize
				if worldY < c.w.Height && worldX < c.w.Width {
					c.w.Tiles[worldY][worldX] = tile
				}
			}
		}
	}
}

func (c *Client) moveCamera(xDelta, yDelta int) {
	width, height := c.r.Screen().Size()
	if yDelta > 0 && c.camY < c.w.Height-height {
		c.camY += yDelta
	} else if yDelta < 0 && c.camY > 0 {
		c.camY += yDelta
	}
	if xDelta > 0 && c.camX < c.w.Width-width {
		c.camX += xDelta
	} else if xDelta < 0 && c.camX > 0 {
		c.camX += xDelta
	}

	// Check if we need a new chunk
	chunkX := c.camX / world.ChunkSize
	chunkY := c.camY / world.ChunkSize

	if xDelta > 0 {
		c.getChunk(chunkX+1, chunkY)
	} else if xDelta < 0 {
		c.getChunk(chunkX-1, chunkY)
	}
	if yDelta > 0 {
		c.getChunk(chunkX, chunkY+1)
	} else if yDelta < 0 {
		c.getChunk(chunkX, chunkY-1)
	}
}

func (c *Client) getChunk(chunkX, chunkY int) {
	if _, ok := c.chunksLoaded[ChunkCoord{X: chunkX, Y: chunkY}]; ok {
		return
	}
	c.chunksLoaded[ChunkCoord{X: chunkX, Y: chunkY}] = struct{}{}

	c.nm.outgoing() <- outgoingMessage{
		getChunkMessage: &message.GetChunkMessage{
			X: chunkX,
			Y: chunkY,
		},
	}
}
