
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy go.mod files
COPY go.mod ./

# Copy source code
COPY cmd/ ./cmd/
COPY internal/ ./internal/

# Build the application
RUN go build -o url-shortener ./cmd/server

# Final stage
FROM alpine:latest

RUN apk add --no-cache ca-certificates

WORKDIR /root/

COPY --from=builder /app/url-shortener .

EXPOSE 8080

CMD ["./url-shortener"]