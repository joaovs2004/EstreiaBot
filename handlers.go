package main

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

var waitingForSearch = make(map[int64]bool)

func ListHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	chatId := update.Message.Chat.ID
	waitingForSearch[chatId] = true

	subscriptions := GetClientSubscriptions(chatId)

	if len(subscriptions) == 0 {
		BotSendMessage("Você não possui nenhuma série cadastrada", chatId, ctx, b)
		return
	}

	BotSendMessage("Suas séries cadastradas:", chatId, ctx, b)

	for _, subscription := range subscriptions {
		callBackDataShow := fmt.Sprintf("removeshow_%s_%s", subscription.ShowID, subscription.Name)
		kb := &models.InlineKeyboardMarkup{
			InlineKeyboard: [][]models.InlineKeyboardButton{
				{
					{Text: "Remover da lista", CallbackData: callBackDataShow},
				},
			},
		}
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:      update.Message.Chat.ID,
			Text:        subscription.Name,
			ReplyMarkup: kb,
		})
	}
}

func SearchHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	chatId := update.Message.Chat.ID
	waitingForSearch[chatId] = true

	BotSendMessage("Digite o nome da serie que voce quer buscar", chatId, ctx, b)
}

func SearchState(ctx context.Context, b *bot.Bot, update *models.Update) {
	data := SearchShow(update.Message.Text)
	chatId := update.Message.Chat.ID

	if len(data.Results) == 0 {
		BotSendMessage("Nenhuma serie encontrada com esse nome", chatId, ctx, b)
	} else if len(data.Results) > 0 {
		for _, show := range data.Results {
			callBackDataShow := fmt.Sprintf("addshow_%s_%s", strconv.Itoa(show.Id), show.OriginalName)
			kb := &models.InlineKeyboardMarkup{
				InlineKeyboard: [][]models.InlineKeyboardButton{
					{
						{Text: "Adicionar na lista", CallbackData: callBackDataShow},
					},
				},
			}
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID:      update.Message.Chat.ID,
				Text:        show.OriginalName,
				ReplyMarkup: kb,
			})
		}
	}
}

func DefaultHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	chatId := update.Message.Chat.ID

	if waitingForSearch[chatId] {
		SearchState(ctx, b, update)
		waitingForSearch[chatId] = false
	} else {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   update.Message.Text,
		})
	}
}

func AddCallbackHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
		ShowAlert:       false,
	})

	showInfo := strings.Split(update.CallbackQuery.Data, "_")
	showId := showInfo[1]
	showName := showInfo[2]

	telegramID := update.CallbackQuery.Message.Message.Chat.ID
	CreateUser(telegramID)
	CreateShow(showId, showName)
	CreateClientSubscription(telegramID, showId)

	message := fmt.Sprintf("Adicionando a série %s na sua lista", showName)

	// Send a message confirming the button selection
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: telegramID,
		Text:   message,
	})
}

func RemoveCallbackHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
		ShowAlert:       false,
	})

	showInfo := strings.Split(update.CallbackQuery.Data, "_")
	showId := showInfo[1]
	showName := showInfo[2]

	telegramID := update.CallbackQuery.Message.Message.Chat.ID

	RemoveClientSubscription(telegramID, showId)

	message := fmt.Sprintf("Removendo a série %s da sua lista", showName)

	// Send a message confirming the button selection
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: telegramID,
		Text:   message,
	})
}
