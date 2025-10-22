package main

import (
	"flag"
	"log"
	"time"

	"github.com/danharasymiw/bit-rail/client"
	"github.com/danharasymiw/bit-rail/engine"
	"github.com/danharasymiw/bit-rail/world/test_worlds"
	"github.com/sirupsen/logrus"
)

func main() {
	serverMode := flag.Bool("server", false, "Run as headless server")
	localMode := flag.Bool("local", false, "Run server and client together")
	flag.Parse()

	if *serverMode {
		w := test_worlds.NewPerlinWorld(123, 123)
		eng := engine.New(w, 150*time.Millisecond)
		eng.Run(make(chan struct{}), make(chan struct{}))
	} else if *localMode {
		w := test_worlds.NewPerlinWorld(123, 123)
		eng := engine.New(w, 150*time.Millisecond)

		c, quitCh := client.New()
		readyCh := make(chan struct{})

		go eng.Run(quitCh, readyCh)

		// Wait for server to be ready
		<-readyCh
		logrus.Info("Server ready, starting client...")

		if err := c.Run(); err != nil {
			logrus.Printf("Client error: %v", err)
		}
	} else {
		// Default: Run as client only
		c, _ := client.New()
		if err := c.Run(); err != nil {
			log.Fatal(err)
		}
	}
}
