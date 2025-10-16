package test_worlds

import "github.com/danharasymiw/trains/world"

func NewBlock() *world.World {
	w := world.New(50, 50)

	return w
}
