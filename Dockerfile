FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o url-shortener ./cmd

FROM scratch

COPY --from=builder /app/url-shortener /url-shortener

ENTRYPOINT ["/url-shortener"]
