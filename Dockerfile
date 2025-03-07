
FROM golang:1.24-alpine AS base
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

FROM base AS tester
CMD ["go", "test", "./..."]

FROM base AS builder
RUN CGO_ENABLED=0 GOOS=linux go build -o url-short ./cmd/main.go

FROM alpine:latest AS final
WORKDIR /root/
COPY --from=builder /app/url-short ./
COPY --from=builder /app/internal/config /root/internal/config
ENV DB_PASSWORD=${DB_PASSWORD}
CMD ["/root/url-short"]
