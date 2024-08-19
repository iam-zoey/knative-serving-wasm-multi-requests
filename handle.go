package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

// Handle HTTP requests (POST, GET)
func Handle(w http.ResponseWriter, r *http.Request) {
	var input string
	var byteCount int64
	var buf bytes.Buffer
	var err error

	switch r.Method {
	// Read input from query parameter for GET requests
	case http.MethodGet:
		input = r.URL.Query().Get("input")
		byteCount = int64(len([]byte(input)))

	// Copy input from request body for POST requests
	case http.MethodPost:
		byteCount, err = io.Copy(&buf, r.Body)
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

	mu.Lock()

	data := fmt.Sprintf("%d %s\n", byteCount, input)
	fmt.Println("Sending to Wasmtime:", data)

	// Write the input to the Wasmtime module's stdin
	if _, err := stdinPipe.Write([]byte(data)); err != nil {
		http.Error(w, "Error writing to Wasmtime module", http.StatusInternalServerError)
		mu.Unlock()
		return
	}
	mu.Unlock()
	go handleOutput(stdoutPipe, byteCount)
	// Clear previous output from the buffer
	mu.Lock()
	stdoutBuf.Reset()
	mu.Unlock()

	output := waitForOutput()

	// fmt.Println("Output from Wasmtime:", output) // optional print statement
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(output))
}

/*
waitForOutput waits for the output from the Wasmtime module
*/
func waitForOutput() string {
	var output string
	for {
		mu.Lock()
		output = stdoutBuf.String()
		mu.Unlock()
		if output != "" {
			break
		}
		fmt.Println(output)
	}
	return output
}

/*
Read the output from the stdoutPipe and save it to the buffer
*/
func handleOutput(stdoutPipe io.ReadCloser, byteCount int64) {
	// Print the output by scanning the stdoutPipe
	buf := make([]byte, byteCount)
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
		// fmt.Println("Received from Wasmtime:", line) // optional print statement
		mu.Lock()
		stdoutBuf.WriteString(line)
		mu.Unlock()
	}

	if err := cmd.Wait(); err != nil {
		fmt.Println("Error waiting for command completion:", err)
		return
	}
}
