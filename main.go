package main

import (
	"time"

	"github.com/danharasymiw/trains/engine"
	"github.com/danharasymiw/trains/world/test_worlds"
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
