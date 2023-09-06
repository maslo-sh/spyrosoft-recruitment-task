package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"spyrosoft_recruitment/logger"
	"spyrosoft_recruitment/types"
	"sync"
	"time"
)

const (
	ApiUrl        = "http://api.nbp.pl/api/exchangerates/rates/a/eur/last/100/"
	FetchInterval = 5
	FetchesAmount = 10
)

func main() {
	logger.InitLogger()
	var mu sync.Mutex

	for {
		var wg sync.WaitGroup
		waitCh := make(chan int)
		wg.Add(FetchesAmount)

		start := time.Now()
		for i := 0; i < FetchesAmount; i++ {
			go apiQueryWorker(i, &mu, &wg)
		}

		go func() {
			wg.Wait()
			close(waitCh)
		}()

		select {
		case <-waitCh:
			elapsed := time.Since(start)
			time.Sleep(FetchInterval*time.Second - elapsed)
		case <-time.After(FetchInterval * time.Second):
			log.Println("Timeout, performing next requests group...")
		}
	}

}

func apiQueryWorker(index int, mu *sync.Mutex, wg *sync.WaitGroup) {
	defer wg.Done()
	req, err := prepareHttpRequest()
	if err != nil {
		log.Fatalf("Failed to prepare GET request: %s", err)
	}

	client := &http.Client{}

	startTime := time.Now()
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to perform GET request: %s", err)
	}

	elapsed := time.Since(startTime)

	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.Fatalf("Failed to close response body: %s", err)
		}
	}()

	statusCode := resp.StatusCode
	contentType := resp.Header.Get("Content-Type")

	// read gzip byte stream and decompress it into readable JSON
	content, err := decompressGzippedResponse(resp)
	if err != nil {
		log.Fatalf("Failed to read compressed body content: %s", err)
	}

	isJsonValid := json.Valid(content)

	var summary types.ExchangeRatesSummary

	err = json.Unmarshal(content, &summary)
	if err != nil {
		log.Fatalf("Failed to unmarshall request content: %s", err)
	}

	//locking mutex to avoid mixing logs from different goroutines
	mu.Lock()
	printReqInfo(index, elapsed, statusCode, contentType, isJsonValid)
	mu.Unlock()
}

func prepareHttpRequest() (*http.Request, error) {
	req, err := http.NewRequest("GET", ApiUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare HTTP GET request: %s", err)
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
		return nil, fmt.Errorf("failed to read compressed body content: %s", err)
	}

	bytesReader := bytes.NewReader(gzipBytes)
	gzipReader, err := gzip.NewReader(bytesReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create gzip reader: %s", err)
	}

	content, err := ioutil.ReadAll(gzipReader)
	if err != nil {
		return nil, fmt.Errorf("failed to read compressed body content: %s", err)
	}

	return content, nil
}

func printReqInfo(index int, elapsed time.Duration, statusCode int, contentType string, isJsonValid bool) {
	log.Printf("<worker-%d> Request Time: %d ms", index, elapsed.Milliseconds())
	log.Printf("<worker-%d> HTTP Status Code: %d", index, statusCode)
	log.Printf("<worker-%d> HTTP Content Type: %s", index, contentType)
	log.Printf("<worker-%d> Is Syntax Valid JSON: %t", index, isJsonValid)
}
