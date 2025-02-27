# Use the official Golang image as the base image
FROM golang:1.23.4 AS builder

RUN apt-get update && apt-get install -y ffmpeg

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod tidy

# Copy the rest of the application code
COPY . .

RUN go build -o server .

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/server .

# Expose the port that the Go application listens on
EXPOSE 8080

# Default command to run your application
CMD ["./server"]