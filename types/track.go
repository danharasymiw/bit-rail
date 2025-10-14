package types

type Track struct {
	Tile      *Tile
	Direction Dir
	HasSignal bool
	SignalDir Dir
	Block     *Block
}
