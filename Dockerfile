# Stage 1: Use Ubuntu to install Wasmtime
FROM ubuntu:latest 
# Install necessary tools
RUN apt-get update && \
    apt-get install -y \
    curl \
    gnupg \
    xz-utils \ 
    golang-go \
    && rm -rf /var/lib/apt/lists/*

# Install Wasmtime
RUN curl https://wasmtime.dev/install.sh -sSf | bash

# Copy the Go modules files
COPY . .

# Set the PATH environment variable
ENV PATH="/root/.wasmtime/bin:${PATH}"

# Example command to show Wasmtime installation
CMD ["go", "run", "main.go", "handle.go"]

EXPOSE 8080