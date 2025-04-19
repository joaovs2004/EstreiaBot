package utils

import (
	"log"
	"os"
	"context"

	"github.com/joho/godotenv"
	"github.com/go-telegram/bot/models"
	"github.com/go-telegram/bot"
)

func GetDotenvValue(desiredValue string) (string) {
  err := godotenv.Load()
  if err != nil {
    log.Fatal("Error loading .env file")
  }

  value := os.Getenv(desiredValue)
  return value
}

func BotSendMessage(text string, chatId int64, ctx context.Context, b *bot.Bot) {
  b.SendMessage(ctx, &bot.SendMessageParams {
		ChatID:    chatId,
		Text:      text,
		ParseMode: models.ParseModeMarkdown,
	})
}
