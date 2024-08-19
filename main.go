package main

import (
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

var stdoutPipe io.ReadCloser

var stdoutBuf bytes.Buffer
var mu sync.Mutex

func main() {
	// Initialize and start the server
	if err := startServer(); err != nil {
		fmt.Printf("Failed to start server: %v\n", err)
	}
}

// Initialize and start the Wasmtime module with persistent input/output handling
func initModule() error {
	var err error

	// Create command
	cmd = exec.Command("wasmtime", "wasm/main.wasm")

	// Create pipes for stdin and stdout
	stdinPipe, err = cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdin pipe: %v", err)
	}
	stdoutPipe, err = cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdout pipe: %v", err)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start Wasmtime module: %v", err)
	}

	// Handle stdout asynchronously; starting a goroutine
	// go HandleOutput(stdoutPipe) // Calling handleOutput from handle.go

	return nil
}

// Start the HTTP server
func startServer() error {
	// Initialize the Wasmtime module
	if err := initModule(); err != nil {
		return fmt.Errorf("failed to initialize Wasmtime module: %v", err)
	}

	// Start the HTTP server
	http.HandleFunc("/", Handle)
	fmt.Println("Listening on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		return fmt.Errorf("failed to start server: %v", err)
	}
	return nil
}
