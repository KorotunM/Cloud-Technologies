# syntax=docker/dockerfile:1

FROM golang:1.23-alpine AS builder
WORKDIR /src
RUN apk add --no-cache ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /out/app ./cmd

FROM alpine:3.20 AS final
WORKDIR /app
RUN addgroup -S app && adduser -S app -G app && apk add --no-cache ca-certificates curl

COPY --from=builder /out/app /app/app
COPY --from=builder /src/web /app/web

ENV SERVER_ADDR=:8080
EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 CMD curl -f http://127.0.0.1:8080/healthz || exit 1

USER app
ENTRYPOINT ["/app/app"]
