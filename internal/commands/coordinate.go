package commands

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	tgbot "github.com/wawan93/bot-framework"
)

func StartCoordinate(bot *tgbot.BotFramework, update *tgbotapi.Update) error {
	msg := tgbotapi.NewMessage(bot.GetChatID(update), "TODO: implement")

	_, err := bot.Send(msg)
	return err
}
