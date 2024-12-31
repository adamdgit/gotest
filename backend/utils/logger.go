package logging

import (
	"fmt"
	"log"
	"os"
	"time"
)

func UpdateLogFile(message interface{}) error {
	// Open the log file in append mode, create it if it doesn't exist
	logFile, err := os.OpenFile("server.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer logFile.Close()

	// Set the log output to the log file
	log.SetOutput(logFile)

	// Set a custom log prefix (optional)
	log.SetPrefix(fmt.Sprintf("[%s] ", time.Now().Format("2006-01-02 15:04:05")))

	// Append message to log file
	log.Println(message)

	return nil
}
