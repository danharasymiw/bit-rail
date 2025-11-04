package message

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

func WriteMessage(ws *websocket.Conn, msgType MessageType, writeBody func(io.Writer) error) error {
	w, err := ws.NextWriter(websocket.BinaryMessage)
	if err != nil {
		return fmt.Errorf("error getting next writer: %v", err)
	}
	defer func() {
		if err := w.Close(); err != nil {
			logrus.Errorf("Error closing writer: %v", err)
		}
	}()

	if _, err := w.Write([]byte{byte(msgType)}); err != nil {
		return fmt.Errorf("error writing message type %d: %v", msgType, err)
	}

	if err := writeBody(w); err != nil {
		return fmt.Errorf("error writing message body: %v", err)
	}

	return nil
}

func GetMessageType(r io.Reader) (MessageType, error) {
	var msgTypeByte [1]byte
	if _, err := io.ReadFull(r, msgTypeByte[:]); err != nil {
		return 0, err
	}
	return MessageType(msgTypeByte[0]), nil
}

// Exists for two reasons:
// 1. To avoid having to type LittleEndian every time we write/read binary
// 2. Makes the binary serialization code easier since we don't have to
// decide whether or not to import binary/encoding
func binaryWrite(w io.Writer, v any) error {
	return binary.Write(w, binary.LittleEndian, v)
}

func binaryRead(r io.Reader, v any) error {
	return binary.Read(r, binary.LittleEndian, v)
}

func writeString(w io.Writer, s string) error {
	length := len(s)
	if length > 255 {
		return fmt.Errorf("string too long: %s", s)
	}
	if err := binaryWrite(w, uint8(length)); err != nil {
		return fmt.Errorf("error writing length for string: %v", err)
	}
	if _, err := w.Write([]byte(s)); err != nil {
		return fmt.Errorf("error writing string contents: %v", err)
	}

	return nil
}

func readString(r io.Reader) (string, error) {
	var length uint8
	if err := binaryRead(r, &length); err != nil {
		return ``, fmt.Errorf("error reading string length: %v", err)
	}

	s := make([]byte, length)
	if _, err := io.ReadFull(r, s); err != nil {
		return ``, fmt.Errorf("error reading string contents: %v", err)
	}
	return string(s), nil
}

func writeUUID(w io.Writer, id *uuid.UUID) error {
	if _, err := w.Write(id[:]); err != nil {
		return fmt.Errorf("error writing UUID: %v", err)
	}
	return nil
}

func readUUID(r io.Reader) (uuid.UUID, error) {
	var id uuid.UUID
	if _, err := io.ReadFull(r, id[:]); err != nil {
		return uuid.Nil, fmt.Errorf("error reading UUID: %v", err)
	}
	return id, nil
}
