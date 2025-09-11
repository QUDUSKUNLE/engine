# Build stage
FROM golang:1.24-alpine AS builder

# Install git and make for go modules and build tools
RUN apk add --no-cache git make

WORKDIR /app

# Copy go mod files first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy code
COPY . .

# Install migration tools
RUN make setup

# Build the application (static binary)
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Production stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates tzdata bash postgresql-client

WORKDIR /root/

# Copy the Go binary
COPY --from=builder /app/main .

# Copy migration tools + files
COPY --from=builder /app/bin/migrate ./bin/migrate
COPY --from=builder /app/adapters/db/migrations ./adapters/db/migrations


# Copy scripts directly into final image
COPY scripts/entrypoint.sh /entrypoint.sh
COPY scripts/migrate.sh /migrate.sh


RUN chmod +x /entrypoint.sh /migrate.sh

# Create logs directory
RUN mkdir -p logs

# Expose port (for local clarity — Railway will override with $PORT)
EXPOSE 8080

# Entrypoint script will run migrations then exec the app
ENTRYPOINT ["/entrypoint.sh"]

# Default command starts the app
CMD ["./main"]
