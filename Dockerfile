FROM golang:1.22.4-alpine3.19 AS builder

WORKDIR /wallet
RUN apk add --no-cache git

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o /wallet /wallet/cmd/server/main.go

FROM alpine:3.18 AS app

COPY --from=builder /wallet/main /wallet/main

CMD ["/wallet/main"]