FROM golang:1.16 AS builder

WORKDIR /app

COPY go.mod .
COPY go.sum .
COPY vendor/ ./vendor/

ENV CGO_ENABLED=0
ENV GOOS=linux

RUN go build -mod=vendor -a -installsuffix cgo -o main ./cmd/api

FROM scratch

COPY --from=builder /app/main .

ENTRYPOINT ["/main"]
