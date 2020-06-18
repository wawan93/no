package commands

import (
	"log"
	"no/internal/models"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/jinzhu/gorm"
	tgbot "github.com/wawan93/bot-framework"
)

func Start(db *gorm.DB) tgbot.CommonHandler {
	return func(bot *tgbot.BotFramework, update *tgbotapi.Update) error {
		if !update.Message.Chat.IsPrivate() {
			msg := tgbotapi.NewMessage(bot.GetChatID(update), "Бот работает только в личке")
			if _, err := bot.Send(msg); err != nil {
				log.Println(err)
			}
			return nil
		}

		var user models.User
		user.ChatID = bot.GetChatID(update)
		if err := db.Where("chat_id=?", user.ChatID).FirstOrCreate(&user).Error; err != nil {
			return err
		}

		photosCfg := tgbotapi.NewUserProfilePhotos(int(bot.GetChatID(update)))

		photos, err := bot.BotAPI.GetUserProfilePhotos(photosCfg)

		if err != nil || photos.TotalCount == 0 {
			msg := tgbotapi.NewMessage(bot.GetChatID(update), "Отправьте картинку")
			_, err := bot.Send(msg)
			return err
		}

		mainPhotos := photos.Photos[0]

		response, err := uploadedPhoto(bot, mainPhotos[len(mainPhotos)-1].FileID)
		if err != nil {
			return err
		}
		defer response.Body.Close()

		// open watermark
		mark, err := os.Open("marks/box.png")
		if err != nil {
			return err
		}
		defer mark.Close()

		// generate new picture
		buf, err := generate(response.Body, mark)
		if err != nil {
			return err
		}

		file := tgbotapi.FileReader{
			Name:   randomName() + ".jpeg",
			Reader: buf,
			Size:   -1,
		}

		msg := tgbotapi.NewPhotoUpload(bot.GetChatID(update), file)

		if _, err = bot.Send(msg); err != nil {
			return err
		}

		user.Photos++

		err = db.Save(&user).Error
		return err
	}
}
