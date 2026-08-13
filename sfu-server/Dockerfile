FROM golang:1.24.4-alpine AS builder

WORKDIR /app

COPY go.mod go.sum* ./
RUN go mod download || go mod tidy && go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/sfu main.go

FROM alpine:latest
WORKDIR /app

COPY --from=builder /app/bin/sfu .

# The SFU uses ephemeral UDP ports for WebRTC ICE negotiation by default.
# It doesn't explicitly bind to a single TCP port, but communicates with the backend via WebSocket.

CMD ["./sfu"]
