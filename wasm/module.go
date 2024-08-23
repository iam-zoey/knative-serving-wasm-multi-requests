package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

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
					break
				}
				fmt.Println("Error reading byte:", err)
				return
			}
			if b == ' ' {
				break
			}
			byteCountStr = append(byteCountStr, b)
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

		// divide the data into action and param
		action := strings.SplitN(string(data), " ", 2)[0]

		// if action is "Sleep", sleep for the specified duration
		if action == "Sleep" {
			param := strings.TrimSpace(strings.SplitN(string(data), " ", 2)[1])
			duration, err := strconv.Atoi(param)
			if err != nil {
				fmt.Println("Invalid sleep duration:", param)
				return
			}
			time.Sleep(time.Duration(duration) * time.Second)
			fmt.Printf("=================WASM MODULE ======================= \n")
			fmt.Printf("WASM says:  Slept for %d seconds\n", duration)

			// Otherwise, Print the data received from the Wasmtime module
		} else {
			fmt.Println("========================\nWASM received:", string(data))
		}

		time.Sleep(1 * time.Second)
	}
}
