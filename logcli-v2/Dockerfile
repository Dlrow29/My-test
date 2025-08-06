FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o cli ./log-cliv2.go

FROM alpine:latest
RUN apk add --no-cache ca-certificates
COPY --from=builder /app/cli /usr/local/bin/
ENTRYPOINT ["cli"]
