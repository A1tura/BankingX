FROM golang:1.24 AS builder

WORKDIR /app/
COPY /microservices/2fa/ /app/service
RUN rm /app/service/go.mod
RUN mv /app/service/go.mod.docker /app/service/go.mod

COPY /modules/middlewares /app/middlewares
COPY /modules/db /app/db
COPY /modules/error /app/error
COPY /modules/auth /app/auth

WORKDIR /app/service

RUN go mod tidy
RUN GOARCH=amd64 GOOS=linux go build -o 2fa-service .

# Build the final image
FROM amd64/alpine:latest

RUN apk add --no-cache libc6-compat

WORKDIR /root/

COPY --from=builder /app/service/2fa-service .

CMD ["./2fa-service"]

