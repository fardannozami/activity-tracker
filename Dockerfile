# syntax=docker/dockerfile:1

FROM golang:1.25-alpine AS builder
WORKDIR /src
RUN apk add --no-cache git ca-certificates
COPY go.mod ./
RUN go mod download
COPY . .
RUN go install github.com/swaggo/swag/cmd/swag@v1.16.6
RUN /go/bin/swag init -g cmd/api/main.go -o internal/docs --parseInternal
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o /out/api ./cmd/api

FROM alpine:3.20
WORKDIR /app
RUN apk add --no-cache ca-certificates postgresql-client tzdata
COPY --from=builder /out/api /app/api
COPY migrations /app/migrations
EXPOSE 8080
CMD ["/app/api"]
