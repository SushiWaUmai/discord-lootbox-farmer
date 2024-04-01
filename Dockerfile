FROM golang:alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build

FROM alpine as runner

WORKDIR /app/

COPY --from=builder /app/discord-lootbox-farmer .

CMD ["./discord-lootbox-farmer"]
