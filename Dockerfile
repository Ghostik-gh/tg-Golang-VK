FROM golang:1.20.2

WORKDIR /app
COPY go.mod go.sum ./

RUN go mod download
COPY *.go config.json ./

EXPOSE 80/tcp

RUN CGO_ENABLED=0 GOOS=linux go build -o /tg-bot-golang .

CMD ["/tg-bot-golang"]

# go mod init vk_intership
# go get -u github.com/go-telegram-bot-api/telegram-bot-api/v5
# go run main.go