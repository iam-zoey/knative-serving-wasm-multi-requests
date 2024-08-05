package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
)

func TestServer(t *testing.T) {
	// Initialize the Wasmtime module
	if err := initModule(); err != nil {
		t.Fatalf("Failed to initialize Wasmtime module: %v", err)
	}

	// Create test server
	ts := httptest.NewServer(http.HandlerFunc(Handle))
	defer ts.Close()

	var wg sync.WaitGroup
	numRequests := 100

	for i := 0; i < numRequests; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			// Prepare the request body
			requestBody := []byte("hello world")
			resp, err := http.Post(ts.URL, "text/plain", bytes.NewBuffer(requestBody))
			if err != nil {
				t.Errorf("Request %d failed: %v", i, err)
				return
			}
			defer resp.Body.Close()
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("Failed to read response body for request %d: %v", i, err)
				return
			}

			expected := "========================\nWASM received: hello world\n"
			actual := string(body)
			if actual != expected {
				t.Errorf("Unexpected response for request %d: got %q, want %q", i, actual, expected)
			} else {
				fmt.Printf("Response for request %d: %q\n", i, actual)
			}
		}(i)
	}
	wg.Wait()
}
