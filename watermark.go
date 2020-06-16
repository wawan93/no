package main

import (
	"bytes"
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/nfnt/resize"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	tgbot "github.com/wawan93/bot-framework"
)

func Watermark(bot *tgbot.BotFramework, update *tgbotapi.Update) error {
	// show "...typing"
	typingMsg := tgbotapi.NewChatAction(bot.GetChatID(update), tgbotapi.ChatTyping)
	bot.Send(typingMsg)

	photos := *update.Message.Photo
	fileID := photos[len(photos)-1].FileID

	// read sent photo
	response, err := uploadedPhoto(bot, fileID)
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

	if _, err := bot.Send(msg); err != nil {
		return err
	}

	return nil
}

func uploadedPhoto(bot *tgbot.BotFramework, fileID string) (*http.Response, error) {
	url, err := bot.GetFileDirectURL(fileID)
	if err != nil {
		return nil, err
	}

	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func generate(background, watermark io.Reader) (io.Reader, error) {
	first, err := jpeg.Decode(background)
	if err != nil {
		return nil, fmt.Errorf("failed to decode: %s", err)
	}

	mark, err := png.Decode(watermark)
	if err != nil {
		return nil, fmt.Errorf("failed to decode: %s", err)
	}

	b := first.Bounds()

	mark = resizeMark(mark, b)

	offsetY := b.Dy() - mark.Bounds().Dy()
	offset := image.Pt(0, offsetY)

	image3 := image.NewRGBA(b)
	draw.Draw(image3, b, first, image.Point{}, draw.Src)
	draw.Draw(image3, mark.Bounds().Add(offset), mark, image.Point{}, draw.Over)

	var buf bytes.Buffer
	jpeg.Encode(&buf, image3, &jpeg.Options{jpeg.DefaultQuality})

	return &buf, nil
}

func resizeMark(img image.Image, b image.Rectangle) image.Image {
	var w, h uint

	if b.Dx() > b.Dy() {
		h = uint(b.Dy() / 2)
		w = h
	}

	if b.Dx() <= b.Dy() {
		w = uint(b.Dx() / 2)
		h = w
	}

	return resize.Resize(w, h, img, resize.Lanczos2)
}

func randomName() string {
	rand.Seed(time.Now().UnixNano())

	u, _ := uuid.NewRandom()

	return u.String()
}
