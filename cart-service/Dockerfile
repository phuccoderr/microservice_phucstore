FROM golang:1.22.5-alpine AS builder

WORKDIR /build

COPY . .

RUN go mod download

RUN go build -o cartService ./cmd/server/

FROM scratch

COPY ./config /config

COPY --from=builder /build/cartService /

ENTRYPOINT ["./cartService", "config/local.yaml"]