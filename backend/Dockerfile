FROM golang:1.24.4-alpine AS builder

WORKDIR /app

# Copy go mod and download dependencies
COPY go.mod go.sum* ./
# Since go.sum might not exist yet, use a wildcard or run go mod tidy
RUN go mod download || go mod tidy && go mod download

# Copy the rest of the application
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/api cmd/api/main.go

# Final stage
FROM alpine:latest
WORKDIR /app

# Copy the compiled binary and migrations directory
COPY --from=builder /app/bin/api .
COPY migrations ./migrations

EXPOSE 8080

CMD ["./api"]
