# Use the official Golang image as the base image for building
FROM golang:1.23 as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o main .

# Start a new stage with only the compiled binary
FROM golang:1.23

WORKDIR /app
COPY --from=builder /app/main .
EXPOSE 8080

CMD ["./main"]
