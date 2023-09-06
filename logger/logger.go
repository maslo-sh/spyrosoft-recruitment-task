package logger

import (
	"io"
	"log"
	"os"
	"strings"
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

func PrintReqInfo(index int, elapsed time.Duration, statusCode int, contentType string, isJsonValid bool, rateOutOfScope []string) {
	log.Printf("<worker-%d> Request Time: %d ms", index, elapsed.Milliseconds())
	log.Printf("<worker-%d> HTTP Status Code: %d", index, statusCode)
	log.Printf("<worker-%d> HTTP Content Type: %s", index, contentType)
	log.Printf("<worker-%d> Is Syntax Valid JSON: %t", index, isJsonValid)
	dates := strings.Join(rateOutOfScope, "; ")
	log.Printf("<worker-%d> Dates with Mid Out Of Scope 4.50 - 4.70: %s", index, dates)
}
