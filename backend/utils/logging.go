package utils

import (
	"log"
	"os"
	"sync"
)

var (
	serverLog     *log.Logger
	serverLogErr  error
	serverLogOnce sync.Once
)

// Writes to server.log using a dedicated logger, so the standard
// log package output is left untouched. log.Logger is safe for concurrent use.
func UpdateServerLogs(message interface{}) error {
	serverLogOnce.Do(func() {
		// Open the log file in append mode, create it if it doesn't exist.
		// Kept open for the lifetime of the process
		logFile, err := os.OpenFile("server.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			serverLogErr = err
			return
		}
		serverLog = log.New(logFile, "", log.LstdFlags)
	})
	if serverLogErr != nil {
		return serverLogErr
	}

	// Append message to log file
	serverLog.Println(message)

	return nil
}
