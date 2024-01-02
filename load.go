package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

// LogEntry represents the structure of the log entry payload
type LogEntry struct {
	Level      string    `json:"level"`
	Message    string    `json:"message"`
	ResourceID string    `json:"resourceId"`
	Timestamp  time.Time `json:"timestamp"`
	TraceID    string    `json:"traceId"`
	SpanID     string    `json:"spanId"`
	Commit     string    `json:"commit"`
	Metadata   MetaData  `json:"metadata"`
}

type MetaData struct {
	ParentResourceID string `json:"parentResourceId"`
}

// generateRandomString generates a random string of given length
func generateRandomString(length int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, length)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// generateRandomTimestamp generates a random timestamp
func generateRandomTimestamp() time.Time {
	return time.Now().UTC().Add(time.Duration(-rand.Intn(365*24*60*60)) * time.Second)
}

// generateRandomPayload generates a random log entry payload
func generateRandomPayload() LogEntry {
	return LogEntry{
		Level:      []string{"error", "info", "debug", "fatal"}[rand.Intn(4)],
		Message:    fmt.Sprintf("Log message: %s", generateRandomString(10)),
		ResourceID: fmt.Sprintf("server-%d", rand.Intn(9000)+1000),
		Timestamp:  generateRandomTimestamp(),
		TraceID:    fmt.Sprintf("%s-%s-%d", generateRandomString(3), generateRandomString(3), rand.Intn(900)+100),
		SpanID:     fmt.Sprintf("span-%d", rand.Intn(900)+100),
		Commit:     generateRandomString(7),
		Metadata:   MetaData{ParentResourceID: fmt.Sprintf("server-%d", rand.Intn(9000)+1000)},
	}
}

// sendRequest sends a POST request to the given URL with the provided payload
func sendRequest(url string, payload LogEntry, successCh, errorCh chan bool) {
	defer func() { successCh <- true }()

	jsonPayload, _ := json.Marshal(payload)
	resp, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	resp.Header.Add("Accept", "*/*")
	resp.Header.Add("Content-Type", "application/json")
	res, err := http.DefaultClient.Do(resp)

	if err != nil {
		fmt.Printf("Error: %s\n", err)
		errorCh <- true
		return
	}
	defer res.Body.Close()
	fmt.Printf("Response: %s\n", res.Status)
}

// runLoadTest runs the load test for the specified duration
func runLoadTest(duration time.Duration, url string, maxConcurrency int) (int, int, int) {
	endTime := time.Now().Add(duration)

	// Use channels to control concurrency and collect results
	successCh := make(chan bool)
	errorCh := make(chan bool)

	var wg sync.WaitGroup
	var successCount, errorCount int

	for time.Now().Before(endTime) {
		payload := generateRandomPayload()

		// Acquire a semaphore to control concurrency
		wg.Add(1)

		// Send the request in a goroutine
		go func(p LogEntry) {
			defer wg.Done()

			select {
			case <-successCh:
				successCount++
			case <-errorCh:
				errorCount++
			}

		}(payload)

		// Send the request in a goroutine
		go sendRequest(url, payload, successCh, errorCh)

	}

	// Wait for all goroutines to finish
	wg.Wait()

	close(successCh)
	close(errorCh)

	return successCount, errorCount, successCount + errorCount
}

func main() {
	// URL of the endpoint to test
	endpointURL := "http://localhost:1323/public/ingest"

	// Duration of the load test in seconds
	testDuration := 30 * time.Second

	// Maximum concurrency (number of parallel requests)
	maxConcurrency := 10

	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())

	// Run the load test
	successCount, errorCount, totalCount := runLoadTest(testDuration, endpointURL, maxConcurrency)

	fmt.Printf("Total requests: %d\n", totalCount)
	fmt.Printf("Successful requests: %d\n", successCount)
	fmt.Printf("Failed requests: %d\n", errorCount)
}
