FROM golang:1.24 AS builder

WORKDIR /app

RUN go install github.com/swaggo/swag/cmd/swag@latest
ENV PATH="/go/bin:$PATH"

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN swag init --parseDependency --parseInternal

RUN CGO_ENABLED=0 go build -o main .

# --- Final Stage ---
FROM alpine:latest

WORKDIR /app

RUN apk update && apk add --no-cache ca-certificates

COPY --from=builder /app/main .
COPY --from=builder /app/.env .
COPY --from=builder /app/docs ./docs

EXPOSE 8000

CMD ["./main"]
