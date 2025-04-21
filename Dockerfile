# Use official Golang image with 1.22+
FROM golang:1.22 AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all modules
RUN go mod download

# Copy the source code to the container
COPY . .

# Build the application
RUN go build -o server ./cmd/main.go

# Use minimal runtime image
FROM alpine:latest

WORKDIR /root/

# Copy the pre-built binary file from the previous stage
COPY --from=builder /app/server .

# Expose the server port
EXPOSE 8080

# Run the app
CMD ["./server"]
