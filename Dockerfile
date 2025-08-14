# =========================================
# 1️⃣ Build Stage
# =========================================
FROM golang:1.24-alpine AS builder

# Install tools required for build
RUN apk add --no-cache git make upx

WORKDIR /app

# Copy go.mod & go.sum first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Install migration tools
RUN make setup

# Build the application (strip debug info, trim paths)
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w" -trimpath \
    -o main .

# Optional: Compress binary to reduce size further
RUN upx --best --lzma main

# =========================================
# 2️⃣ Final Production Image
# =========================================
FROM gcr.io/distroless/static-debian12

WORKDIR /

# Copy compiled app from builder
COPY --from=builder /app/main /main

# Copy migrations (only if needed in production)
COPY --from=builder /app/adapters/db/migrations /adapters/db/migrations

# ✅ Copy entrypoint.sh directly from local context to avoid .dockerignore issues
COPY scripts/entrypoint.sh /entrypoint.sh

# Make entrypoint executable
USER 0
RUN chmod +x /entrypoint.sh

# Expose service port
EXPOSE 8080

# Entrypoint & CMD
ENTRYPOINT ["/entrypoint.sh"]
CMD ["/main"]
