package main

import (
    "fmt"
    "net/http"
    "sync"
    "time"
)

const (
    URL = "http://127.0.0.1:10550/proxy/2/process/500" // Change this to the URL you want to test
    ConcurrentRequests = 1000    // Number of concurrent requests
    TestDuration = 5           // Duration of the load test in seconds
)

func main() {
	var wg sync.WaitGroup
	reqChan := make(chan struct{}, ConcurrentRequests)
	var totalRequestTime time.Duration
	var totalRequests, errorCount int
	var minRequestTime, maxRequestTime time.Duration

	// Start the load test
	fmt.Println("Starting load test...")

	startTime := time.Now()

	for time.Since(startTime) < time.Second*time.Duration(TestDuration) {
		// Start ConcurrentRequests number of goroutines
		for i := 0; i < ConcurrentRequests; i++ {
			wg.Add(1)
			reqChan <- struct{}{} // Send signal to start a new request
			go func() {
				defer wg.Done()
				start := time.Now()
				err := sendRequest(URL)
				if err != nil {
					fmt.Printf("Error sending GET request: %v\n", err)
					errorCount++
					return
				}
				elapsed := time.Since(start)
				totalRequestTime += elapsed
				totalRequests++
				if elapsed < minRequestTime || minRequestTime == 0 {
					minRequestTime = elapsed
				}
				if elapsed > maxRequestTime {
					maxRequestTime = elapsed
				}
				<-reqChan // Receive signal indicating request is done
			}()
		}

		// Wait for all goroutines to finish before starting the next batch
		wg.Wait()

		// Sleep for 1 second before starting the next batch of requests
		time.Sleep(time.Second)
	}

	// Calculate metrics
	averageRequestTime := totalRequestTime / time.Duration(totalRequests)

	// Summarize metrics
	fmt.Printf("\nLoad test summary:\n")
	fmt.Printf("Average request time: %v\n", averageRequestTime)
	fmt.Printf("Max request time: %v\n", maxRequestTime)
	fmt.Printf("Min request time: %v\n", minRequestTime)
	fmt.Printf("Total number of requests: %d\n", totalRequests)
	fmt.Printf("Error count: %d\n", errorCount)
}


func sendRequest(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// You can add more logging or metrics here if needed
	// fmt.Printf("Response status code: %d\n", resp.StatusCode)

	return nil
}
