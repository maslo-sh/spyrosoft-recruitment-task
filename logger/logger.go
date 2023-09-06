package logger

import (
	"io"
	"log"
	"os"
	"time"
)

func InitLogger() {
	log.SetFlags(0)
	file, err := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("Failed to create log file: %s", err)
	}

	log.SetPrefix(time.Now().Format("[01-02-2006 15:04:05] "))

	multi := io.MultiWriter(file, os.Stdout)
	log.SetOutput(multi)
}
