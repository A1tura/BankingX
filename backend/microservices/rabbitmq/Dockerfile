FROM golang:1.24 AS builder

WORKDIR /app
COPY . .

RUN go mod tidy
RUN GOARCH=amd64 GOOS=linux go build -o rabbitmq-service .

# Build the final image
FROM amd64/alpine:latest

RUN apk add --no-cache libc6-compat

WORKDIR /root/

COPY --from=builder /app/rabbitmq-service .

CMD ["./rabbitmq-service"]

