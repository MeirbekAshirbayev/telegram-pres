package main

import (
	"encoding/json"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	// Load config manually to get token
	file, err := os.Open("config.json")
	if err != nil {
		log.Fatalf("Error opening config.json: %v", err)
	}
	defer file.Close()

	var cfg Config
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&cfg)
	if err != nil {
		log.Fatalf("Error decoding config.json: %v", err)
	}

	bot, err := tgbotapi.NewBotAPI(cfg.BotToken)
	if err != nil {
		log.Panic(err)
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)
	log.Println("Waiting for messages... Please post something in your channel.")

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.ChannelPost != nil {
			log.Printf("CHANNEL FOUND! Title: [%s], ID: [%d]", update.ChannelPost.Chat.Title, update.ChannelPost.Chat.ID)
		} else if update.Message != nil {
			log.Printf("MESSAGE FOUND! Chat Title: [%s], ID: [%d]", update.Message.Chat.Title, update.Message.Chat.ID)
		} else if update.MyChatMember != nil {
			log.Printf("ADDED TO CHAT! Title: [%s], ID: [%d]", update.MyChatMember.Chat.Title, update.MyChatMember.Chat.ID)
		}
	}
}

// Config struct to match main.go (simplified)
type Config struct {
	BotToken string `json:"bot_token"`
}
