package commands

import (
	"no/internal/models"
	"no/internal/repo"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	tgbot "github.com/wawan93/bot-framework"
)

func StartAPM(users *repo.UserRepo, ticks *repo.TickRepo) tgbot.CommonHandler {
	return func(bot *tgbot.BotFramework, update *tgbotapi.Update) error {
		msg := tgbotapi.NewMessage(bot.GetChatID(update), "Скачать макет листовки для печати")
		msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
			// tgbotapi.NewInlineKeyboardRow(
			// 	tgbotapi.NewInlineKeyboardButtonURL(
			// 		"Взять листовки у координатора",
			// 		"https://drive.google.com/drive/folders/1xBLVab1GJSrEaPR1bvXNjEIfacUi0rgh",
			// 	),
			// ),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonURL(
					"Напечатать самому",
					"https://drive.google.com/drive/folders/1xBLVab1GJSrEaPR1bvXNjEIfacUi0rgh",
				),
			),
		)
		if _, err := bot.Send(msg); err != nil {
			return err
		}

		bot.RegisterLocationHandler(TickLocation(users, ticks), bot.GetChatID(update))
		bot.RegisterPhotoHandler(TickPhoto(users, ticks), bot.GetChatID(update))
		bot.RegisterCommand("❌ Отмена", Start(users, ticks), bot.GetChatID(update))

		text := "Отпраьте Location, где вы наклеили листовки и стикеры, а также фото для наших соцсетей"
		msg = tgbotapi.NewMessage(bot.GetChatID(update), text)

		kb := tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButtonLocation("Send Location"),
				tgbotapi.NewKeyboardButton("❌ Отмена"),
			),
		)
		kb.OneTimeKeyboard = true

		msg.ReplyMarkup = kb

		_, err := bot.Send(msg)
		return err
	}
}

func TickLocation(users *repo.UserRepo, ticks *repo.TickRepo) tgbot.CommonHandler {
	return func(bot *tgbot.BotFramework, update *tgbotapi.Update) error {
		user, err := users.Get(bot.GetChatID(update))
		if err != nil {
			return err
		}

		tick := &models.Tick{
			User:      *user,
			Latitude:  update.Message.Location.Latitude,
			Longitude: update.Message.Location.Longitude,
		}

		if err := ticks.Save(tick); err != nil {
			return err
		}

		return nil
	}
}

func TickPhoto(users *repo.UserRepo, ticks *repo.TickRepo) tgbot.CommonHandler {
	return func(bot *tgbot.BotFramework, update *tgbotapi.Update) error {

		return nil
	}
}
