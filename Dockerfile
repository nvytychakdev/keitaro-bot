# Start with the official Golang image as a build stage
FROM golang:1.24.2 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the Go application (replace 'main.go' with your main file if different)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o keitaro-bot

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/keitaro-bot .
COPY --from=builder /app/assets ./assets

# Expose port (optional, match your app's port)
EXPOSE 8080

# Command to run the application
CMD ["./keitaro-bot"]