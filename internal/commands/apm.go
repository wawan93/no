package commands

import (
	"context"
	"no/internal/models"
	"no/internal/repo"
	"strconv"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	tgbot "github.com/wawan93/bot-framework"
)

func StartAPM(users *repo.UserRepo, ticks *repo.TickRepo, cities *repo.CityRepo) tgbot.CommonHandler {
	return func(bot *tgbot.BotFramework, update *tgbotapi.Update) error {
		bot.AnswerCallbackQuery(tgbotapi.CallbackConfig{
			CallbackQueryID: update.CallbackQuery.ID,
		})
		user, err := users.Get(bot.GetChatID(update))
		if err != nil {
			return err
		}

		if user.CityID == 0 {
			return SelectRegion(users, ticks, cities)(bot, update)
		}

		bot.RegisterCallbackQueryHandler(GetAPM(users, ticks), "apm_get_from_coordinator", bot.GetChatID(update))
		bot.RegisterCallbackQueryHandler(AskAPMCount(users, ticks, cities), "apm_start", bot.GetChatID(update))
		bot.RegisterCallbackQueryHandler(SelectRegion(users, ticks, cities), "change_region", bot.GetChatID(update))

		msg := tgbotapi.NewMessage(bot.GetChatID(update), "Ты можешь распечатать стикеры и листовки сам, а можешь взять уже готовые")
		kb := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonURL(
					"Напечатать самому",
					"https://drive.google.com/drive/folders/1xBLVab1GJSrEaPR1bvXNjEIfacUi0rgh",
				),
			),
		)

		if user.City.Coordinator != "" {
			kb.InlineKeyboard = append(kb.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(
					"Взять листовки у координатора",
					"apm_get_from_coordinator",
				),
			),
			)
		}
		kb.InlineKeyboard = append(kb.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				"Я начал раздавать",
				"apm_start",
			),
		),
		)
		kb.InlineKeyboard = append(kb.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				"Сменить район",
				"change_region",
			),
		),
		)

		msg.ReplyMarkup = kb

		if _, err := bot.Send(msg); err != nil {
			return err
		}

		return nil
	}
}

func AskAPMCount(users *repo.UserRepo, ticks *repo.TickRepo, cities *repo.CityRepo) tgbot.CommonHandler {
	return func(bot *tgbot.BotFramework, update *tgbotapi.Update) error {
		bot.AnswerCallbackQuery(tgbotapi.CallbackConfig{
			CallbackQueryID: update.CallbackQuery.ID,
		})
		bot.Send(tgbotapi.NewMessage(bot.GetChatID(update), "Сколько материалов вы взяли (напечатали)"))
		bot.RegisterPlainTextHandler(SaveAPMCount(users, ticks, cities), bot.GetChatID(update))
		return nil
	}
}

func SaveAPMCount(users *repo.UserRepo, ticks *repo.TickRepo, cities *repo.CityRepo) tgbot.CommonHandler {
	return func(bot *tgbot.BotFramework, update *tgbotapi.Update) error {
		user, err := users.Get(bot.GetChatID(update))
		if err != nil {
			return err
		}

		material, err := strconv.Atoi(update.Message.Text)
		if err != nil {
			return err
		}

		user.Materials += uint(material)
		users.Update(user)

		bot.Send(tgbotapi.NewMessage(bot.GetChatID(update), "Сохранено. Можете раздавать"))
		return StartDistributionAPM(users, ticks, cities)(bot, update)
	}
}

func StartDistributionAPM(users *repo.UserRepo, ticks *repo.TickRepo, cities *repo.CityRepo) tgbot.CommonHandler {
	return func(bot *tgbot.BotFramework, update *tgbotapi.Update) error {
		chatID := bot.GetChatID(update)
		ctx, cancel := context.WithCancel(context.Background())

		bot.RegisterCommand("❌Отмена", func(bot *tgbot.BotFramework, update *tgbotapi.Update) error {
			cancel()
			return Start(users, ticks, cities)(bot, update)
		}, chatID)

		resetCh := make(chan struct{})

		go func() {
			for {
				timeout := time.After(15 * time.Minute)
				select {
				case <-ctx.Done():
					break
				case <-resetCh:
					continue
				case <-timeout:
					bot.Send(tgbotapi.NewMessage(
						chatID,
						"Вы расклеиваете уже 15 минут. Пожалуйста, скиньте location и фото!",
					))
				}
			}
		}()

		bot.RegisterLocationHandler(TickLocation(users, ticks, resetCh), bot.GetChatID(update))

		text := "Отпраьте Location, где вы наклеили листовки и стикеры, чтобы отметить этот дом на карте"
		msg := tgbotapi.NewMessage(bot.GetChatID(update), text)

		kb := tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("❌Отмена"),
			),
		)
		kb.OneTimeKeyboard = true

		msg.ReplyMarkup = kb

		_, err := bot.Send(msg)
		return err
	}
}

func TickLocation(users *repo.UserRepo, ticks *repo.TickRepo, resetCh chan struct{}) tgbot.CommonHandler {
	return func(bot *tgbot.BotFramework, update *tgbotapi.Update) error {
		user, err := users.Get(bot.GetChatID(update))
		if err != nil {
			return err
		}

		tick := &models.Tick{
			UserID:    user.ID,
			CityID:    user.CityID,
			Latitude:  update.Message.Location.Latitude,
			Longitude: update.Message.Location.Longitude,
		}

		if err := ticks.Save(tick); err != nil {
			return err
		}

		resetCh <- struct{}{}

		bot.RegisterPhotoHandler(TickPhoto(users, ticks, tick, resetCh), bot.GetChatID(update))

		text := "Отлично! Дом будет отмечен на карте! Теперь можете прислать фото наклеенной листовки или стикера для наших соцсетей"
		msg := tgbotapi.NewMessage(bot.GetChatID(update), text)

		_, err = bot.Send(msg)

		return err
	}
}

func TickPhoto(users *repo.UserRepo, ticks *repo.TickRepo, tick *models.Tick, resetCh chan struct{}) tgbot.CommonHandler {
	return func(bot *tgbot.BotFramework, update *tgbotapi.Update) error {
		photos := *update.Message.Photo
		photo := photos[len(photos)-1]

		url, err := bot.GetFileDirectURL(photo.FileID)
		if err != nil {
			return err
		}

		tick.Photo = url

		if err := ticks.Save(tick); err != nil {
			return err
		}

		fwd := tgbotapi.NewForward(-483425949, bot.GetChatID(update), update.Message.MessageID)
		bot.Send(fwd)

		text := "Отлично! Фото будет опубликовано в наших соцсетях."
		msg := tgbotapi.NewMessage(bot.GetChatID(update), text)

		if _, err := bot.Send(msg); err != nil {
			return err
		}

		resetCh <- struct{}{}

		bot.RegisterPlainTextHandler(TickCount(users, ticks, tick, resetCh), bot.GetChatID(update))

		text = "Отметьте, сколько вы расклеили"
		msg = tgbotapi.NewMessage(bot.GetChatID(update), text)
		_, err = bot.Send(msg)
		return err
	}
}

func TickCount(users *repo.UserRepo, ticks *repo.TickRepo, tick *models.Tick, resetCh chan struct{}) tgbot.CommonHandler {
	return func(bot *tgbot.BotFramework, update *tgbotapi.Update) error {
		materials, err := strconv.Atoi(update.Message.Text)
		if err != nil {
			return err
		}

		tick.Materials += uint(materials)

		if err := ticks.Save(tick); err != nil {
			return err
		}

		text := "Сохранено. Отметьте следующий дом:"
		msg := tgbotapi.NewMessage(bot.GetChatID(update), text)
		if _, err := bot.Send(msg); err != nil {
			return err
		}

		resetCh <- struct{}{}

		text = "Отпраьте Location, где вы наклеили листовки и стикеры, чтобы отметить этот дом на карте"
		msg = tgbotapi.NewMessage(bot.GetChatID(update), text)
		_, err = bot.Send(msg)
		return err
	}
}

func GetAPM(users *repo.UserRepo, ticks *repo.TickRepo) tgbot.CommonHandler {
	return func(bot *tgbot.BotFramework, update *tgbotapi.Update) error {
		bot.AnswerCallbackQuery(tgbotapi.CallbackConfig{
			CallbackQueryID: update.CallbackQuery.ID,
		})
		user, err := users.Get(bot.GetChatID(update))
		if err != nil {
			return err
		}

		text := "Напишите координатору @" + user.City.Coordinator
		msg := tgbotapi.NewMessage(bot.GetChatID(update), text)

		_, err = bot.Send(msg)
		return err
	}
}
