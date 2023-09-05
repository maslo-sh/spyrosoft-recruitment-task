package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"spyrosoft_recruitment/types"
	"sync"
	"time"
)

const (
	API_URL        = "http://api.nbp.pl/api/exchangerates/rates/a/eur/last/100/"
	FETCH_INTERVAL = 5
	FETCHES_AMOUNT = 10
)

func main() {
	initLogger()
	var mu sync.Mutex

	req, err := prepareHttpRequest()
	if err != nil {
		log.Fatalf("Failed to prepare GET request: %e", err)
	}

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

	content, err := decompressGzippedResponse(resp)
	if err != nil {
		log.Fatalf("Failed to read compressed body content: %e", err)
	}

	isJsonValid := json.Valid(content)

	var summary types.ExchangeRatesSummary

	err = json.Unmarshal(content, &summary)
	if err != nil {
		log.Fatalf("Failed to unmarshall request content: %e", err)
	}

	mu.Lock()
	log.Printf("REQUEST TIME: %d ms; STATUS CODE: %d; CONTENT TYPE: %s; VALID JSON SYNTAX: %t", elapsed.Milliseconds(), statusCode, contentType, isJsonValid)
	mu.Unlock()
}

func prepareHttpRequest() (*http.Request, error) {
	req, err := http.NewRequest("GET", API_URL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare HTTP GET request: %e", err)
	}

	addHeaders(req)

	return req, nil
}

func addHeaders(req *http.Request) {
	req.Header.Set("Host", "api.nbp.pl")
	req.Header.Set("User-Agent", "Golang Program")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Accept-Language", "pl-PL,pl;q=0.9,en-US;q=0.8,en;q=0.7")

	//gzip encoding results in a much smaller response body
	req.Header.Set("Accept-Encoding", "deflate, gzip")
}

func decompressGzippedResponse(response *http.Response) ([]byte, error) {
	gzipBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read compressed body content: %e", err)
	}

	bytesReader := bytes.NewReader(gzipBytes)
	gzipReader, err := gzip.NewReader(bytesReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create gzip reader: %e", err)
	}

	content, err := ioutil.ReadAll(gzipReader)
	if err != nil {
		return nil, fmt.Errorf("failed to read compressed body content: %e", err)
	}

	return content, nil
}

func initLogger() {
	log.SetFlags(0)
	file, err := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("Failed to create log file: %e", err)
	}

	log.SetPrefix(time.Now().Format("[01-02-2006 15:04:05] "))

	multi := io.MultiWriter(file, os.Stdout)
	log.SetOutput(multi)
}
