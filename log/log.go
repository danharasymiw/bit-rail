package log

import (
	"fmt"
	"os"
	"sync"
)

var (
	mu       sync.Mutex
	buffer   []string
	maxLen   = 30
	fileOnce sync.Once
	logFile  *os.File
)

func Log(msg string, args ...any) {
	mu.Lock()
	defer mu.Unlock()
	text := fmt.Sprintf(msg, args...)
	buffer = append(buffer, text)
	if len(buffer) > maxLen {
		buffer = buffer[1:]
	}

	initFile()
	fmt.Fprintln(logFile, text)
	logFile.Sync()
}

func GetLogs() []string {
	mu.Lock()
	defer mu.Unlock()
	cpy := make([]string, len(buffer))
	copy(cpy, buffer)
	return cpy
}

func initFile() {
	fileOnce.Do(func() {
		var err error
		logFile, err = os.Create("train.log")
		if err != nil {
			panic(err)
		}
	})
}
