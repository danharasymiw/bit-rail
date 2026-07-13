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

func setupLogging(logFile string, debugMode bool, toStdout bool) *os.File {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	if debugMode {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}

	f, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o666)
	if err != nil {
		logrus.Fatalf("Failed to open log file %s: %v", logFile, err)
	}

	if toStdout {
		logrus.SetOutput(io.MultiWriter(os.Stdout, f))
	} else {
		logrus.SetOutput(f)
	}

	return f
}

func main() {
	serverMode := flag.Bool("server", false, "Run as headless server")
	localMode := flag.Bool("local", false, "Run server and client together")
	debugMode := flag.Bool("debug", false, "Enable debug logging")
	addr := flag.String("addr", ":2977", "Server listen address")
	connect := flag.String("connect", "ws://localhost:2977/ws", "Server WebSocket URL (client mode)")
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
		logFile = setupLogging("server.log", *debugMode, true)
		logrus.Info("Starting server in headless mode")
		eng := engine.New(w, tickDur, *addr)
		eng.Run(make(chan struct{}), make(chan struct{}))
	} else if *localMode {
		logFile = setupLogging("server.log", *debugMode, false)
		logrus.Info("Starting server and client in local mode")
		eng := engine.New(w, tickDur, *addr)

		c, quitCh := client.New(*debugMode, *connect)
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
		logFile = setupLogging("client.log", *debugMode, false)
		logrus.Info("Starting client")
		c, _ := client.New(*debugMode, *connect)
		if err := c.Run(); err != nil {
			log.Fatal(err)
		}
	}
}
