# Use a Debian base image that has glibc
FROM golang:1.21-bullseye AS builder

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum from the container directory 
COPY ../../go.mod ../../go.sum ./

# Download dependencies
RUN go mod download

# Copy the entire project into the WORKDIR
COPY ../../ ./

# Build the application
RUN go build -o grpc_client ./cmd/grpcclient

# Use a minimal base image for the final container
FROM debian:bullseye-slim

# Set the working directory
WORKDIR /app

# Copy the compiled binary from the builder stage
COPY --from=builder /app/grpc_client .

# Run the compiled binary
CMD ["./grpc_client"]
