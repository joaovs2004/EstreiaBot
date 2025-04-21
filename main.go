package main

import (
	"context"
	"os"
	"os/signal"
	"time"

	"estreiaBot/database"
	"estreiaBot/handlers"
	"estreiaBot/utils"

	"github.com/go-telegram/bot"
)

func main() {
	database.InitDb()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	opts := []bot.Option{
		bot.WithDefaultHandler(handlers.DefaultHandler),
		bot.WithCallbackQueryDataHandler("addshow", bot.MatchTypePrefix, handlers.AddCallbackHandler),
		bot.WithCallbackQueryDataHandler("removeshow", bot.MatchTypePrefix, handlers.RemoveCallbackHandler),
	}

	telegramToken := utils.GetDotenvValue("TELEGRAM_TOKEN")

	b, err := bot.New(telegramToken, opts...)
	if err != nil {
		panic(err)
	}

	b.RegisterHandler(bot.HandlerTypeMessageText, "listar", bot.MatchTypeCommand, handlers.ListHandler)
	b.RegisterHandler(bot.HandlerTypeMessageText, "buscar", bot.MatchTypeCommand, handlers.SearchHandler)

	// Periodically check for new seasons
	go func() {
		ticker := time.NewTicker(24 * time.Hour) // Check once a day
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				handlers.CheckForNewSeasons(ctx, b)
			case <-ctx.Done():
				return
			}
		}
	}()

	b.Start(ctx)
}
