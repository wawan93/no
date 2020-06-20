package commands

import (
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	tgbot "github.com/wawan93/bot-framework"

	"no/internal/img"
	"no/internal/repo"
)

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
