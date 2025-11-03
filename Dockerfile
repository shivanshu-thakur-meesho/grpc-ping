# Step 1: Build the Go binary
FROM golang:1.24 AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY proto/ ./proto/
COPY gen/ ./gen/
COPY server/ ./server/

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -o grpc-ping ./server

# Step 2: Run stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/grpc-ping .

EXPOSE 50051

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD echo "Health check not implemented" || exit 1

# Run the binary
CMD ["./grpc-ping"]
