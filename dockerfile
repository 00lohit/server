# Use an official Go runtime as a parent image
# FROM golang:1.21-alpine AS builder  # Previous builder stage

# Use the official golang image with a specific version and Debian base
FROM golang:1.23-bookworm AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go source code into the container
COPY . .

# Download Go modules
RUN go mod download

# Build the Go application
RUN go build -o main .

# Use a minimal Alpine Linux image for the final image
# FROM alpine:latest  # Previous final stage

# Use the same base image as the builder for consistency and to avoid libc issues
FROM golang:1.23-bookworm

# Set the working directory
WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /app/main .

# Expose the port the application listens on
EXPOSE 8080

# Command to run the executable
CMD ["./main"]
