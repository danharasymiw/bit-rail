package main

import (
	"time"

	"github.com/danharasymiw/bit-rail/engine"
	"github.com/danharasymiw/bit-rail/world/test_worlds"
)

func main() {
	w := test_worlds.NewPerlinWorld(123, 123)

	eng := engine.New(
		w,
		150*time.Millisecond,
	)
	eng.Run()
	println("bye.")
}
