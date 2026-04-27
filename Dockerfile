# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy dependency files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the Master engine binary
RUN CGO_ENABLED=0 GOOS=linux go build -o faro-master ./cmd/faro/main.go

# Run stage
FROM alpine:latest

WORKDIR /root/
COPY --from=builder /app/faro-master .
# Copy internal assets for the dashboard
COPY --from=builder /app/internal/pkg/server/assets ./internal/pkg/server/assets

# Expose the dashboard port
EXPOSE 8089

CMD ["./faro-master"]
