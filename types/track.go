package types

type Track struct {
	Direction Dir
	HasSignal bool
	SignalDir Dir
	Block     *Block
}
