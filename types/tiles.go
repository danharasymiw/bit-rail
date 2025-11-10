package types

type TileType uint8

const (
	TileGrass TileType = iota
	TileTrack
	TileIron
	TileWater
	TileTree
	TileMountain
)

type Tile struct {
	Type TileType
}
