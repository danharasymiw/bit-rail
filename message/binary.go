package message

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/danharasymiw/bit-rail/trains"
	"github.com/danharasymiw/bit-rail/types"
	"github.com/google/uuid"
)

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

func writeTileType(w io.Writer, v *types.TileType) error {
	if err := binaryWrite(w, uint8(*v)); err != nil {
		return fmt.Errorf("error writing TileType: %v", err)
	}
	return nil
}

func readTileType(r io.Reader) (*types.TileType, error) {
	var b uint8
	if err := binaryRead(r, &b); err != nil {
		return nil, fmt.Errorf("error reading TileType: %v", err)
	}
	v := types.TileType(b)
	return &v, nil
}

func writeDir(w io.Writer, v *types.Dir) error {
	// Dir is uint8, so just use binary.Write - simple and efficient for single byte
	if err := binaryWrite(w,  uint8(*v)); err != nil {
		return fmt.Errorf("error writing Dir: %v", err)
	}
	return nil
}

func readDir(r io.Reader) (*types.Dir, error) {
	var b uint8
	if err := binaryRead(r, &b); err != nil {
		return nil, fmt.Errorf("error reading Dir: %v", err)
	}
	v := types.Dir(b)
	return &v, nil
}

func writeCarType(w io.Writer, v *trains.CarType) error {
	if err := binaryWrite(w,  uint8(*v)); err != nil {
		return fmt.Errorf("error writing CarType: %v", err)
	}
	return nil
}

func readCarType(r io.Reader) (*trains.CarType, error) {
	var b uint8
	if err := binaryRead(r,  &b); err != nil {
		return nil, fmt.Errorf("error reading CarType: %v", err)
	}
	v := trains.CarType(b)
	return &v, nil
}