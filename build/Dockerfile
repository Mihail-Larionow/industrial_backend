FROM golang:1.21-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app
COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/calculator ./cmd/api

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/calculator .

COPY config/main.yaml /app/config/main.yaml

EXPOSE 8080 9090

CMD ["./calculator"]
