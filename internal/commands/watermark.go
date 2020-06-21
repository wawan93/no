package commands

import (
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	tgbot "github.com/wawan93/bot-framework"

	"no/internal/img"
	"no/internal/repo"
)

func StartWatermark(users *repo.UserRepo) tgbot.CommonHandler {
	return func(bot *tgbot.BotFramework, update *tgbotapi.Update) error {
		photosCfg := tgbotapi.NewUserProfilePhotos(int(bot.GetChatID(update)))

		photos, err := bot.BotAPI.GetUserProfilePhotos(photosCfg)

		if err != nil || photos.TotalCount == 0 {
			msg := tgbotapi.NewMessage(bot.GetChatID(update), "Отправьте фото")
			_, err := bot.Send(msg)
			return err
		}

		mainPhotos := photos.Photos[0]

		url, err := bot.GetFileDirectURL(mainPhotos[len(mainPhotos)-1].FileID)
		if err != nil {
			return err
		}

		buf, err := img.Generate(url)
		if err != nil {
			return err
		}

		file := tgbotapi.FileReader{
			Name:   strconv.Itoa(int(bot.GetChatID(update))) + ".jpeg",
			Reader: buf,
			Size:   -1,
		}

		msg := tgbotapi.NewPhotoUpload(bot.GetChatID(update), file)

		if _, err := bot.Send(msg); err != nil {
			return err
		}

		user, err := users.Get(bot.GetChatID(update))
		if err != nil {
			return err
		}
		users.IncrementPhotos(user)

		msgSend := tgbotapi.NewMessage(bot.GetChatID(update), "Отправьте фото")
		_, err = bot.Send(msgSend)
		return err
	}
}

func Watermark(users *repo.UserRepo) tgbot.CommonHandler {
	return func(bot *tgbot.BotFramework, update *tgbotapi.Update) error {
		// show "...typing"
		typingMsg := tgbotapi.NewChatAction(bot.GetChatID(update), tgbotapi.ChatTyping)
		bot.Send(typingMsg)

		photos := *update.Message.Photo
		fileID := photos[len(photos)-1].FileID

		url, err := bot.GetFileDirectURL(fileID)
		if err != nil {
			return err
		}

		buf, err := img.Generate(url)
		if err != nil {
			return err
		}

		file := tgbotapi.FileReader{
			Name:   strconv.Itoa(int(bot.GetChatID(update))) + ".jpeg",
			Reader: buf,
			Size:   -1,
		}

		msg := tgbotapi.NewPhotoUpload(bot.GetChatID(update), file)

		if _, err := bot.Send(msg); err != nil {
			return err
		}

		user, err := users.Get(bot.GetChatID(update))
		if err != nil {
			return err
		}
		return users.IncrementPhotos(user)
	}
}
