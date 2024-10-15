package config

import (
	"fmt"
	"io"
	"log"
	"os"
)

var logFile *os.File

func initializeLogFile() {
	file, err := openLogFile(fmt.Sprintf("./%s.log", GetAppName()))
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}

	logFile = file
	log.SetOutput(io.MultiWriter(file, os.Stdout))
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)

	log.Println("Log file initialized!")
}

func openLogFile(path string) (*os.File, error) {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func closeLogFile(logFile *os.File) {
	if logFile == nil {
		return
	}
	log.Println("Closing log file...")
	err := logFile.Close()
	if err != nil {
		log.Fatal(err, "Error closing log file")
	}
	log.Println("Log file closed successfully!")
}
