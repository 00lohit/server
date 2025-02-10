# Use the official Golang image as the base image with the required version
FROM golang:1.23 as builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go modules manifests
COPY go.mod go.sum ./

# Download Go module dependencies
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go application for ARM64
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o main .

# Use the same base image for the final container
FROM golang:1.23

# Set the working directory inside the container
WORKDIR /app

# Copy the compiled binary from the builder stage
COPY --from=builder /app/main .

# Expose the port your application listens on
EXPOSE 8080

# Command to run the application
CMD ["./main"]
