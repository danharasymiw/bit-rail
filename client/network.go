package client

import (
	"encoding/json"

	"github.com/danharasymiw/bit-rail/message"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

type incomingMessage struct {
	chatMessage        *message.ChatMessage
	chunksMessage      *message.ChunksMessage
	initialLoadMessage *message.InitialLoadMessage
}

type outgoingMessage struct {
	loginMessage    *message.LoginMessage
	chatMessage     *message.ChatMessage
	getChunkMessage *message.GetChunksMessage
}

type clientNetworkManager struct {
	ws         *websocket.Conn
	incomingCh chan incomingMessage
	outgoingCh chan outgoingMessage
}

func newClientNetworkManager() (*clientNetworkManager, error) {
	ws, _, err := websocket.DefaultDialer.Dial("ws://localhost:2977/ws", nil)
	if err != nil {
		return nil, err
	}
	return &clientNetworkManager{
		ws:         ws,
		incomingCh: make(chan incomingMessage, 100),
		outgoingCh: make(chan outgoingMessage, 100),
	}, nil
}

// Incoming returns the channel for receiving messages from server
func (nm *clientNetworkManager) incoming() <-chan incomingMessage {
	return nm.incomingCh
}

func (nm *clientNetworkManager) outgoing() chan<- outgoingMessage {
	return nm.outgoingCh
}

func (nm *clientNetworkManager) start() {
	go nm.readLoop()
	go nm.writeLoop()
}

func (nm *clientNetworkManager) readLoop() {
	defer close(nm.incomingCh)

	for {
		var msg message.Message
		err := nm.ws.ReadJSON(&msg)
		if err != nil {
			logrus.Infof("WebSocket read error: %v", err)
			return
		}

		var incoming incomingMessage

		switch msg.Type {
		case message.MessageTypeInitialLoad:
			var initialLoadMsg message.InitialLoadMessage
			if err := json.Unmarshal(msg.Data, &initialLoadMsg); err != nil {
				logrus.Errorf("Error unmarshaling initial load message: %v", err)
				continue
			}
			incoming.initialLoadMessage = &initialLoadMsg

		case message.MessageTypeChat:
			var chatMsg message.ChatMessage
			if err := json.Unmarshal(msg.Data, &chatMsg); err != nil {
				logrus.Errorf("Error unmarshaling chat message: %v", err)
				continue
			}
			incoming.chatMessage = &chatMsg

		case message.MessageTypeChunks:
			var chunksMsg message.ChunksMessage
			if err := json.Unmarshal(msg.Data, &chunksMsg); err != nil {
				logrus.Errorf("Error unmarshaling chunks message: %v", err)
				continue
			}
			incoming.chunksMessage = &chunksMsg

		default:
			logrus.Debugf("Unknown message type: %d", msg.Type)
			continue
		}

		nm.incomingCh <- incoming
	}
}

func (nm *clientNetworkManager) writeLoop() {
	for outgoing := range nm.outgoingCh {
		var msgType message.MessageType
		var data []byte
		var err error

		// Determine message type and marshal
		if outgoing.loginMessage != nil {
			msgType = message.MessageTypeLogin
			data, err = json.Marshal(outgoing.loginMessage)
		} else if outgoing.chatMessage != nil {
			msgType = message.MessageTypeChat
			data, err = json.Marshal(outgoing.chatMessage)
		} else if outgoing.getChunkMessage != nil {
			msgType = message.MessageTypeGetChunks
			data, err = json.Marshal(outgoing.getChunkMessage)
		} else {
			logrus.Warn("Unknown outgoing message type")
			continue
		}
		if err != nil {
			logrus.Errorf("Error marshaling message: %v", err)
			continue
		}

		msg := message.Message{
			Type: msgType,
			Data: data,
		}
		if err := nm.ws.WriteJSON(msg); err != nil {
			logrus.Infof("WebSocket write error: %v", err)
			return
		}
	}
}

func (nm *clientNetworkManager) close() {
	close(nm.outgoingCh)
	nm.ws.Close()
}
