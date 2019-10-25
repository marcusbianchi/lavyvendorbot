package main

import (
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func processCadMessage(message string, chatID int64, messageID int) tgbotapi.MessageConfig {
	var msg tgbotapi.MessageConfig
	i, ok := statesCad[chatID]
	if !ok {
		msg = tgbotapi.NewMessage(chatID, "Comando não encontrado!")
		msg.ReplyToMessageID = messageID
		return msg
	}
	switch i {
	case "inicio":
		//salva texto Premissa
		p := premise{description: message}
		statesCad[chatID] = "confirma_texto"
		premiseByChat[chatID] = p
		msg = tgbotapi.NewMessage(chatID, "Este é o Texto da Premissa: \""+message+"\"\nSim ou não?")

	case "confirma_texto":
		if strings.TrimSpace(strings.ToLower(message)) == "sim" {
			statesCad[chatID] = "topico"
			msg = tgbotapi.NewMessage(chatID, "Qual é o tópico da premissa?\nProjeto ou Premissas Gerais")
		} else {
			msg = tgbotapi.NewMessage(chatID, "Digite o texto da Premissa.")
			statesCad[chatID] = "inicio"
		}

	case "topico":
		text := strings.TrimSpace(strings.ToLower(message))
		switch text {
		case "projeto", "premissas gerais":
			statesCad[chatID] = "confirma_topico"
			p := premiseByChat[chatID]
			p.topic = text
			premiseByChat[chatID] = p
			msg = tgbotapi.NewMessage(chatID, "Este é o Tópico da Premissa: \""+text+"\"\nSim ou não?")
		default:
			statesCad[chatID] = "topico"
			msg = tgbotapi.NewMessage(chatID, "Qual é o tópico da premissa?\nProjeto ou Premissas Gerais")
		}

	case "confirma_topico":
		if strings.TrimSpace(strings.ToLower(message)) == "sim" {
			statesCad[chatID] = "tags"
			msg = tgbotapi.NewMessage(chatID, "Digite as tags da premissa separadas por vírgulas (Ex:java,c#,CRM)")
		} else {
			statesCad[chatID] = "topico"
			msg = tgbotapi.NewMessage(chatID, "Qual é o tópico da premissa?\nProjeto ou Premissas Gerais")
		}
	case "tags":
		statesCad[chatID] = "confirma_tags"
		tags := strings.Split(strings.TrimSpace(strings.ToLower(message)), ", ")
		p := premiseByChat[chatID]
		p.tags = tags
		premiseByChat[chatID] = p
		msg = tgbotapi.NewMessage(chatID, "Estas são as tags da premissa: \""+strings.Join(tags, " ")+"\"\nSim ou não?")
	case "confirma_tags":
		if strings.TrimSpace(strings.ToLower(message)) == "sim" {
			statesCad[chatID] = "confirma_final"
			p := premiseByChat[chatID]
			mensagem := fmt.Sprint("Por favor confirme a criação da premisa:",
				"\nDescrição: ", p.description,
				"\nTópico: ", p.topic,
				"\nTags: ", strings.Join(p.tags, ","),
				"\nSim ou não?")
			msg = tgbotapi.NewMessage(chatID, mensagem)
		} else {
			statesCad[chatID] = "tags"
			msg = tgbotapi.NewMessage(chatID, "Digite as tags da premissa separadas por vírgulas (Ex:java,c#,CRM)")
		}
	case "confirma_final":
		if strings.TrimSpace(strings.ToLower(message)) == "sim" {
			statesCad[chatID] = "inicio"
			msg = tgbotapi.NewMessage(chatID, "/cadastrar Para cadastrar Novas Premissas\n/buscar Para buscar premissas")
		} else {
			statesCad[chatID] = "confirma_final"
			p := premiseByChat[chatID]
			mensagem := fmt.Sprint("Por favor confirme a criação da premisa:",
				"\nDescrição: ", p.description,
				"\nTópico: ", p.topic,
				"\nTags: ", strings.Join(p.tags, ","),
				"\nSim ou não?")
			msg = tgbotapi.NewMessage(chatID, mensagem)
		}
	}
	return msg
}
