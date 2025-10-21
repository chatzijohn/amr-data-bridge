# -------------------------
# Stage 1: Build the Go binary
# -------------------------
FROM golang:1.25.1-alpine AS builder

# Install git and build dependencies
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy dependency files first for caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the Go binary from cmd/server/main.go
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o amr-data-bridge ./cmd/server

# -------------------------
# Stage 2: Minimal runtime image
# -------------------------
FROM gcr.io/distroless/base-debian12

# Create working directory
WORKDIR /app

# Copy the compiled binary from builder stage
COPY --from=builder /app/amr-data-bridge .



# Expose the application port
EXPOSE 7050

# Run the compiled binary
ENTRYPOINT ["/app/amr-data-bridge"]


