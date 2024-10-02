package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

// Main function reads input from standard input and processes commands.
func main() {
	reader := bufio.NewReader(os.Stdin)
	var byteCountStr []byte

	for {
		byteCountStr = byteCountStr[:0]

		// Get the Byte Count
		for {
			b, err := reader.ReadByte()
			if err != nil {
				if err == io.EOF {
					return // Exit if we reach the end of input
				}
				fmt.Println("Error reading byte:", err)
				return
			}
			if b == ' ' {
				break
			}
			byteCountStr = append(byteCountStr, b)
		}

		// Check if byteCountStr is empty before conversion
		if len(byteCountStr) == 0 {
			fmt.Println("Error: byte count is empty")
			return
		}

		// Convert byte count to integer
		byteCount, err := strconv.Atoi(strings.TrimSpace(string(byteCountStr)))
		if err != nil {
			fmt.Println("Error converting byteCount:", err)
			return
		}

		// Read the specified number of bytes
		data := make([]byte, byteCount)
		n, err := io.ReadFull(reader, data)
		if err != nil {
			if err == io.EOF {
				fmt.Println("Unexpected end of input")
			} else {
				fmt.Println("Error reading data:", err)
			}
			return
		}

		// Check if the number of bytes read is equal to the byte count
		if n != byteCount {
			fmt.Printf("Byte count mismatch: expected %d, got %d\n", byteCount, n)
			return
		}

		fmt.Println("========================\nWASM received:", string(data))

	}
}
