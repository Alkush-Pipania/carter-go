# Stage 1: Build
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install git (required for some Go dependencies)
RUN apk add --no-cache git

# Copy go mod files first for better caching
COPY go.mod go.sum ./

# Enable automatic toolchain download
ENV GOTOOLCHAIN=auto

RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/server ./cmd/api

# Stage 2: Runtime
FROM alpine:3.19

WORKDIR /app

# Install CA certificates for HTTPS requests
RUN apk add --no-cache ca-certificates tzdata

# Create non-root user
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# Copy binary from builder
COPY --from=builder /app/server .

# Change ownership
RUN chown -R appuser:appgroup /app

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 8080

# Run the application
CMD ["./server"]
