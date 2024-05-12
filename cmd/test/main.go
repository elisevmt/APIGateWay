package main

import (
	"net/http"
)

func main() {
	for {
		go func() {
			_, err := http.Get("http://127.0.0.1:10550/proxy/1")
			if err != nil {
			}
		}()
	}
}
