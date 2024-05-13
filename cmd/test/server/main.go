package main

import (
	"fmt"
	"net/http"
	"log"
	"strconv"
	"time"
)

func processHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the integer value from the URL path
	durationStr := r.URL.Path[len("/process/"):]
	duration, err := strconv.Atoi(durationStr)
	if err != nil {
		http.Error(w, "Invalid duration", http.StatusBadRequest)
		return
	}

	// Simulate processing by sleeping for the specified duration
	time.Sleep(time.Duration(duration) * time.Millisecond)

	// Respond to the client
	fmt.Fprintf(w, "Processed request in %d milliseconds\n", duration)
}



func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		handler.ServeHTTP(w, r)
		log.Printf("%s %s %s", r.Method, r.RequestURI, time.Since(start))
	})
}

func main() {
	port := 9000
	for {
		// Attempt to start the server on the specified port
		addr := fmt.Sprintf(":%d", port)
		fmt.Printf("Attempting to start server on port %d...\n", port)
		if err := http.ListenAndServe(addr, logRequest(http.HandlerFunc(processHandler))); err != nil {
			// If the port is already in use, increase the port number and retry
			if err.Error() == "listen tcp :"+strconv.Itoa(port)+": bind: address already in use" {
				fmt.Printf("Port %d is already in use. Trying the next port...\n", port)
				port++
			} else {
				fmt.Printf("Failed to start server: %v\n", err)
				return
			}
		} else {
			break // Server started successfully
		}
	}

	fmt.Printf("Server started successfully on port %d\n", port)
}
