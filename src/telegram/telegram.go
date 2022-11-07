package telegram

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/fabritsius/potd-telegrapher/src/telegraph"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func RunBot() {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_APITOKEN"))
	if err != nil {
		panic(err)
	}

	bot.Debug = true

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30

	updates := bot.GetUpdatesChan(updateConfig)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		reply := func(msgText string) {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgText)
			msg.ReplyToMessageID = update.Message.MessageID

			_, err := bot.Send(msg)
			if err != nil {
				fmt.Println(err)
			}
		}

		if update.Message.From.UserName != os.Getenv("BOT_ADMIN") {
			reply("Your are not in the list of admins, please don't talk to me")
			continue
		}

		if update.Message.IsCommand() {
			replyMsg, err := handleCommands(update.Message.Command())
			if err != nil {
				fmt.Println(err)
				continue
			}

			reply(replyMsg)
			continue
		}

		reply(update.Message.Text)
	}
}

func handleCommands(command string) (string, error) {
	if command == "potd" {
		return makePOTD()
	}

	return "", fmt.Errorf("got unsupported command: %s", command)
}

func makePOTD() (string, error) {
	today := time.Now().Format("2006-01-02")

	result, err := telegraph.MakeArticle(today)
	if err != nil {
		log.Fatal(err)
	}

	return result.URL, nil
}
