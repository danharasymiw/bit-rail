package types

//go:generate go run ../cmd/generator/main.go -type=Track -output=../message/track_gen.go
type Track struct {
	Direction Dir
	HasSignal bool
	SignalDir Dir
	Block     *Block
}
