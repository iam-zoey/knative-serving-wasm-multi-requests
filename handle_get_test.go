package main

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

/*
TestGetRequest tests the GET request functionality of the Wasmtime module
*/
func TestGetMethod(t *testing.T) {
	// Initialize the Wasmtime module
	if err := initModule(); err != nil {
		t.Fatalf("Failed to initialize Wasmtime module: %v", err)
	}

	// Create test server
	ts := httptest.NewServer(http.HandlerFunc(Handle))
	defer ts.Close()

	numRequests := 20

	for i := 0; i < numRequests; i++ {
		// Create a GET HTTP request with query parameter
		req, err := http.NewRequest(http.MethodGet, ts.URL+"/path?input=HelloWorld", nil)
		if err != nil {
			t.Fatal(err)
		}

		// Send the request
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()

		// Read the response body
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}

		// Check the response
		fmt.Printf("%d test\n", i+1) // optional print statement
		expected := "========================\nWASM received: HelloWorld\n"
		actual := string(body)
		if actual != expected {
			t.Errorf("Unexpected response for request %d: got %q, want %q", i, actual, expected)
		} else {
			fmt.Printf("Response for request %d: %q\n", i, actual)
		}
	}
}
