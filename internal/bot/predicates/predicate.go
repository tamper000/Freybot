package predicate

import (
	"context"

	"github.com/mymmrac/telego"
)

func OnlyAdmin(adminID int64) func(ctx context.Context, update telego.Update) bool {
	return func(ctx context.Context, update telego.Update) bool {
		if update.Message != nil && update.Message.Chat.ID == adminID {
			return true
		} else if update.CallbackQuery != nil && update.CallbackQuery.Message.GetChat().ID == adminID {
			return true
		}

		return false
	}
}

func OnlyPrivate(ctx context.Context, update telego.Update) bool {
	if update.Message != nil && update.Message.Chat.Type == "private" {
		return true
	} else if update.CallbackQuery != nil && update.CallbackQuery.Message.GetChat().Type == "private" {
		return true
	}

	return false
}

func OnlyVoice(ctx context.Context, update telego.Update) bool {
	if update.Message != nil && update.Message.Voice != nil {
		return true
	}

	return false
}

func OnlyPhoto(ctx context.Context, update telego.Update) bool {
	if update.Message != nil && len(update.Message.Photo) > 0 {
		return true
	}

	return false
}
