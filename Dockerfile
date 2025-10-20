FROM golang:alpine AS builder

WORKDIR /build/

COPY go.mod go.sum /build/

RUN go mod download

COPY ./cmd /build/cmd
COPY ./internal /build/internal

RUN CGO_ENABLED=0 GOOS=linux go build /build/cmd/app/main.go

FROM alpine:3.16.8

WORKDIR /app/

COPY --from=builder /build/main /app/main