FROM golang:1.24 AS builder

WORKDIR /app
COPY . .

RUN go mod tidy
RUN GOARCH=amd64 GOOS=linux go build -o email-service .

# Build the final image
FROM amd64/alpine:latest

RUN apk add --no-cache libc6-compat

WORKDIR /root/

COPY --from=builder /app/email-service .
RUN mkdir templates
COPY --from=builder /app/templates/templates.json /root/templates/

CMD ["./email-service"]

