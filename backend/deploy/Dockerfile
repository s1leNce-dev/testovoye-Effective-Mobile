# --- Build Stage ---
FROM golang:1.24 AS builder

WORKDIR /app

COPY ./backend/go.mod ./backend/go.sum ./
RUN go mod download

COPY ./backend/ ./

RUN CGO_ENABLED=0 go build -o main ./cmd/server_app/main.go

# --- Final Stage ---
FROM alpine:latest

WORKDIR /app
RUN apk update && apk add --no-cache ca-certificates

COPY --from=builder /app/main .
COPY --from=builder /app/.env .

EXPOSE 8000
CMD ["./main"]
