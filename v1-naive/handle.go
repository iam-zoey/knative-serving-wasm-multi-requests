package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Handle HTTP requests (POST, GET)
func Handle(w http.ResponseWriter, r *http.Request) {
	var input string
	var buf bytes.Buffer

	switch r.Method {
	case http.MethodGet:
		// Read input from query parameter for GET requests
		input = r.URL.Query().Get("input")
	case http.MethodPost:
		// Copy input from request body for POST requests
		_, err := io.Copy(&buf, r.Body)
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusInternalServerError)
			return
		}

		defer r.Body.Close()

		// Convert buffer to string
		input = buf.String()
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if input == "" {
		http.Error(w, "Error: 'input' parameter is required", http.StatusBadRequest)
		return
	}

	// Write the input to the Wasmtime module's stdin
	mu.Lock()
	// fmt.Println("Sending to Wasmtime:", buf.Bytes())
	// if _, err := stdinPipe.Write(buf.Bytes()); err != nil {
	// if _, err := io.Copy(stdinPipe, &buf); err != nil {
	if _, err := stdinPipe.Write([]byte(input + "\n")); err != nil {
		http.Error(w, "Error writing to Wasmtime module", http.StatusInternalServerError)
		mu.Unlock()
		return
	}
	mu.Unlock()

	// Clear previous output from the buffer
	mu.Lock()
	stdoutBuf.Reset()
	mu.Unlock()

	// Wait for output from the Wasmtime module
	time.Sleep(100 * time.Millisecond)
	// mu.Lock()
	// output := stdoutBuf.String()
	// mu.Unlock()

	output := waitForOutput()
	// Respond to the HTTP request with the output from the Wasmtime module
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(output))
}
func waitForOutput() string {
	var output string
	for {
		mu.Lock()
		output = stdoutBuf.String()
		mu.Unlock()
		if output != "" {
			break
		}
	}
	return output
}

// Read stdoutPipe and print the output
func HandleOutput(stdoutPipe io.ReadCloser) {
	// Print the output by scanning the stdoutPipe
	buf := make([]byte, 1024)
	for {
		// Read and save into buf
		n, err := stdoutPipe.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Printf("Error reading stdout: %v\n", err)
			return
		}

		line := string(buf[:n])
		fmt.Println("Received from Wasmtime:", line)
		mu.Lock()
		stdoutBuf.WriteString(line)
		mu.Unlock()
	}

	if err := cmd.Wait(); err != nil {
		fmt.Println("Error waiting for command completion:", err)
		return
	}

}
