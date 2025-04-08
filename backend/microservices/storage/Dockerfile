FROM golang:1.24 AS builder

WORKDIR /app/
COPY /microservices/storage/ /app/service
RUN rm /app/service/go.mod
RUN mv /app/service/go.mod.docker /app/service/go.mod

COPY /modules/db /app/db
COPY /modules/middlewares /app/middlewares
COPY /modules/error /app/error

WORKDIR /app/service

RUN go mod tidy
RUN GOARCH=amd64 GOOS=linux go build -o storage-service .

# Build the final image
FROM amd64/alpine:latest

RUN apk add --no-cache libc6-compat

WORKDIR /root/

COPY --from=builder /app/service/storage-service .

CMD ["./storage-service"]

