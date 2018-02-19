package main

import (
	"github.com/jinzhu/configor"
	"gopkg.in/telegram-bot-api.v4"
	"log"
	"os/exec"
	"strings"
)

type Configuration struct {
	BotToken       string `env:"BOT_TOKEN"`
	AllowedUserIDs []string
	BotSubnet      string
}

var Config = Configuration{}

func main() {

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	err := configor.Load(&Config, "config.json")
	if err != nil {
		log.Panic(err)
	}

	if len(Config.AllowedUserIDs) == 0 {
		log.Panic("No allowed users in configuration.",
			"Set AllowedUserIDs in config.json")
	}

	bot, err := tgbotapi.NewBotAPI(Config.BotToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if strings.HasPrefix(update.Message.Text, "/whoishome") {
			out, err := exec.Command("nmap", "-sP", Config.BotSubnet).Output()
			if err != nil {
				log.Printf(err.Error())
			}

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, string(out))
			msg.ReplyToMessageID = update.Message.MessageID
			bot.Send(msg)

		} else {

			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			msg.ReplyToMessageID = update.Message.MessageID

			bot.Send(msg)
		}
	}
}
