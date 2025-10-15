package world

import (
	"github.com/aquilax/go-perlin"
	"github.com/danharasymiw/trains/types"
)

const (
	alpha = 1.0
	beta  = 2.0
	n     = 4
)

func Generate(w *World, seed int64) {
	// Base elevation map
	pElev := perlin.NewPerlin(alpha, beta, n, seed)
	scaleElev := 0.01

	// Tree/noise map
	pTree := perlin.NewPerlin(alpha, beta, n, seed+1234)
	scaleTree := 0.05 // higher frequency for smaller forest patches

	for y := 0; y < len(w.Tiles); y++ {
		for x := 0; x < len(w.Tiles[y]); x++ {
			// --- Elevation ---
			v := pElev.Noise2D(float64(x)*scaleElev, float64(y)*scaleElev)
			v = (v + 1) * 0.5

			var tileType types.TileType

			switch {
			case v < 0:
				tileType = types.TileWater
			case v < 0.75:
				tileType = types.TileGrass
			default:
				tileType = types.TileMountain
			}

			// --- Trees overlay ---
			if tileType == types.TileGrass {
				// Generate tree density noise
				t := pTree.Noise2D(float64(x)*scaleTree, float64(y)*scaleTree)
				t = (t + 1) * 0.5 // normalize 0–1

				if t > 0.8 {
					tileType = types.TileTree
				}
			}

			w.Tiles[y][x] = &types.Tile{Type: tileType}
		}
	}
}
