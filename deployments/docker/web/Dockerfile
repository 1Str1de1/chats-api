FROM golang:alpine AS builder

WORKDIR /chats-api

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o /bin/chats-api ./cmd/app

FROM alpine

WORKDIR /app

COPY --from=builder /bin/chats-api /app/chats-api

COPY --from=builder /chats-api/migrations /app/migrations

EXPOSE 80

ENTRYPOINT ["/app/chats-api"]


