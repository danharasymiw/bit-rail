package test_worlds

import "github.com/danharasymiw/bit-rail/world"

func NewBlock() *world.World {
	w := world.New(50, 50)

	return w
}
