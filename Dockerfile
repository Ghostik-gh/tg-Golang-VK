FROM golang:1.20.2

WORKDIR /app
COPY config.json go.mod go.sum ./

RUN go mod download

RUN mkdir cmd
ADD cmd ./cmd/

RUN mkdir internal
ADD internal ./internal

EXPOSE 80/tcp


RUN CGO_ENABLED=0 GOOS=linux go build -o /tg-bot-golang ./cmd/tg-bot/main.go

CMD ["/tg-bot-golang"]

