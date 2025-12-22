package main

import (
	"flag"
	"io"
	"log"
	"os"
	"time"

	"github.com/danharasymiw/bit-rail/client"
	"github.com/danharasymiw/bit-rail/engine"
	"github.com/danharasymiw/bit-rail/world/test_worlds"
	"github.com/sirupsen/logrus"
)

func setupLogging(logFile string, debugMode bool) *os.File {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	if debugMode {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}

	// Open log file
	f, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o666)
	if err != nil {
		logrus.Fatalf("Failed to open log file %s: %v", logFile, err)
	}

	// Write to both file and stdout
	multiWriter := io.MultiWriter(os.Stdout, f)
	logrus.SetOutput(multiWriter)

	return f
}

func main() {
	serverMode := flag.Bool("server", false, "Run as headless server")
	localMode := flag.Bool("local", false, "Run server and client together")
	debugMode := flag.Bool("debug", false, "Enable debug logging")
	flag.Parse()

	var logFile *os.File
	defer func() {
		if logFile != nil {
			logFile.Close()
		}
	}()

	// w := test_worlds.NewPerlinWorld(123, 123)
	w := test_worlds.IntersectingLoopsTestWorld()
	tickDur := 500 * time.Millisecond
	if *serverMode {
		logFile = setupLogging("server.log", *debugMode)
		logrus.Info("Starting server in headless mode")
		eng := engine.New(w, tickDur)
		eng.Run(make(chan struct{}), make(chan struct{}))
	} else if *localMode {
		logFile = setupLogging("server.log", *debugMode)
		logrus.Info("Starting server and client in local mode")
		eng := engine.New(w, tickDur)

		c, quitCh := client.New(*debugMode)
		readyCh := make(chan struct{})

		go eng.Run(quitCh, readyCh)

		// Wait for server to be ready
		<-readyCh
		logrus.Info("Server ready, starting client...")

		if err := c.Run(); err != nil {
			logrus.Errorf("Client error: %v", err)
		}
	} else {
		// Default: Run as client only
		logFile = setupLogging("client.log", *debugMode)
		logrus.Info("Starting client")
		c, _ := client.New(*debugMode)
		if err := c.Run(); err != nil {
			log.Fatal(err)
		}
	}
}
