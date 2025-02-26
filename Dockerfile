# Use a lightweight base image with Go installed
FROM golang:1.21-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the Go application
RUN go build -o main .

# Create a minimal runtime image
FROM alpine:latest

# Set the working directory in the runtime image
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/main .

# Expose the port your application listens on (if applicable)
EXPOSE 3000

# Run the application
CMD ["./main"]