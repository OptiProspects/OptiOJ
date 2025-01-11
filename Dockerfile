FROM golang:1.23-alpine AS builder

WORKDIR /OptiOJ
COPY . .

RUN apk add --no-cache \
    gcc \
    musl-dev \
    libc-dev

RUN go mod download
RUN CGO_ENABLED=1 GOOS=linux go build -o OptiOJ .

FROM alpine:latest

WORKDIR /OptiOJ
COPY --from=builder /OptiOJ/OptiOJ .
COPY --from=builder /OptiOJ/sql ./sql

EXPOSE 2550

CMD ["./OptiOJ"]
