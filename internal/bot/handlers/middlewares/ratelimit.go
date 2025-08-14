package middlewares

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/allegro/bigcache/v3"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

var (
	cache   *bigcache.BigCache
	adminID int64
	rate    int
)

func ConfigureRatelimit(ctx context.Context, admin int64, ratelimit int, ratetime time.Duration) error {
	adminID = admin
	rate = ratelimit

	cacheClient, err := bigcache.New(ctx, bigcache.DefaultConfig(ratetime))
	if err != nil {
		return err
	}

	cache = cacheClient
	return nil
}

func Ratelimit(ctx *th.Context, update telego.Update) error {
	var userID int64
	if update.Message != nil {
		userID = update.Message.Chat.ID
	}

	if update.CallbackQuery != nil || userID == adminID || (update.Message != nil && strings.HasPrefix(update.Message.Text, "/gen")) {
		return ctx.Next(update)
	}

	key := "rate:" + strconv.Itoa(int(userID))
	entry, err := cache.Get(key)
	if err != nil {
		cache.Set(key, []byte("1"))
		return ctx.Next(update)
	}

	count, _ := strconv.Atoi(string(entry))
	if count >= rate {
		chatID := tu.ID(userID)
		msg := tu.Message(chatID, "К сожалению вы достигли максимальное кол-во запросов.").WithParseMode(telego.ModeMarkdown)
		_, err := ctx.Bot().SendMessage(ctx, msg)
		return err
	}

	cache.Set(key, []byte(strconv.Itoa(count+1)))
	return ctx.Next(update)
}
