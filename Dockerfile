# Use official Golang image as base
FROM golang:1.21.2-alpine

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod tidy

# Copy the source code
COPY . .

# Build the Go application
RUN go build -o main main.go


# Expose the application port
EXPOSE 8080

# Run the application
CMD ["./main"]
