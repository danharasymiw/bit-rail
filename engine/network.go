package engine

import (
	"io"
	"net"
	"net/http"
	"sync"

	"github.com/danharasymiw/bit-rail/message"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

// PlayerMessage wraps an incoming message with player context
type playerMessage struct {
	playerID   string
	message    *incomingMessage
	responseCh chan *outgoingMessage
}

type incomingMessage struct {
	loginMessage     *message.LoginMessage
	chatMessage      *message.ChatMessage
	getChunksMessage *message.GetChunksMessage
}

type outgoingMessage struct {
	initialLoadMessage *message.InitialLoadMessage
	chatMessage        *message.ChatMessage
	chunksMessage      *message.ChunksMessage
	worldUpdateMessage *message.WorldUpdateMessage
}

type playerConnection struct {
	playerID   string
	ws         *websocket.Conn
	outgoingCh chan *outgoingMessage
}

type networkManager struct {
	addr        string
	players     map[string]*playerConnection
	playersMu   sync.RWMutex
	upgrader    websocket.Upgrader
	incomingCh  chan *playerMessage   // Shared channel for ALL players
	broadcastCh chan *outgoingMessage // Shared channel for ALL players
}

func newNetworkManager(addr string) *networkManager {
	return &networkManager{
		addr:        addr,
		players:     make(map[string]*playerConnection),
		incomingCh:  make(chan *playerMessage, 100),
		broadcastCh: make(chan *outgoingMessage, 100),
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin:     func(r *http.Request) bool { return true },
		},
	}
}

func (nm *networkManager) startServer(readyCh chan<- struct{}) {
	mux := http.NewServeMux()
	mux.HandleFunc("/ws", nm.wsHandler)

	listener, err := net.Listen("tcp", nm.addr)
	if err != nil {
		panic(err)
	}
	logrus.Infof("Server ready on %s", nm.addr)
	close(readyCh)

	go nm.broadcastLoop()
	server := &http.Server{Handler: mux}
	server.Serve(listener)
}

func (nm *networkManager) wsHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := nm.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	// Read the first message (must be login)
	_, reader, err := ws.NextReader()
	if err != nil {
		logrus.Errorf("WebSocket read error on initial connection: %v", err)
		ws.Close()
		return
	}

	msgType, err := message.GetMessageType(reader)
	if err != nil {
		logrus.Errorf("Error reading message type on initial connection: %v", err)
		ws.Close()
		return
	}

	if msgType != message.MessageTypeLogin {
		logrus.Warnf("First message was not login, got type: %d", msgType)
		ws.Close()
		return
	}

	loginMsg, err := message.ReadLoginMessage(reader)
	if err != nil {
		logrus.Errorf("Failed to read login message: %v", err)
		ws.Close()
		return
	}

	responseCh := make(chan *outgoingMessage, 100)
	playerConn := &playerConnection{
		playerID:   loginMsg.Username,
		ws:         ws,
		outgoingCh: responseCh,
	}

	nm.playersMu.Lock()
	nm.players[loginMsg.Username] = playerConn
	nm.playersMu.Unlock()

	go nm.handleRead(playerConn)
	// Send login message to engine for processing
	nm.incomingCh <- &playerMessage{
		playerID:   loginMsg.Username,
		message:    &incomingMessage{loginMessage: loginMsg},
		responseCh: responseCh,
	}

	nm.handleWrite(playerConn)
}

func (nm *networkManager) handleRead(playerConn *playerConnection) {
	defer nm.disconnectPlayer(playerConn.playerID)

	logEntry := logrus.WithField("player", playerConn.playerID)

	for {
		_, reader, err := playerConn.ws.NextReader()
		if err != nil {
			logEntry.Errorf("WebSocket read error: %v", err)
			return
		}

		msgType, err := message.GetMessageType(reader)
		if err != nil {
			logEntry.Errorf("Error reading message type: %v", err)
			continue
		}

		var incoming incomingMessage

		switch msgType {
		case message.MessageTypeChat:
			logrus.Debug("received chat message")
			incoming.chatMessage, err = message.ReadChatMessage(reader)
			if err != nil {
				logEntry.Errorf("Error reading chat message: %v", err)
				continue
			}

		case message.MessageTypeGetChunks:
			incoming.getChunksMessage, err = message.ReadGetChunksMessage(reader)
			if err != nil {
				logEntry.Errorf("Error reading get chunks message: %v", err)
				continue
			}

		default:
			logEntry.Errorf("Unknown incoming message type: %d", msgType)
			continue
		}

		nm.incomingCh <- &playerMessage{
			playerID:   playerConn.playerID,
			message:    &incoming,
			responseCh: playerConn.outgoingCh,
		}
	}
}

func (nm *networkManager) handleWrite(playerConn *playerConnection) {
	logEntry := logrus.WithField("playerID", playerConn.playerID)
	for outgoing := range playerConn.outgoingCh {
		var writeErr error
		if outgoing.initialLoadMessage != nil {
			writeErr = message.WriteMessage(playerConn.ws, message.MessageTypeInitialLoad, func(w io.Writer) error {
				err := message.WriteInitialLoadMessage(w, outgoing.initialLoadMessage)
				if err != nil {
					logEntry.Errorf("WriteInitialLoadMessage returned error: %v", err)
				}
				return err
			})
			if writeErr != nil {
				logEntry.Errorf("Error writing initial load message: %v", writeErr)
			}
		} else if outgoing.chatMessage != nil {
			writeErr = message.WriteMessage(playerConn.ws, message.MessageTypeChat, func(w io.Writer) error {
				return message.WriteChatMessage(w, outgoing.chatMessage)
			})
			if writeErr != nil {
				logEntry.Errorf("Error writing chat message: %v", writeErr)
			}
		} else if outgoing.chunksMessage != nil {
			writeErr = message.WriteMessage(playerConn.ws, message.MessageTypeChunks, func(w io.Writer) error {
				return message.WriteChunksMessage(w, outgoing.chunksMessage)
			})
			if writeErr != nil {
				logEntry.Errorf("Error writing chunks message: %v", writeErr)
			}
		} else if outgoing.worldUpdateMessage != nil {
			writeErr = message.WriteMessage(playerConn.ws, message.MessageTypeWorldUpdate, func(w io.Writer) error {
				return message.WriteWorldUpdateMessage(w, outgoing.worldUpdateMessage)
			})
			if writeErr != nil {
				logEntry.Errorf("Error writing world update message: %v", writeErr)
			}
		} else {
			logEntry.Warn("Unknown outgoing message type")
			continue
		}
	}
}

func (nm *networkManager) broadcastLoop() {
	for msg := range nm.broadcastCh {
		nm.playersMu.RLock()
		for _, player := range nm.players {
			select {
			case player.outgoingCh <- msg:
			default:
				logrus.Warnf("Outgoing channel full for player %s, dropping message", player.playerID)
			}
		}
		nm.playersMu.RUnlock()
	}
}

func (nm *networkManager) disconnectPlayer(playerID string) {
	nm.playersMu.Lock()
	if player, exists := nm.players[playerID]; exists {
		close(player.outgoingCh)
		delete(nm.players, playerID)
		logrus.Infof("Player %s disconnected", playerID)
		// Close websocket connection after removing from map to avoid concurrent access
		ws := player.ws
		nm.playersMu.Unlock()
		ws.Close()
		return
	}
	nm.playersMu.Unlock()
}
