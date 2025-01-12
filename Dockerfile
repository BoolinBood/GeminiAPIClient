# Use the official Golang image as the base image
FROM golang:1.20

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application code
COPY . .

# Expose the port that the Go application listens on
EXPOSE 8080

# Default command to run your application
CMD ["go", "run", "main.go"]