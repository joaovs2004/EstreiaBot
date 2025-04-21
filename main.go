package main

import (
	"context"
	"os"
	"os/signal"

	// "estreiaBot/database"
	"github.com/go-telegram/bot"
)

func main() {
	InitDb()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	opts := []bot.Option{
		bot.WithDefaultHandler(DefaultHandler),
		bot.WithCallbackQueryDataHandler("addshow", bot.MatchTypePrefix, AddCallbackHandler),
		bot.WithCallbackQueryDataHandler("removeshow", bot.MatchTypePrefix, RemoveCallbackHandler),
	}

	telegramToken := GetDotenvValue("TELEGRAM_TOKEN")

	b, err := bot.New(telegramToken, opts...)
	if err != nil {
		panic(err)
	}

	b.RegisterHandler(bot.HandlerTypeMessageText, "listar", bot.MatchTypeCommand, ListHandler)
	b.RegisterHandler(bot.HandlerTypeMessageText, "buscar", bot.MatchTypeCommand, SearchHandler)

	b.Start(ctx)
}
