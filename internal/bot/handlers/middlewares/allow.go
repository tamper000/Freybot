package middlewares

import (
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	"github.com/tamper000/freybot/internal/repository"
)

func OnlyAllowUsers(userRepo repository.UserRepository, adminID int64) func(ctx *th.Context, update telego.Update) error {
	return func(ctx *th.Context, update telego.Update) error {
		var userID int64
		if update.Message != nil {
			userID = update.Message.Chat.ID
		} else if update.CallbackQuery != nil {
			userID = update.CallbackQuery.Message.GetChat().ID
		} else if update.InlineQuery != nil {
			userID = update.InlineQuery.From.ID
		} else if update.ChosenInlineResult != nil {
			userID = update.ChosenInlineResult.From.ID
		}

		if userID == adminID {
			return ctx.Next(update)
		}

		_, err := userRepo.GetUser(userID)
		if err != nil {
			return err
		}

		return ctx.Next(update)
	}
}
