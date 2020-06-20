package commands

import (
	"log"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	tgbot "github.com/wawan93/bot-framework"

	"no/internal/img"
	"no/internal/repo"
)

func Start(users *repo.UserRepo) tgbot.CommonHandler {
	return func(bot *tgbot.BotFramework, update *tgbotapi.Update) error {
		if !update.Message.Chat.IsPrivate() {
			msg := tgbotapi.NewMessage(bot.GetChatID(update), "Бот работает только в личке")
			if _, err := bot.Send(msg); err != nil {
				log.Println(err)
			}
			return nil
		}

		user, err := users.Get(bot.GetChatID(update))

		photosCfg := tgbotapi.NewUserProfilePhotos(int(bot.GetChatID(update)))

		photos, err := bot.BotAPI.GetUserProfilePhotos(photosCfg)

		if err != nil || photos.TotalCount == 0 {
			msg := tgbotapi.NewMessage(bot.GetChatID(update), "Отправьте картинку")
			_, err := bot.Send(msg)
			return err
		}

		mainPhotos := photos.Photos[0]

		url, err := bot.GetFileDirectURL(mainPhotos[len(mainPhotos)-1].FileID)
		if err != nil {
			return err
		}

		// generate new picture
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

		if _, err = bot.Send(msg); err != nil {
			return err
		}

		return users.IncrementPhotos(user)
	}
}
