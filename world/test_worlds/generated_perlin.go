package test_worlds

import "github.com/danharasymiw/trains/world"

func NewPerlinWorld(seed, smoothness int64) *world.World {
	w := world.New(500, 500)
	world.Generate(w, seed)
	return w
}
