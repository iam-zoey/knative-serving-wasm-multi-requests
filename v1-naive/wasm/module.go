package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"time"
)

// TOOD: depending on the input, run the program for different periods of time
// Make the module.go more complex
// Test with the more concurrent requests
func main() {
	reader := bufio.NewReader(os.Stdin)
	delimiter := byte('\n') // Newline as delimiter

	for {
		request, err := reader.ReadBytes(delimiter)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("Error reading from stdin:", err)
			return
		}

		// Remove the delimiter
		request = bytes.TrimSuffix(request, []byte{delimiter})
		str_request := string(request)

		if str_request == "0" {
			fmt.Println("Exiting...")
			break
		}

		fmt.Println("========================\nWASM received:", string(request))

		// Introduce a delay between requests
		time.Sleep(time.Millisecond * 10)

	}
}
