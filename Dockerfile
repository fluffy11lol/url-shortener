FROM golang:1.23-alpine AS builder

RUN apk add --no-cache gcc musl-dev sqlite-dev

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o url-shortener ./cmd/url-shortener

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/url-shortener ./url-shortener

COPY config/local.yaml ./config/local.yaml

COPY storage ./storage

EXPOSE 8081

CMD ["./url-shortener"]