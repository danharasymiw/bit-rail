package client

import (
	"os/user"
	"time"

	"github.com/danharasymiw/bit-rail/message"
	"github.com/danharasymiw/bit-rail/world"
	"github.com/gdamore/tcell"
)

type Client struct {
	w            *world.World
	chunksLoaded map[world.ChunkCoord]struct{}
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

	c.nm.outgoingCh <- outgoingMessage{
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

		case incoming := <-c.nm.incomingCh:
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
	for incoming := range c.nm.incomingCh {
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
	c.chunksLoaded = make(map[world.ChunkCoord]struct{})

	for _, chunk := range msg.Chunks {
		c.chunksLoaded[chunk.Coord] = struct{}{}
		for i, tile := range chunk.Tiles {
			worldY := chunk.Coord.Y*world.ChunkSize + i/world.ChunkSize
			worldX := chunk.Coord.X*world.ChunkSize + i%world.ChunkSize
			if worldY < c.w.Height && worldX < c.w.Width {
				c.w.Tiles[worldY][worldX] = tile
			}
		}
	}

	for _, train := range msg.Trains {
		c.w.AddTrain(train)
	}

	// Ensure we have full chunk buffer (in case initial load didn't include all)
	c.loadChunksAroundCamera()

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
			c.chunksLoaded[chunk.Coord] = struct{}{}
			for i, tile := range chunk.Tiles {
				worldY := chunk.Coord.Y*world.ChunkSize + i/world.ChunkSize
				worldX := chunk.Coord.X*world.ChunkSize + i%world.ChunkSize
				if worldY < c.w.Height && worldX < c.w.Width {
					c.w.Tiles[worldY][worldX] = tile
				}
			}
		}
	}
}

func (c *Client) moveCamera(xDelta, yDelta int) {
	width, height := c.r.Screen().Size()
	newCamX := c.camX + xDelta
	newCamY := c.camY + yDelta
	if newCamX < 0 {
		newCamX = 0
	} else if newCamX > c.w.Width-width {
		newCamX = c.w.Width - width
	}
	if newCamY < 0 {
		newCamY = 0
	} else if newCamY > c.w.Height-height {
		newCamY = c.w.Height - height
	}

	c.camX = newCamX
	c.camY = newCamY

	// Ensure we have a buffer of chunks around the camera
	c.loadChunksAroundCamera()
}

// loadChunksAroundCamera ensures a radius of chunks is loaded around the camera
func (c *Client) loadChunksAroundCamera() {
	const chunkRadius = 3

	centerChunk := world.TileToChunkCoords(c.camX, c.camY)

	chunkCoords := make([]world.ChunkCoord, 0, (2*chunkRadius+1)*(2*chunkRadius+1))
	for dx := -chunkRadius; dx <= chunkRadius; dx++ {
		for dy := -chunkRadius; dy <= chunkRadius; dy++ {
			chunkCoords = append(chunkCoords, world.ChunkCoord{
				X: centerChunk.X + dx,
				Y: centerChunk.Y + dy,
			})
		}
	}

	c.getChunks(chunkCoords)
}

func (c *Client) getChunks(chunkCoords []world.ChunkCoord) {
	missingCoords := make([]world.ChunkCoord, 0)
	for _, coord := range chunkCoords {
		if _, ok := c.chunksLoaded[coord]; ok {
			continue
		}
		missingCoords = append(missingCoords, coord)
		// Technically we don't have it yet but it's been requested to avoid requesting it again
		// Might need to make this more intelligent later
		c.chunksLoaded[coord] = struct{}{}
	}
	if len(missingCoords) == 0 {
		return
	}

	c.nm.outgoingCh <- outgoingMessage{
		getChunkMessage: &message.GetChunksMessage{
			Coords: missingCoords,
		},
	}
}
