# Start with Go base image
FROM golang:1.24.2-alpine

# Create app directory
WORKDIR /build

# Copy go.mod and go.sum first
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of your code
COPY . .

# Build the Go binary
RUN go build -o main .

# Use a minimal image to run the binary
FROM alpine:latest

WORKDIR /root/

# Copy the binary and other folders
COPY --from=0 /build/main .
COPY --from=0 /build/public ./public
COPY --from=0 /build/storage ./storage
COPY --from=0 /build/view ./view

# Expose port 8080 (required by Cloud Run)
EXPOSE 8080

CMD ["./main"]
