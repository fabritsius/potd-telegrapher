include .env
export $(shell sed 's/=.*//' .env)

wikipedia:
	go run ./cmd/get-wikipedia

telegraph:
	go run ./cmd/make-page

bot:
	go run ./cmd/telegram-bot
