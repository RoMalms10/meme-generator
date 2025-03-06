FROM golang:1.19-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o meme-generator .

# Use a minimal alpine image for the final stage
FROM alpine:3.16

WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /app/meme-generator .

# Copy templates directory
COPY templates/ /app/templates/

# Create non-root user
RUN adduser -D -H -u 1000 appuser
USER appuser

# Expose the gRPC port
EXPOSE 50051

CMD ["./meme-generator"]