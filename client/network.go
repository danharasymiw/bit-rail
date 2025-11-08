package client

import (
	"io"

	"github.com/danharasymiw/bit-rail/message"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

type incomingMessage struct {
	chatMessage        *message.ChatMessage
	chunksMessage      *message.ChunksMessage
	initialLoadMessage *message.InitialLoadMessage
	worldUpdateMessage *message.WorldUpdateMessage
}

type outgoingMessage struct {
	loginMessage     *message.LoginMessage
	chatMessage      *message.ChatMessage
	getChunksMessage *message.GetChunksMessage
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

func (nm *clientNetworkManager) start() {
	go nm.readLoop()
	go nm.writeLoop()
}

func (nm *clientNetworkManager) readLoop() {
	defer close(nm.incomingCh)

	for {
		_, reader, err := nm.ws.NextReader()
		if err != nil {
			logrus.Debugf("WebSocket read error: %v", err)
			return
		}

		msgType, err := message.GetMessageType(reader)
		if err != nil {
			logrus.Errorf("Error reading message type: %v", err)
			continue
		}

		var incoming incomingMessage
		switch msgType {
		case message.MessageTypeInitialLoad:
			logrus.Debug("Reading initial load message...")
			incoming.initialLoadMessage, err = message.ReadInitialLoadMessage(reader)
			if err != nil {
				logrus.Errorf("Error reading initial load message: %v", err)
				continue
			}
			logrus.Debugf("Successfully read initial load message (chunks: %d, trains: %d, tracks: %d)",
				len(incoming.initialLoadMessage.Chunks), len(incoming.initialLoadMessage.Trains), len(incoming.initialLoadMessage.Tracks))

		case message.MessageTypeChat:
			incoming.chatMessage, err = message.ReadChatMessage(reader)
			if err != nil {
				logrus.Errorf("Error reading chat message: %v", err)
				continue
			}

		case message.MessageTypeChunks:
			incoming.chunksMessage, err = message.ReadChunksMessage(reader)
			if err != nil {
				logrus.Errorf("Error reading chunks message: %v", err)
				continue
			}

		case message.MessageTypeWorldUpdate:
			incoming.worldUpdateMessage, err = message.ReadWorldUpdateMessage(reader)
			if err != nil {
				logrus.Errorf("Error reading world update message: %v", err)
				continue
			}

		default:
			logrus.Debugf("Unknown incoming message type: %d", msgType)
			continue
		}

		nm.incomingCh <- incoming
	}
}

func (nm *clientNetworkManager) writeLoop() {
	for outgoing := range nm.outgoingCh {
		var writeErr error
		if outgoing.loginMessage != nil {
			writeErr = message.WriteMessage(nm.ws, message.MessageTypeLogin, func(w io.Writer) error {
				return message.WriteLoginMessage(w, outgoing.loginMessage)
			})
		} else if outgoing.chatMessage != nil {
			writeErr = message.WriteMessage(nm.ws, message.MessageTypeChat, func(w io.Writer) error {
				return message.WriteChatMessage(w, outgoing.chatMessage)
			})
		} else if outgoing.getChunksMessage != nil {
			writeErr = message.WriteMessage(nm.ws, message.MessageTypeGetChunks, func(w io.Writer) error {
				return message.WriteGetChunksMessage(w, outgoing.getChunksMessage)
			})
		} else {
			logrus.Warn("Unknown outgoing message type")
			return
		}

		if writeErr != nil {
			logrus.Errorf("Error writing message: %v", writeErr)
		}
	}
}

func (nm *clientNetworkManager) close() {
	close(nm.outgoingCh)
	nm.ws.Close()
}
