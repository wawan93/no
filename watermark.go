package main

import (
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"net/http"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	tgbot "github.com/wawan93/bot-framework"
)

func Watermark(bot *tgbot.BotFramework, update *tgbotapi.Update) error {
	photos := *update.Message.Photo

	fileID := photos[len(photos)-1].FileID

	url, err := bot.GetFileDirectURL(fileID)
	if err != nil {
		return err
	}

	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	mark, err := os.Open("marks/box.png")
	if err != nil {
		return err
	}
	defer mark.Close()

	first, err := jpeg.Decode(response.Body)
	if err != nil {
		return fmt.Errorf("failed to decode: %s", err)
	}

	second, err := png.Decode(mark)
	if err != nil {
		return fmt.Errorf("failed to decode: %s", err)
	}

	offset := image.Pt(0, 0)
	b := first.Bounds()
	image3 := image.NewRGBA(b)
	draw.Draw(image3, b, first, image.ZP, draw.Src)
	draw.Draw(image3, second.Bounds().Add(offset), second, image.Point{}, draw.Over)

	file, err := ioutil.TempFile("/tmp", "*.jpeg")
	if err != nil {
		return err
	}
	defer os.Remove(file.Name())

	jpeg.Encode(file, image3, &jpeg.Options{jpeg.DefaultQuality})

	msg := tgbotapi.NewPhotoUpload(bot.GetChatID(update), file.Name())

	if _, err := bot.Send(msg); err != nil {
		return err
	}

	return nil
}
