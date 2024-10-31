# Stage 1: Build the Go application
FROM golang:1.21.0 AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -o websocket-chat .

# Stage 2: Create a minimal image for the runtime
FROM alpine:3.19

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/websocket-chat .

# Expose the port the app runs on
EXPOSE 8080

# Command to run the executable
CMD ["./websocket-chat"]
