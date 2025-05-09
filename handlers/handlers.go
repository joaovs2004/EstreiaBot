package handlers

import (
	"context"
	"estreiaBot/database"
	"estreiaBot/tmdb"
	"estreiaBot/utils"
	"fmt"
	"strconv"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

var waitingForSearch = make(map[int64]bool)

func ListHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	chatId := update.Message.Chat.ID

	subscriptions := database.GetClientSubscriptions(chatId)

	if len(subscriptions) == 0 {
		utils.BotSendMessage("Você não possui nenhuma série cadastrada", chatId, ctx, b)
		return
	}

	utils.BotSendMessage("Suas séries cadastradas:", chatId, ctx, b)

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

	utils.BotSendMessage("Digite o nome da serie que voce quer buscar", chatId, ctx, b)
}

func SearchState(ctx context.Context, b *bot.Bot, update *models.Update) {
	data := tmdb.SearchShow(update.Message.Text)
	chatId := update.Message.Chat.ID

	if len(data.Results) == 0 {
		utils.BotSendMessage("Nenhuma serie encontrada com esse nome", chatId, ctx, b)
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
			Text:   "Utilize o comando /buscar para buscar uma série ou /listar para listar suas séries cadastradas",
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
	database.CreateUser(telegramID)
	database.CreateShow(showId, showName)
	database.CreateClientSubscription(telegramID, showId)

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

	database.RemoveClientSubscription(telegramID, showId)

	message := fmt.Sprintf("Removendo a série %s da sua lista", showName)

	// Send a message confirming the button selection
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: telegramID,
		Text:   message,
	})
}

func CheckForNewSeasons(ctx context.Context, b *bot.Bot) {
	subscriptions := database.GetAllSubscriptions()

	for _, subscription := range subscriptions {
		latestSeason := tmdb.GetLastSeason(subscription.ShowID)
		show := database.GetShow(subscription.ShowID)

		if latestSeason > show.LastSeason {
			// Update the last season in the database
			database.UpdateLastSeason(subscription.ShowID, latestSeason)

			// Notify the user about the new season
			message := fmt.Sprintf("A nova temporada (%d) da série %s foi lançada!", latestSeason, show.Name)
			utils.BotSendMessage(message, subscription.ClientID, ctx, b)
		}
	}
}
