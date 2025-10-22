package engine

import (
	"encoding/json"
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
	responseCh *chan outgoingMessage
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
}

type playerConnection struct {
	playerID   string
	ws         *websocket.Conn
	outgoingCh chan outgoingMessage
}

type networkManager struct {
	players     map[string]*playerConnection
	playersMu   sync.RWMutex
	upgrader    websocket.Upgrader
	incomingCh  chan playerMessage   // Shared channel for ALL players
	broadcastCh chan outgoingMessage // Shared channel for ALL players
}

func newNetworkManager() *networkManager {
	return &networkManager{
		players:     make(map[string]*playerConnection),
		incomingCh:  make(chan playerMessage, 100),
		broadcastCh: make(chan outgoingMessage, 100),
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin:     func(r *http.Request) bool { return true },
		},
	}
}

func (nm *networkManager) startServer(readyCh chan<- struct{}) {
	http.HandleFunc("/ws", nm.wsHandler)

	listener, err := net.Listen("tcp", ":2977")
	if err != nil {
		panic(err)
	}
	logrus.Info("Server ready on :2977")
	close(readyCh)

	go nm.broadcastLoop()
	http.Serve(listener, nil)
}

func (nm *networkManager) wsHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := nm.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	var msg message.Message
	if err := ws.ReadJSON(&msg); err != nil {
		ws.Close()
		return
	}

	if msg.Type != message.MessageTypeLogin {
		logrus.Warn("First message was not login")
		ws.Close()
		return
	}

	var loginMsg message.LoginMessage
	if err := json.Unmarshal(msg.Data, &loginMsg); err != nil {
		logrus.Errorf("Failed to unmarshal login message: %v", err)
		ws.Close()
		return
	}

	// Create response channel for this player
	responseCh := make(chan outgoingMessage, 100)

	playerConn := &playerConnection{
		playerID:   loginMsg.Username,
		ws:         ws,
		outgoingCh: responseCh,
	}

	nm.playersMu.Lock()
	nm.players[loginMsg.Username] = playerConn
	nm.playersMu.Unlock()

	// Send login message to engine for processing
	nm.incomingCh <- playerMessage{
		playerID:   loginMsg.Username,
		message:    &incomingMessage{loginMessage: &loginMsg},
		responseCh: &responseCh,
	}

	go nm.handleRead(playerConn)
	nm.handleWrite(playerConn)
}

func (nm *networkManager) handleRead(playerConn *playerConnection) {
	logEntry := logrus.WithField("player", playerConn.playerID)
	for {
		var msg message.Message
		if err := playerConn.ws.ReadJSON(&msg); err != nil {
			logEntry.Debugf("WebSocket read error: %v", err)
			return
		}

		var incoming incomingMessage

		switch msg.Type {
		case message.MessageTypeChat:
			var chatMsg message.ChatMessage
			if err := json.Unmarshal(msg.Data, &chatMsg); err != nil {
				logEntry.Errorf("Error unmarshaling chat message: %v", err)
				continue
			}
			incoming.chatMessage = &chatMsg

		case message.MessageTypeGetChunks:
			var getChunksMsg message.GetChunksMessage
			if err := json.Unmarshal(msg.Data, &getChunksMsg); err != nil {
				logEntry.Errorf("Error unmarshaling get chunks message: %v", err)
				continue
			}
			incoming.getChunksMessage = &getChunksMsg

		default:
			logEntry.Debugf("Unknown message type: %d", msg.Type)
			continue
		}

		nm.incomingCh <- playerMessage{
			playerID:   playerConn.playerID,
			message:    &incoming,
			responseCh: &playerConn.outgoingCh,
		}
	}
}

func (nm *networkManager) handleWrite(playerConn *playerConnection) {
	logEntry := logrus.WithField("player", playerConn.playerID)
	defer nm.disconnectPlayer(playerConn.playerID)
	defer playerConn.ws.Close()

	for outgoing := range playerConn.outgoingCh {
		var msgType message.MessageType
		var data []byte
		var err error

		if outgoing.initialLoadMessage != nil {
			msgType = message.MessageTypeInitialLoad
			data, err = json.Marshal(outgoing.initialLoadMessage)
		} else if outgoing.chatMessage != nil {
			msgType = message.MessageTypeChat
			data, err = json.Marshal(outgoing.chatMessage)
		} else if outgoing.chunksMessage != nil {
			msgType = message.MessageTypeChunks
			data, err = json.Marshal(outgoing.chunksMessage)
		} else {
			logEntry.Warn("Unknown outgoing message type")
			continue
		}

		if err != nil {
			logEntry.Errorf("Error marshaling message: %v", err)
			continue
		}

		msg := message.Message{
			Type: msgType,
			Data: data,
		}

		if err := playerConn.ws.WriteJSON(msg); err != nil {
			logEntry.Errorf("WebSocket write error: %v", err)
			return
		}
	}
}

func (nm *networkManager) broadcastLoop() {
	for msg := range nm.broadcastCh {
		nm.playersMu.RLock()
		for _, player := range nm.players {
			player.outgoingCh <- msg
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
	}
	nm.playersMu.Unlock()
}
