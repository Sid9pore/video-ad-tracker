# Use official Golang image  
FROM golang:1.21 AS builder  

WORKDIR /app  

# Copy go mod and sum files  
COPY go.mod go.sum ./  

# Download all modules  
RUN go mod download  

# Copy the source code to the container  
COPY . .  

# Build the application  
RUN go build -o server ./cmd/main.go  

# Start a new stage from scratch  
FROM alpine:latest  

WORKDIR /root/  

# Copy the Pre-built binary file from the previous stage  
COPY --from=builder /app/server .  

# Expose the server port  
EXPOSE 8080  

# Command to run the executable  
CMD ["./server"]  