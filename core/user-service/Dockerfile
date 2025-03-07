FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

RUN go build -o user-service ./core/user-service/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/user-service .

COPY ./core/user-service/config.yaml .

EXPOSE 8081

CMD ["./user-service"]

