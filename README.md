# R
This project shows how to compile Go code into WebAssembly (WASM) and run it while embedded within an HTTP server The WebAssembly module, compiled from module.go, processes and prints the data it receives from HTTP server requests.

## Project Structure 
```
.
├── handle.go        # Code to handle HTTP requests and pass data to the WASM module
├── main.go          # Starts the HTTP server and initializes the WASM module
└── wasm
    ├── main.wasm    # Compiled WebAssembly binary
    └── module.go    # Go code for the WebAssembly module

```
### WebAssembly Module
The WebAssembly module (module.go) is embedded within the HTTP server. It receives data through incoming requests, processes it, and prints the received data

---
## Getting started 

###  Compiling Go code to WebAssembly 
```shell 
GOOS=wasip1 GOARCH=wasm go build -o wasm/main.wasm wasm/module.go
```

### Building and Running the Docker Container
```shell 
docker build -t . <IMAGE_NAME>
docker run -p 8080:8080 -t <IMAGE_NAME>
```
---
## Testing 
With the server running in the Docker container, you can send HTTP requests to http://localhost:8080 for testing. If you would like to run locally, `go run main.go handle.go`

```
# Test GET method 
curl "http://localhost:8080/path?input=HelloWorld"

# Test POST method
curl -X POST "http://localhost:8080" -d "Sleep 2" 
```

Or you can simply run test
```
go test
```
Note: The test file sends multiple reqeusts for basic operation (Sleep, print), it takea a while (about a minute)