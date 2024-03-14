FROM golang:1.22.0-alpine AS builder

COPY . /github.com/arivlav/auth/source/
WORKDIR /github.com/arivlav/auth/source/

RUN go mod download
RUN go build -o ./bin/auth_server cmd/grpc_server/main.go

FROM alpine:latest

WORKDIR /root/
COPY --from=builder /github.com/arivlav/auth/source/bin/auth_server .

CMD ["./auth_server"]