package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

/*
TestBasicRequest tests the basic POST Requests functionality of the Wasmtime module
by sending multiple requests (20 in this case) to the test server.
*/
func TestBasicRequest(t *testing.T) {
	// Initialize the Wasmtime module
	if err := initModule(); err != nil {
		t.Fatalf("Failed to initialize Wasmtime module: %v", err)
	}

	// Create test server
	ts := httptest.NewServer(http.HandlerFunc(Handle))
	defer ts.Close()

	numRequests := 20

	for i := 0; i < numRequests; i++ {

		// Send a POST request
		requestBody := []byte("hello world")
		resp, err := http.Post(ts.URL, "text/plain", bytes.NewBuffer(requestBody))
		if err != nil {
			t.Errorf("Request %d failed: %v", i, err)
			return
		}
		defer resp.Body.Close()

		// Read the response body
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Errorf("Failed to read response body for request %d: %v", i, err)
			return
		}

		// Check the response
		fmt.Printf("%d test\n", i) // optional print statement
		expected := "========================\nWASM received: hello world\n"
		actual := string(body)
		if actual != expected {
			t.Errorf("Unexpected response for request %d: got %q, want %q", i, actual, expected)
		} else {
			fmt.Printf("Response for request %d: %q\n", i, actual)
		}
	}
}

/*
TestSleepRequests tests the sleep functionality for POST request to the Wasmtime module
by sending multiple requests (5 in this case) to the test server.
*/
func TestSleepRequests(t *testing.T) {

	numRequests := 5

	// Start the Wasmtime module
	if err := initModule(); err != nil {
		t.Fatalf("Failed to initialize Wasmtime module: %v", err)
	}

	// Create test server
	ts := httptest.NewServer(http.HandlerFunc(Handle))
	defer ts.Close()

	// Test sleep functionality
	for i := 0; i < numRequests; i++ {

		// Send Post request for sleeping
		requestBody := []byte(fmt.Sprintf("Sleep %d", i))
		resp, err := http.Post(ts.URL, "text/plain", bytes.NewBuffer(requestBody))
		if err != nil {
			t.Errorf("Request %d failed: %v", i, err)
			return
		}
		defer resp.Body.Close()

		// Read the response body
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Errorf("Failed to read response body for request %d: %v", i, err)
			return
		}

		// Sleep for i seconds
		time.Sleep(time.Duration(i) * time.Second)
		expected := fmt.Sprintf("Slept for %d seconds\n", i)

		// Check the response
		actual := string(body)
		if actual != expected {
			t.Errorf("Unexpected response for request %d: got %q, want %q", i, actual, expected)
		} else {
			fmt.Printf("Response for a test request %d: %q", i, actual)
		}
	}
}
