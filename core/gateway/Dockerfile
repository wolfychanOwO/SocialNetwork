FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum .
RUN go mod download
COPY . .

ENV CGO_ENABLED=0
ENV GOOS=linux

RUN go build -o gateway ./core/gateway/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/gateway .

EXPOSE 8080

CMD ["./gateway"]

