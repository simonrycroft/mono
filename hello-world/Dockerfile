FROM golang:1.22-alpine

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files first (if present)
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application code
COPY .. .

# Build the Go application
RUN go build -o hello-world

# Expose the port the app runs on
EXPOSE 8080

# Command to run the executable
CMD ["./hello-world"]
