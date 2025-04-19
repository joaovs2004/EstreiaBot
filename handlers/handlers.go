package handlers

import (
	"context"
	"strconv"

	"estreiaBot/api"
  "estreiaBot/utils"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

var waitingForSearch = make(map[int64]bool)

func ListHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
  chatId := update.Message.Chat.ID
  waitingForSearch[chatId] = true

  utils.BotSendMessage("Comando list", chatId, ctx, b)
}

func SearchHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
  chatId := update.Message.Chat.ID
  waitingForSearch[chatId] = true

  utils.BotSendMessage("Digite o nome da serie que voce quer buscar", chatId, ctx, b)
}

func SearchState(ctx context.Context, b *bot.Bot, update *models.Update) {
    data := api.SearchShow(update.Message.Text)
    chatId := update.Message.Chat.ID
    
    if len(data.Results) == 0 {
      utils.BotSendMessage("Nenhuma serie encontrada com esse nome", chatId, ctx, b)
    } else if len(data.Results) > 0 {
      for _, show := range data.Results {
        kb := &models.InlineKeyboardMarkup{
          InlineKeyboard: [][]models.InlineKeyboardButton{
            {
              {Text: "Adicionar na lista", CallbackData: "show_" + strconv.Itoa(show.Id)},
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

func CallbackHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
  b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
		ShowAlert:       false,
	})

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.CallbackQuery.Message.Message.Chat.ID,
		Text:   "You selected the button: " + update.CallbackQuery.Data,
	})
}
