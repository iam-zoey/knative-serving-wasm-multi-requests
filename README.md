# WebAssembly in Knative Service 
This project explores the integration of WebAssembly (WASM) within Knative Serving. The Go HTTP server handles http requests and forwards data to the WASM module via standard input. The WASM module processes the input (sleep for specific duration or print the data) and returns results through standard output, which are then sent back as the HTTP response.

This model handles multiple requests to avoid the overhead of creating and initializing new instances for each request. However, it requires careful concurrency management to ensure correct processing of multiple requests and may not handle sudden spikes effectively.

## Project Structure 
```
.
├── handle.go        # Code to handle HTTP requests and pass data to the WASM module
├── main.go          # Starts the HTTP server and initializes the WASM module
└── wasm
    ├── main.wasm    # Compiled WebAssembly binary
    └── module.go    # Go code for the WebAssembly module

```
---
## Getting Started 

### Prerequisites
- [Wasmtime installed](https://docs.wasmtime.dev/cli-install.html) 
-  [Knative Cluster](https://knative.dev/docs/getting-started/quickstart-install/) with a registry configured
- Go (v 1.21+)

###  Compile Go code to WASM
To compile the Go code into a WASM binary (`.wasm`):
```shell 
GOOS=wasip1 GOARCH=wasm go build -o wasm/main.wasm wasm/module.go
```

### Build and Push an image  
Build the Docker image and push it to your container registry:
```
docker build . -t <REGISTRY/IMAGE_NAME>
docker push <REGISTRY/IMAGE_NAME>

```
### Apply a Knative service 
#### Edit`service.yaml`
Follow the `CONFIGUREME`tag and provide the name of your Docker image. 

#### Apply the Knative service configuration
After updating the service.yaml with your Docker image name, apply the configuration to your Knative cluster:
```
kubectl apply -f service.yaml
```
Once applied, you can check the service URL by running:
```
kubectl get kservice
```

---
## Testing 
```shell
# TESTING: POST  request
curl -X POST -d "Hi WebAssembly" <SERVICE-URL>
curl -X POST -d "Sleep 2" <SERVICE-URL>


# TESTING: GET request
curl "<SERVICE-URL>/?input=Hello%20from%20GET"
```

---

####  Testing with Docker
Run the container: 
```
 docker run --rm -p 8080:8080 <REGISTRY/IMAGE_NAME>
```
Test the WASM container with GET and POST request:
```shell 
# TESTING: GET request 
curl -X POST -d "Hello from curl" http://localhost:8080/

# TESTING: POST request 
curl "http://localhost:8080/?input=Hello%20from%20GET"

```

#### Testing it locally 
Run the following commands to test in local environment 
```
go run main.go handle.go 

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