package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

type ExchangeRate struct {
	No            string  `json:"no"`
	EffectiveDate string  `json:"effectiveDate"`
	Mid           float64 `json:"mid"`
}

type ExchangeRatesSummary struct {
	Table    string         `json:"table"`
	Currency string         `json:"currency"`
	Code     string         `json:"code"`
	Rates    []ExchangeRate `json:"rates"`
}

const (
	API_URL        = "http://api.nbp.pl/api/exchangerates/rates/a/eur/last/100/"
	FETCH_INTERVAL = 5
	FETCHES_AMOUNT = 10
)

func main() {
	initLogger()
	var mu sync.Mutex

	req, err := http.NewRequest("GET", API_URL, nil)

	if err != nil {
		log.Fatalf("Failed to prepare HTTP GET request: %e", err)
	}

	addHeaders(req)

	client := &http.Client{}

	startTime := time.Now()
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to perform GET request: %e", err)
	}

	elapsed := time.Since(startTime)

	defer resp.Body.Close()

	statusCode := resp.StatusCode
	contentType := resp.Header.Get("Content-Type")

	gzipBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read compressed body content: %e", err)
	}

	bytesReader := bytes.NewReader(gzipBytes)
	gzipReader, err := gzip.NewReader(bytesReader)
	if err != nil {
		log.Fatalf("Failed to create gzip reader: %e", err)
	}

	content, err := ioutil.ReadAll(gzipReader)
	if err != nil {
		log.Fatalf("Failed to read compressed body content: %e", err)
	}

	isJsonValid := json.Valid(content)

	var summary ExchangeRatesSummary

	err = json.Unmarshal(content, &summary)
	if err != nil {
		log.Fatalf("Failed to unmarshall request content: %e", err)
	}

	mu.Lock()
	log.Printf("REQUEST TIME: %d ms; STATUS CODE: %d; CONTENT TYPE: %s; VALID JSON SYNTAX: %t", elapsed.Milliseconds(), statusCode, contentType, isJsonValid)
	mu.Unlock()
}

func addHeaders(req *http.Request) {
	req.Header.Set("Host", "api.nbp.pl")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Golang Program")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Sec-GPC", "1")
	req.Header.Set("Accept-Encoding", "deflate, gzip")
	req.Header.Set("Accept-Language", "pl-PL,pl;q=0.9,en-US;q=0.8,en;q=0.7")
}

func initLogger() {
	log.SetFlags(0)
	file, err := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("Failed to create logger: %e", err)
	}

	log.SetPrefix(time.Now().Format("[01-02-2006 15:04:05] "))

	multi := io.MultiWriter(file, os.Stdout)
	log.SetOutput(multi)
}
