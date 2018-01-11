// expense-tracker project main.go
package main

import (
	"encoding/json"
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

	db, err := NewDB(conf.Db.User + ":" + conf.Db.Password + "@tcp(" + conf.Db.Host + ":3306)/" + conf.Db.Schema)
	if err != nil {
		log.Panic(err)
	}
	env := Env{DB: db}

	bot, err := tgbotapi.NewBotAPI(conf.BotApi)
	if err != nil {
		log.Panic(err)
	}
	//bot.Debug = true

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

		if update.Message.Text == "/balance" {
			b, _ := GetBalance(env.DB, user.UserId)
			if err != nil {
				log.Panic(err)
			}

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, ConvertBalance(b.CurrentBalance))
			bot.Send(msg)
			continue
		}

		if update.Message.Text == "/log" {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, conf.NotImplemented)
			bot.Send(msg)
			continue
		}

		result := strings.Split(update.Message.Text, ",")
		if len(result) == 2 {
			if result[0] != "" && result[1] != "" {
				amount, err := strconv.ParseFloat(result[0], 64)
				if err != nil {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, err.Error())
					bot.Send(msg)
					continue
				}

				insertId, err := InsertAmount(env.DB, user.UserId, amount, strings.Trim(result[1], " "))
				if err != nil {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, err.Error())
					bot.Send(msg)
					continue
				}

				if insertId != 0 {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, conf.Added)
					msg.ReplyToMessageID = update.Message.MessageID
					bot.Send(msg)
				}
			}
		}
	}
}
