FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY . .

RUN go build -o main .

FROM alpine:latest  

RUN apk add --no-cache ca-certificates

WORKDIR /root/

COPY --from=builder /app/main .

CMD ["./main"]
