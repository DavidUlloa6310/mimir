# Dockerfile for Go backend
FROM golang:1.23-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o /goapp

FROM alpine:3.17
WORKDIR /
COPY --from=builder /goapp /goapp

EXPOSE 8080

CMD ["/goapp"]
