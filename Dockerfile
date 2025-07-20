FROM golang:1.24.4-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main

FROM alpine:latest
WORKDIR /app

RUN apk add --no-cache postgresql-client netcat-openbsd

COPY --from=builder /app/main /app/main
COPY --from=builder /app/migrations /app/migrations
COPY wait-for.sh /app/wait-for.sh
RUN chmod +x /app/wait-for.sh

EXPOSE 8080
CMD ["./wait-for.sh", "db", "5432", "./main"]