// expense-tracker project main.go
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"gopkg.in/telegram-bot-api.v4"
)

type Configuration struct {
	Db             Db
	BotApi         string
	Intro          string
	NotImplemented string
	Added          string
}

func main() {
	file, _ := os.Open("configuration.json")
	decoder := json.NewDecoder(file)
	conf := Configuration{}
	err := decoder.Decode(&conf)
	if err != nil {
		log.Println("error:", err)
	}

	db, err := NewDB(conf.Db.User + ":" + conf.Db.Password + "@tcp(" + conf.Db.Host + ":" + conf.Db.Port + ")/" + conf.Db.Schema)
	if err != nil {
		log.Panic(err)
	}
	env := Env{DB: db}

	bot, err := tgbotapi.NewBotAPI(conf.BotApi)
	if err != nil {
		log.Panic(err)
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message == nil {
			continue
		}

		// get user data if not exists then create
		user, err := GetUserById(env.DB, update.Message.From.ID)
		if err != nil {
			// TODO: refactor this
			if err.Error() == "sql: no rows in result set" {
				_, err := InsertUser(env.DB, update.Message.From.ID)
				if err != nil {
					log.Panic(err)
				}

				msg := tgbotapi.NewMessage(update.Message.Chat.ID, conf.Intro)
				bot.Send(msg)
				continue
			}

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, err.Error())
			bot.Send(msg)
			log.Panic(err)
		}

		switch update.Message.Text {
		case "/balance":
			b, _ := GetBalance(env.DB, user.UserId)
			if err != nil {
				log.Panic(err)
			}

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, ConvertBalance(b.CurrentBalance))
			msg.ParseMode = "markdown"
			bot.Send(msg)
			continue
		case "/log":
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, conf.NotImplemented)
			msg.ParseMode = "markdown"
			bot.Send(msg)
			continue
		}

		result := strings.Split(update.Message.Text, ",")
		if len(result) == 2 {
			if result[0] != "" && result[1] != "" {
				// convert amount to float
				amount, err := strconv.ParseFloat(result[0], 64)
				if err != nil {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, err.Error())
					bot.Send(msg)
					continue
				}

				// insert amount into db and throw error if it was unsuccessful
				insertId, err := InsertAmount(env.DB, user.UserId, amount, strings.Trim(result[1], " "))
				if err != nil {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, err.Error())
					bot.Send(msg)
					continue
				}

				// send reply message with status and balance
				if insertId != 0 {
					// get balance
					b, _ := GetBalance(env.DB, user.UserId)
					if err != nil {
						log.Panic(err)
					}

					msgToReply := fmt.Sprintf(conf.Added, ConvertBalance(b.CurrentBalance))
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgToReply)
					msg.ReplyToMessageID = update.Message.MessageID
					msg.ParseMode = "markdown"

					bot.Send(msg)
				}
			}
		}
	}
}
