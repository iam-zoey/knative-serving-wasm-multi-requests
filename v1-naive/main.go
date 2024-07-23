package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os/exec"

	"sync"
)

// Command execution variables
var cmd *exec.Cmd

// Use stdinPipe to pass input to the Wasmtime module
var stdinPipe io.WriteCloser

// Use stdoutPipe to read output from the Wasmtime module
var stdoutPipe io.ReadCloser

var stdoutBuf bytes.Buffer
var mu sync.Mutex

func main() {
	// Initialize the Wasmtime module
	err := initModule()
	if err != nil {
		fmt.Printf("Failed to initialize Wasmtime module: %v\n", err)
		return
	}
	// Start the HTTP server
	http.HandleFunc("/", handle)
	fmt.Println("Listening on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Failed to start server: %v\n", err)
	}

}

// Initialize and start the Wasmtime module with persistent input/output handling
func initModule() error {
	var err error

	// Create command
	cmd = exec.Command("wasmtime", "main.wasm")

	// Create pipes for stdin and stdout
	stdinPipe, err = cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdin pipe: %v", err)
	}
	stdoutPipe, err = cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdout pipe: %v", err)
	}

	// Start the command
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start Wasmtime module: %v", err)
	}

	// Handle stdout asynchronously; starting a goroutine
	go handleOutput(stdoutPipe)

	return nil
}

// FIXME Naive approach: use buffer to handle output from the Wasmtime module
func handleOutput(stdoutPipe io.ReadCloser) {
	// Print the output by scanning the stdoutPipe
	scanner := bufio.NewScanner(stdoutPipe)
	for scanner.Scan() {
		line := scanner.Text() // Read the line
		fmt.Println("Received from Wasmtime: ", line)
		mu.Lock()                          // lock to prevent overwrite of the stdoutbuffer
		stdoutBuf.WriteString(line + "\n") // Write the line to the buffer
		mu.Unlock()                        // release the lock

	}
	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading stdout: %v\n", err)
	}

	if err := cmd.Wait(); err != nil {
		fmt.Println("Error waiting for command completion:", err)
		return
	}
}

// Handle HTTP requests (POST, GET)
func handle(w http.ResponseWriter, r *http.Request) {
	var input string

	switch r.Method {
	case http.MethodGet:
		// Read input from query parameter for GET requests
		input = r.URL.Query().Get("input")
	case http.MethodPost:
		// Copy input from request body for POST requests
		var buf bytes.Buffer
		if _, err := io.Copy(&buf, r.Body); err != nil {
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
	if _, err := stdinPipe.Write([]byte(input + "\n")); err != nil {
		http.Error(w, "Error writing to Wasmtime module", http.StatusInternalServerError)
		mu.Unlock()
		return
	}
	mu.Unlock()

	// ============================================
	// // Capture the output from the Wasmtime module
	var output string
	output = stdoutBuf.String()
	// ============================================
	// var output string
	// for i := 0; i < 3; i++ {
	// 	mu.Lock()
	// 	output = stdoutBuf.String()
	// 	if output != "" {
	// 		stdoutBuf.Reset()
	// 		mu.Unlock()
	// 		break
	// 	}
	// 	mu.Unlock()
	// }

	// Respond to the HTTP request with the output from the Wasmtime module
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(output))
}
