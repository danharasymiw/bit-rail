package engine

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/danharasymiw/bit-rail/message"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

type playerConnection struct {
	playerID string
	ws       *websocket.Conn
	sendChan chan message.Message
}

type initialDataProvider func() message.InitialLoadMessage

type networkManager struct {
	getInitialData initialDataProvider
	players        map[string]*playerConnection
	playersMu      sync.RWMutex
	upgrader       websocket.Upgrader
}

func newNetworkManager(getInitialData initialDataProvider) *networkManager {
	return &networkManager{
		getInitialData: getInitialData,
		players:        make(map[string]*playerConnection),
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin:     func(r *http.Request) bool { return true },
		},
	}
}

func (nm *networkManager) startServer() {
	http.HandleFunc("/ws", nm.wsHandler)
	logrus.Info("Server running on :2977")
	err := http.ListenAndServe(":2977", nil)
	if err != nil {
		panic(err)
	}
}

func (nm *networkManager) wsHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := nm.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	// Read login message
	var loginMsg message.LoginMessage
	if err := ws.ReadJSON(&loginMsg); err != nil {
		ws.Close()
		return
	}

	nm.playersMu.Lock()
	playerConn := &playerConnection{
		playerID: loginMsg.Username,
		ws:       ws,
		sendChan: make(chan message.Message, 10),
	}
	nm.players[loginMsg.Username] = playerConn
	nm.playersMu.Unlock()

	initialLoad := nm.getInitialData()

	if err := ws.WriteJSON(initialLoad); err != nil {
		nm.disconnectPlayer(loginMsg.Username)
		ws.Close()
		return
	}

	// Start message handlers - handleWrite will handle cleanup
	go nm.handleRead(playerConn)
	nm.handleWrite(playerConn) // Block here until connection closes
}

func (nm *networkManager) handleRead(playerConnection *playerConnection) {
	for {
		_, data, err := playerConnection.ws.ReadMessage()
		if err != nil {
			return
		}
		var msg message.Message
		err = json.Unmarshal(data, &msg)
		if err != nil {
			return
		}
		switch msg.Type {
		case message.MessageTypeChat:
			var chatMessage message.ChatMessage
			err = json.Unmarshal(msg.Data, &chatMessage)
			if err != nil {
				return
			}
			playerConnection.sendChan <- msg
		}
	}
}

func (nm *networkManager) handleWrite(playerConnection *playerConnection) {
	defer nm.disconnectPlayer(playerConnection.playerID)
	defer playerConnection.ws.Close()

	for msg := range playerConnection.sendChan {
		if err := playerConnection.ws.WriteJSON(msg); err != nil {
			return
		}
	}
}

func (nm *networkManager) disconnectPlayer(playerID string) {
	nm.playersMu.Lock()
	if player, exists := nm.players[playerID]; exists {
		close(player.sendChan)
		delete(nm.players, playerID)
	}
	nm.playersMu.Unlock()
}
