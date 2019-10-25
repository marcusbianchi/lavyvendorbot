package main

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var statesCad map[int64]string
var statesBus map[int64]string
var userState map[int64]string

type premise struct {
	description string
	topic       string
	tags        []string
}

var premiseByChat map[int64]premise

func main() {
	statesCad = make(map[int64]string)
	statesBus = make(map[int64]string)
	userState = make(map[int64]string)
	premiseByChat = make(map[int64]premise)

	bot, err := tgbotapi.NewBotAPI("")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = false

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {

		if update.Message == nil { // ignore any non-Message Updates
			continue
		}
		chatID := update.Message.Chat.ID
		MessageID := update.Message.MessageID
		msg := processMessage(update.Message.Text, chatID, MessageID)
		bot.Send(msg)
	}
}

func processMessage(message string, chatID int64, messageID int) tgbotapi.MessageConfig {
	var msg tgbotapi.MessageConfig
	u, _ := userState[chatID]
	if message == "/start" || message == "/help" {
		msg = tgbotapi.NewMessage(chatID, "/cadastrar Para cadastrar Novas Premissas\n/buscar Para buscar premissas")
		return msg
	}
	if message == "/cadastrar" {
		_, ok := statesCad[chatID]
		if u == "buscar" || !ok {
			statesCad[chatID] = "inicio"
		}
		msg = tgbotapi.NewMessage(chatID, "Digite o texto da Premissa.")

		userState[chatID] = "cadastrar"
	} else if message == "/buscar" {
		_, ok := statesBus[chatID]
		if u == "cadastrar" || !ok {
			statesBus[chatID] = "inicio"
		}
		msg = tgbotapi.NewMessage(chatID, "Qual o tópico?")

		userState[chatID] = "buscar"
	} else if u == "cadastrar" {
		msg = processCadMessage(message, chatID, messageID)
	} else {
		msg = tgbotapi.NewMessage(chatID, "Comando não encontrado!")
		msg.ReplyToMessageID = messageID
	}
	return msg
}

//msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Comando não encontrado!")
//msg.ReplyToMessageID = update.Message.MessageID
