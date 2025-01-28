FROM golang:1.23-alpine AS builder

WORKDIR /OptiOJ
COPY . .

RUN apk add --no-cache \
    gcc \
    musl-dev \
    libc-dev

RUN go mod download

# use different build parameters for different architectures
ARG TARGETARCH
RUN if [ "$TARGETARCH" = "arm64" ]; then \
    CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o OptiOJ . ; \
    else \
    CGO_ENABLED=1 GOOS=linux go build -o OptiOJ . ; \
    fi

FROM alpine:latest

WORKDIR /OptiOJ
COPY --from=builder /OptiOJ/OptiOJ .
COPY --from=builder /OptiOJ/sql ./sql

EXPOSE 2550

CMD ["./OptiOJ"]
