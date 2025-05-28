# Build stage
FROM golang:1.24-alpine AS builder

# Set working directory
WORKDIR /app

# Install git and other dependencies
RUN apk add --no-cache git

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

RUN sed -i 's/go 1.24.2/go 1.21/' go.mod

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o vk_butilka ./cmd/main.go

# Runtime stage
FROM alpine:latest

# Set working directory
WORKDIR /app

# Install ca-certificates for HTTPS connections
RUN apk --no-cache add ca-certificates tzdata

# Create content directory
RUN mkdir -p /app/content

# Copy the binary from builder stage
COPY --from=builder /app/vk_butilka /app/vk_butilka

# Copy .env file if it exists (will be overridden by environment variables if provided)
COPY .env* /app/

# Set user for security
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
RUN chown -R appuser:appgroup /app
USER appuser

# Set environment variables (can be overridden at runtime)
ENV CONTENT_DIR="/app/content"
ENV DONUT_FREQUENCY="5"
ENV POST_INTERVAL_HOURS="3"
ENV DONUT_DURATION="-1"
ENV CONTENT_PER_POST="5"

# Command to run
CMD ["/app/vk_butilka"]
