# Build stage
FROM golang:1.24-alpine AS builder

# Install git and make for go modules and build tools
RUN apk add --no-cache git make

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Install migration tools
RUN make setup

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Production stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/main .

# Copy migration tools
COPY --from=builder /app/bin/migrate ./bin/migrate

# Copy migration files
COPY --from=builder /app/adapters/db/migrations ./adapters/db/migrations

# Copy scripts
COPY --from=builder /app/scripts/entrypoint.sh ./entrypoint.sh
RUN chmod +x ./entrypoint.sh

# Copy config files
# COPY .env .env

# Create logs directory
RUN mkdir -p logs

# Expose port
EXPOSE 8080

# Use entrypoint script
ENTRYPOINT ["./entrypoint.sh"]
CMD ["./main"]
