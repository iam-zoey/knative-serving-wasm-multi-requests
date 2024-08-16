# WebAssembly in Knative Service 
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

### Building and pushing the image to a registry 
```shell 
docker build -t . <registry>/<IMAGE_NAME>
docker push <registry>/<IMAGE_NAME>
```

### Deploying a Knative service 
Fix `service.yaml` with your image name (<registry>/<IMAGE_NAME>)
```
export FUNC_ENABLE_HOST_BUILDER=truek
kubectl apply -f service.yaml 
```

---
## Testing 
```
kubectl get kservice 
```
Copy kservice url and send request 
```
 curl -X POST  http://wasm-module.default.127.0.0.1.sslip.io -d "Sleep 2"
 curl "http://wasm-module.default.127.0.0.1.sslip.io/path?input=HelloWorld"

```


#### Testing it locally 
With the server running in your local environment with `go run main.go handle.go` command
,  you can send HTTP requests to http://localhost:8080 for testing. 
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