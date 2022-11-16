include .env
export $(shell sed 's/=.*//' .env)

wikipedia:
	go run ./cmd/get-wikipedia

telegraph:
	go run ./cmd/make-page

bot:
	go run ./cmd/telegram-bot

post-to-channel:
	go run ./cmd/post-to-channel
