FROM golang:1.23.1 as builder

WORKDIR /app
COPY . .

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /service ./cmd

FROM alpine:latest

WORKDIR /
COPY --from=builder /service /service

CMD ["/service"]