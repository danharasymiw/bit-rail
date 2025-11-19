package client

import (
	"os/user"
	"time"

	"github.com/danharasymiw/bit-rail/message"
	"github.com/danharasymiw/bit-rail/types"
	"github.com/danharasymiw/bit-rail/world"
	"github.com/gdamore/tcell"
	"github.com/sirupsen/logrus"
)

type Client struct {
	w            *world.World
	chunksLoaded map[types.Pos]struct{}
	chatMessages []ChatMessage
	username     string
	debugMode    bool

	running bool
	nm      *clientNetworkManager

	camPos   types.Pos
	camSpeed int
	r        Renderer

	quitCh chan struct{}
}

func New(debugMode bool) (*Client, chan struct{}) {
	quitCh := make(chan struct{})
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}
	c := &Client{
		quitCh:       quitCh,
		running:      false,
		username:     usr.Username,
		camSpeed:     2,
		chatMessages: make([]ChatMessage, 0),
		debugMode:    debugMode,
	}

	logrus.AddHook(&ChatLogHook{client: c})

	return c, quitCh
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

	c.nm.outgoingCh <- &outgoingMessage{
		loginMessage: &message.LoginMessage{
			Username: c.username,
		},
	}

	if err := c.waitForInitialLoad(); err != nil {
		logrus.Errorf("Error waiting for initial load: %v", err)
		return err
	}

	c.r = NewSimpleRenderer(screen, c.w, c.debugMode)

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
			c.r.Render(c.camPos, c.chatMessages)
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
	c.camPos = msg.CameraPos
	c.chunksLoaded = make(map[types.Pos]struct{})

	for _, chunk := range msg.Chunks {
		c.chunksLoaded[chunk.Pos] = struct{}{}
		for i, tile := range chunk.Tiles {
			worldY := chunk.Pos.Y*world.ChunkSize + i/world.ChunkSize
			worldX := chunk.Pos.X*world.ChunkSize + i%world.ChunkSize
			if worldY < c.w.Height && worldX < c.w.Width {
				c.w.Tiles[worldY][worldX] = tile
			}
		}
	}

	for _, train := range msg.Trains {
		c.w.AddTrain(train)
	}
	for _, trackUpdate := range msg.Tracks {
		c.w.AddTrack(trackUpdate.Pos, &trackUpdate.Track)
	}

	// Ensure we have full chunk buffer (in case initial load didn't include all)
	c.loadChunksAroundCamera()

	return nil
}

func (c *Client) handleIncomingMessage(incoming *incomingMessage) {
	switch {
	case incoming.chatMessage != nil:
		c.handleChatMessage(incoming.chatMessage)

	case incoming.chunksMessage != nil:
		for _, chunk := range incoming.chunksMessage.Chunks {
			c.chunksLoaded[chunk.Pos] = struct{}{}
			for i, tile := range chunk.Tiles {
				worldY := chunk.Pos.Y*world.ChunkSize + i/world.ChunkSize
				worldX := chunk.Pos.X*world.ChunkSize + i%world.ChunkSize
				if worldY < c.w.Height && worldX < c.w.Width {
					c.w.Tiles[worldY][worldX] = tile
				}
			}
		}
	case incoming.worldUpdateMessage != nil:
		for _, tu := range incoming.worldUpdateMessage.TilesUpdated {
			c.w.Tiles[tu.Pos.Y][tu.Pos.X] = &tu.Tile
		}
		for _, tu := range incoming.worldUpdateMessage.TracksUpdated {
			c.w.Tracks[tu.Pos] = &tu.Track
		}
		c.w.Trains = incoming.worldUpdateMessage.Trains
	}
}

func (c *Client) handleChatMessage(msg *message.ChatMessage) {
	c.chatMessages = append(c.chatMessages, ChatMessage{
		Author:  msg.Author,
		Message: msg.Message,
	})

	// Keep only last N messages
	const maxChatMessages = 10
	if len(c.chatMessages) > maxChatMessages {
		c.chatMessages = c.chatMessages[len(c.chatMessages)-maxChatMessages:]
	}
}

func (c *Client) moveCamera(xDelta, yDelta int) {
	width, height := c.r.Screen().Size()
	newCamX := c.camPos.X + xDelta
	newCamY := c.camPos.Y + yDelta
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

	c.camPos.X = newCamX
	c.camPos.Y = newCamY
	// Ensure we have a buffer of chunks around the camera
	c.loadChunksAroundCamera()
}

// loadChunksAroundCamera ensures a radius of chunks is loaded around the camera
func (c *Client) loadChunksAroundCamera() {
	const chunkRadius = 3

	centerChunk := world.TileToChunkPos(c.camPos)

	minChunkX, minChunkY := 0, 0
	if centerChunk.X > chunkRadius {
		minChunkX = centerChunk.X - chunkRadius
	}
	if centerChunk.Y > chunkRadius {
		minChunkX = centerChunk.Y - chunkRadius
	}
	maxChunkX := c.w.Width / world.ChunkSize
	maxChunkY := c.w.Height / world.ChunkSize

	chunkPositions := make([]types.Pos, 0, (2*chunkRadius+1)*(2*chunkRadius+1))
	for dx := minChunkX; dx <= maxChunkX; dx++ {
		for dy := minChunkY; dy <= maxChunkY; dy++ {
			chunkPositions = append(chunkPositions, types.Pos{X: centerChunk.X + dx, Y: centerChunk.Y + dy})
		}
	}

	c.getChunks(chunkPositions)
}

func (c *Client) getChunks(positions []types.Pos) {
	missingChunkPositions := make([]types.Pos, 0)
	for _, coord := range positions {
		if _, ok := c.chunksLoaded[coord]; ok {
			continue
		}
		missingChunkPositions = append(missingChunkPositions, coord)
		// Technically we don't have it yet but it's been requested to avoid requesting it again
		// Might need to make this more intelligent later
		c.chunksLoaded[coord] = struct{}{}
	}
	if len(missingChunkPositions) == 0 {
		return
	}

	c.nm.outgoingCh <- &outgoingMessage{
		getChunksMessage: &message.GetChunksMessage{
			Positions: missingChunkPositions,
		},
	}
}
