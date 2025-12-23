# Build Stage
FROM golang:alpine AS builder

WORKDIR /app

# Copy dependency files first to leverage Docker cache
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the binary (CGO_ENABLED=0 for static linking)
RUN CGO_ENABLED=0 GOOS=linux go build -o amr-bridge ./cmd/server

# Final Stage - minimal image
FROM alpine:latest

WORKDIR /app

# Copy binary and default preferences
COPY --from=builder /app/amr-bridge .
COPY preferences.yaml .

EXPOSE 8080

CMD ["./amr-bridge"]