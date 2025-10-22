package main

import (
	"flag"
	"log"
	"time"

	"github.com/danharasymiw/bit-rail/client"
	"github.com/danharasymiw/bit-rail/engine"
	"github.com/danharasymiw/bit-rail/world/test_worlds"
)

func main() {
	serverMode := flag.Bool("server", false, "Run as headless server")
	localMode := flag.Bool("local", false, "Run server and client together")
	flag.Parse()

	if *serverMode {
		// Run as headless server
		w := test_worlds.NewPerlinWorld(123, 123)
		eng := engine.New(w, 150*time.Millisecond)
		eng.RunHeadless()
	} else if *localMode {
		// Run server and client together in this process
		w := test_worlds.NewPerlinWorld(123, 123)
		eng := engine.New(w, 150*time.Millisecond)
		eng.RunLocal()
	} else {
		// Default: Run as client only
		c, _ := client.New()
		if err := c.Run(); err != nil {
			log.Fatal(err)
		}
	}
}
