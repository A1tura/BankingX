FROM golang:1.24 AS builder

WORKDIR /app/
COPY . /app/service

WORKDIR /app/service

RUN go mod tidy
RUN GOARCH=amd64 GOOS=linux go build -o docs-service .

# Build the final image
FROM amd64/alpine:latest

RUN apk add --no-cache libc6-compat

WORKDIR /root/

COPY --from=builder /app/service/src ./src/
COPY --from=builder /app/service/docs-service .

CMD ["./docs-service"]

