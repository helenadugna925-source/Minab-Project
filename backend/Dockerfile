# Use correct Go version
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build binary
RUN go build -o backend .

# Final stage
FROM alpine:latest
WORKDIR /app

COPY --from=builder /app/backend .

EXPOSE 8082

CMD ["./backend"]
