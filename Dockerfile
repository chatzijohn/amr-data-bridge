# -------------------------
# Stage 1: Build the Go binary
# -------------------------
FROM golang:1.22-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git

# Set the working directory inside the container
WORKDIR /app

# Copy go mod files first (for dependency caching)
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the entire source code
COPY . .

# Build the Go binary statically
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o amr-data-bridge ./cmd || go build -o amr-data-bridge .

# -------------------------
# Stage 2: Create minimal runtime image
# -------------------------
FROM gcr.io/distroless/base-debian12

# Set the working directory
WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /app/amr-data-bridge .

# Expose application port (adjust to your API port)
EXPOSE 8080

# Run the compiled binary
ENTRYPOINT ["/app/amr-data-bridge"]
