# Use a Debian base image that has glibc
FROM golang:1.21-bullseye AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main ./cmd/grpcserver

# Use a minimal base image for the final container
FROM debian:bullseye-slim
WORKDIR /app
COPY --from=builder /app/main .

CMD ["./main"]
