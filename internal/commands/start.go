package commands

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	tgbot "github.com/wawan93/bot-framework"

	"no/internal/repo"
)

func SelectRegion(users *repo.UserRepo, ticks *repo.TickRepo, cities *repo.CityRepo) tgbot.CommonHandler {
	return func(bot *tgbot.BotFramework, update *tgbotapi.Update) error {
		bot.AnswerCallbackQuery(tgbotapi.CallbackConfig{
			CallbackQueryID: update.CallbackQuery.ID,
		})
		if update.Message != nil && update.Message.Chat != nil && !update.Message.Chat.IsPrivate() {
			msg := tgbotapi.NewMessage(bot.GetChatID(update), "Бот работает только в личке")
			msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonURL(
						"Перейти в личку",
						"https://t.me/"+bot.Self.UserName,
					),
				),
			)
			if _, err := bot.Send(msg); err != nil {
				log.Println(err)
			}
			return nil
		}

		bot.RegisterPlainTextHandler(SaveRegion(users, ticks, cities), bot.GetChatID(update))

		c, err := cities.Regions()
		if err != nil {
			return err
		}

		kb := tgbotapi.NewReplyKeyboard()
		for i := range c {
			kb.Keyboard = append(kb.Keyboard, tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton(c[i].Region),
			))
		}
		kb.OneTimeKeyboard = true

		msg := tgbotapi.NewMessage(bot.GetChatID(update), "Выберите регион")
		msg.ReplyMarkup = kb

		_, err = bot.Send(msg)
		return err
	}
}

func SaveRegion(users *repo.UserRepo, ticks *repo.TickRepo, cities *repo.CityRepo) tgbot.CommonHandler {
	return func(bot *tgbot.BotFramework, update *tgbotapi.Update) error {
		region := update.Message.Text
		bot.RegisterPlainTextHandler(SaveCity(users, ticks, cities, region), bot.GetChatID(update))

		text := "Выберите регион"
		if region == "Москва" || region == "Санкт-Петербург" {
			text = "Выберите район"
		}

		bot.RegisterPlainTextHandler(SaveCity(users, ticks, cities, region), bot.GetChatID(update))

		c, err := cities.Cities(region)
		if err != nil {
			return err
		}

		kb := tgbotapi.NewReplyKeyboard()
		for i := range c {
			kb.Keyboard = append(kb.Keyboard, tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton(c[i].Name),
			))
		}

		kb.OneTimeKeyboard = true

		msg := tgbotapi.NewMessage(bot.GetChatID(update), text)
		msg.ReplyMarkup = kb

		_, err = bot.Send(msg)
		return err
	}
}

func SaveCity(users *repo.UserRepo, ticks *repo.TickRepo, cities *repo.CityRepo, region string) tgbot.CommonHandler {
	return func(bot *tgbot.BotFramework, update *tgbotapi.Update) error {
		user, err := users.Get(bot.GetChatID(update))
		if err != nil {
			return err
		}

		name := update.Message.Text
		city, err := cities.Find(name, region)
		if err != nil {
			return err
		}

		user.City = *city
		user.CityID = city.ID

		if err := users.Update(user); err != nil {
			return err
		}

		return Start(users, ticks, cities)(bot, update)
	}
}

func Start(users *repo.UserRepo, ticks *repo.TickRepo, cities *repo.CityRepo) tgbot.CommonHandler {
	return func(bot *tgbot.BotFramework, update *tgbotapi.Update) error {

		_, err := users.Get(bot.GetChatID(update))
		if err != nil {
			log.Println(err)
		}

		bot.RegisterCallbackQueryHandler(StartWatermark(users), "watermark", bot.GetChatID(update))
		bot.RegisterCallbackQueryHandler(StartAPM(users, ticks, cities), "apm_main", bot.GetChatID(update))
		bot.RegisterCallbackQueryHandler(StartCoordinate(users, ticks, cities), "coordinate", bot.GetChatID(update))

		text := `Привет, я бот кампании «НЕТ»
Я помогу тебе:
`
		msg := tgbotapi.NewMessage(bot.GetChatID(update), text)
		msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(
					"Поменять аватарку",
					"watermark",
				),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(
					"Расклеить стикеры и листовки",
					"apm_main",
				),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(
					"Стать координатором в городе",
					"coordinate",
				),
			),
		)
		if _, err := bot.Send(msg); err != nil {
			log.Println(err)
		}

		return nil
	}
}
