package main

import (
	"log"
	"net/http"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	tgbot "github.com/wawan93/bot-framework"

	"no/internal/commands"
	"no/internal/db"
)

func main() {
	db.Connect(
		"mysql",
		os.Getenv("DBHOST"),
		os.Getenv("DBPORT"),
		os.Getenv("DBUSER"),
		os.Getenv("DBPASS"),
		os.Getenv("DBNAME"),
	)
	go db.Migrate()

	token := os.Getenv("TOKEN")

	webhookAddress := os.Getenv("WEBHOOK_ADDRESS")
	if webhookAddress == "" {
		log.Panic("WEBHOOK_ADDRESS is empty")
	}

	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatal(err)
	}

	api.Debug = true

	log.Printf("logged in as %v", api.Self.UserName)

	bot := tgbot.NewBotFramework(api)

	updates := getUpdatesChannel(api, webhookAddress)

	if err := bot.RegisterCommand("/start", commands.Start(db.Conn), 0); err != nil {
		log.Fatalf("can't register command: %+v", err)
	}

	if err := bot.RegisterPhotoHandler(commands.Watermark(db.Conn), 0); err != nil {
		log.Fatalf("can't register handler: %+v", err)
	}

	bot.HandleUpdates(updates)
}

func getUpdatesChannel(api *tgbotapi.BotAPI, webhookAddress string) tgbotapi.UpdatesChannel {
	var updates tgbotapi.UpdatesChannel
	if os.Getenv("APP_ENV") == "production" {
		_, err := api.SetWebhook(tgbotapi.NewWebhook(
			"https://" + webhookAddress + "/no",
		))
		if err != nil {
			log.Fatal(err)
		}

		updates = api.ListenForWebhook("/no")

		go http.ListenAndServe("0.0.0.0:80", nil)

		return updates
	}

	_, err := api.RemoveWebhook()
	if err != nil {
		log.Fatalf("can't remove webhook: %v", err)
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, _ = api.GetUpdatesChan(u)
	return updates
}
