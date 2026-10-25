# Use a newer, supported version of Go
FROM golang:1.20-bookworm AS builder

WORKDIR /app

# Copy the Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the server code
COPY . .

# Build the Go binary
RUN CGO_ENABLED=0 GOOS=linux go build -o evilqr-server main.go

# Use a lightweight Debian image for the final run
FROM debian:bookworm-slim

WORKDIR /app

# Copy the compiled binary from the builder
COPY --from=builder /app/evilqr-server .

# Copy the templates folder (if the server uses it for the fake pages)
COPY --from=builder /app/templates ./templates

# Expose the port the server runs on
EXPOSE 8080

# Run the server
CMD ["./evilqr-server", "-d", "templates"]
