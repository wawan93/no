package commands

import (
	"fmt"
	"no/internal/repo"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	tgbot "github.com/wawan93/bot-framework"
)

func StartCoordinate(users *repo.UserRepo, ticks *repo.TickRepo, cities *repo.CityRepo) tgbot.CommonHandler {
	return func(bot *tgbot.BotFramework, update *tgbotapi.Update) error {
		user, err := users.Get(bot.GetChatID(update))
		if err != nil {
			return err
		}

		if user.CityID == 0 {
			return SelectRegion(users, ticks, cities)(bot, update)
		}

		text := fmt.Sprintf(
			"[user](tg://user?id=%d) хочет стать координатором в %s, %s",
			user.ChatID,
			user.City.Name,
			user.City.Region,
		)
		msg := tgbotapi.NewMessage(-439649564, text)
		msg.ParseMode = tgbotapi.ModeMarkdown

		if _, err := bot.Send(msg); err != nil {
			return err
		}

		replyMsg := tgbotapi.NewMessage(bot.GetChatID(update), "Заявка отправлена")
		_, err = bot.Send(replyMsg)
		return err
	}
}
