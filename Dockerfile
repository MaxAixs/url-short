FROM golang:1.22 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o url-short ./cmd/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/url-short .

COPY --from=builder /app/internal/config /root/internal/config

ENV DB_PASSWORD=${DB_PASSWORD}

CMD ["/root/url-short"]
