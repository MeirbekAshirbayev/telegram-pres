package main

import (
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var bot *tgbotapi.BotAPI

func initBot(token string) {
	var err error
	bot, err = tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatalf("Error initializing bot: %v", err)
	}
	log.Printf("Authorized on account %s", bot.Self.UserName)
}

func isUserMember(userID int64, channelID int64) (bool, error) {
	chatConfig := tgbotapi.GetChatMemberConfig{
		ChatConfigWithUser: tgbotapi.ChatConfigWithUser{
			ChatID: channelID,
			UserID: userID,
		},
	}
	member, err := bot.GetChatMember(chatConfig)
	if err != nil {
		return false, err
	}

	// Status can be creator, administrator, member, restricted, left, or kicked
	if member.Status == "member" || member.Status == "administrator" || member.Status == "creator" {
		return true, nil
	}
	return false, nil
}

func PostPresentationToChannel(channelID int64, topicID int64, title, presentationURL string) error {
	// Create a formatted message: [Title](URL)
	// Using HTML parse mode for cleaner links: <a href="URL">Title</a>
	msgText := fmt.Sprintf("📚 <b>%s</b>\n\n👇 Көру үшін басыңыз:\n<a href=\"%s\">%s</a>", title, presentationURL, title)

	msg := tgbotapi.NewMessage(channelID, msgText)
	msg.ParseMode = "HTML"
	if topicID != 0 {
		msg.ReplyToMessageID = int(topicID)
	}

	_, err := bot.Send(msg)
	return err
}
