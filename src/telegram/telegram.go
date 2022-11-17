package telegram

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/fabritsius/potd-telegrapher/src/telegraph"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func RunBot() {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_TOKEN"))
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
				log.Printf("error: %v\n", err)
			}
		}

		if !userAllowed(update.Message.From) {
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

		if update.ChatJoinRequest != nil {
			if !userAllowed(&update.ChatJoinRequest.From) {
				if _, err = bot.Send(tgbotapi.DeclineChatJoinRequest{
					ChatConfig: update.FromChat().ChatConfig(),
					UserID:     update.ChatJoinRequest.From.ID,
				}); err != nil {
					log.Printf("error: %v\n", err)
				}
			} else {
				if _, err = bot.Send(tgbotapi.ApproveChatJoinRequestConfig{
					ChatConfig: update.FromChat().ChatConfig(),
					UserID:     update.ChatJoinRequest.From.ID,
				}); err != nil {
					log.Printf("error: %v\n", err)
				}
			}

			continue
		}
	}
}

func PostTodayArticle() {
	today := time.Now().Format("2006-01-02")
	PostArticle(today)
}

func PostArticle(date string) {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_TOKEN"))
	if err != nil {
		log.Fatalln(err)
	}

	channelsString, found := os.LookupEnv("CHANNEL_IDS")
	if !found {
		log.Fatalln("CHANNEL_IDS env variable is empty")
	}

	articleURL, err := makePOTD(date)
	if err != nil {
		log.Printf("error: %v\n", err)
		return
	}

	channels := parseChannels(channelsString)

	for _, channel := range channels {
		chanID, err := strconv.ParseInt(channel, 10, 64)
		if err != nil {
			continue
		}

		msg := tgbotapi.NewMessage(chanID, articleURL)
		_, err = bot.Send(msg)
		if err != nil {
			log.Printf("error: %v", err)
		}
	}
}

func userAllowed(user *tgbotapi.User) bool {
	return user.UserName == os.Getenv("BOT_ADMIN")
}

func handleCommands(command string) (string, error) {
	if command == "potd" {
		today := time.Now().Format("2006-01-02")
		return makePOTD(today)
	}

	return "", fmt.Errorf("got unsupported command: %s", command)
}

func makePOTD(date string) (string, error) {
	result, err := telegraph.MakeArticle(date)
	if err != nil {
		return "", err
	}

	return result.URL, nil
}

func parseChannels(channels string) []string {
	return strings.Split(strings.ReplaceAll(channels, " ", ""), ",")
}
