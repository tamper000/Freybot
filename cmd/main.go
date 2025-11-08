package main

import (
	"context"
	"fmt"
	"log"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	"github.com/tamper000/freybot/internal/bot/handlers"
	"github.com/tamper000/freybot/internal/bot/handlers/middlewares"
	predicate "github.com/tamper000/freybot/internal/bot/predicates"
	"github.com/tamper000/freybot/internal/config"
	"github.com/tamper000/freybot/internal/database"
	"github.com/tamper000/freybot/internal/metrics"
	"github.com/tamper000/freybot/internal/providers"
	"github.com/tamper000/freybot/internal/repository"
	"github.com/valyala/fasthttp"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	if err := cfg.Validate(); err != nil {
		log.Fatal(err)
	}

	db, err := database.LoadDatabase(cfg.Database.Path)
	if err != nil {
		log.Fatal(err)
	}
	userRepo := repository.NewUserRepository(db)
	userRepo.AddUser(cfg.Telegram.AdminID)
	dialogRepo := repository.NewDialogRepository(db, cfg.Database.MaxHistory)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	bot, err := telego.NewBot(cfg.Telegram.Token, telego.WithDefaultLogger(false, true))
	if err != nil {
		log.Fatal(err)
	}
	bot.DeleteWebhook(ctx, nil)

	if err := SetCommands(ctx, bot); err != nil {
		log.Fatal(err)
	}

	var bh *th.BotHandler
	var srv *fasthttp.Server

	if !cfg.Webhook.Enabled {
		bh = LongPoll(ctx, bot)
	} else {
		bh, srv = WebHook(ctx, bot, cfg.Webhook.Domain, strconv.Itoa(cfg.Webhook.Port))
	}

	flux, err := providers.NewFluxClient(cfg.Proxy)
	if err != nil {
		log.Fatal(err)
	}

	ionet, pollinations, openrouter, llm7 := providers.CreateClients(cfg)
	handlers := handlers.NewHandlers(ionet, pollinations, openrouter, llm7,
		userRepo, dialogRepo, flux)

	if err = middlewares.ConfigureRatelimit(ctx, cfg.Telegram.AdminID, cfg.Ratelimit.Rate, cfg.Ratelimit.Time); err != nil {
		log.Fatal(err)
	}

	bh.Use(th.PanicRecovery())
	bh.Use(middlewares.OnlyAllowUsers(userRepo, cfg.Telegram.AdminID))
	bh.Use(middlewares.Ratelimit)

	AddAdminHandlers(bh, handlers, cfg.Telegram.AdminID)
	AddPrivateHandlers(bh, handlers)

	if cfg.Prometheus.Enabled {
		go metrics.StartServer(ctx, ":"+strconv.Itoa((cfg.Prometheus.Port)))
	}

	go func() {
		fmt.Println("Starting bot handler...")
		if err := bh.Start(); err != nil {
			log.Printf("Bot handler error: %v", err)
		}
	}()

	// Graceful shutdown
	done := make(chan struct{}, 1)
	go func() {
		<-ctx.Done()
		fmt.Println("Stopping...")

		stopCtx, stopCancel := context.WithTimeout(context.Background(), time.Second*10)
		defer stopCancel()

		if srv != nil {
			_ = srv.ShutdownWithContext(stopCtx)
			fmt.Println("Server done")
		}

		_ = bh.StopWithContext(stopCtx)
		fmt.Println("Bot handler done")

		sqlDB, err := db.DB()
		if err != nil {
			sqlDB.Close()
			fmt.Println("Database connection done")
		}

		done <- struct{}{}
	}()

	<-done
	fmt.Println("Done")

}

func AddPrivateHandlers(bh *th.BotHandler, handlers *handlers.Handler) {
	private := bh.Group(predicate.OnlyPrivate)

	// Start handler
	private.HandleMessage(handlers.StartHandler, th.CommandEqual("start"))

	// Choose models
	private.HandleMessage(handlers.ChooseGroup, th.TextEqual("Текстовые модели"))
	private.HandleMessage(handlers.ChooseGroup, th.CommandEqual("text"))
	private.HandleCallbackQuery(handlers.ChooseModel, th.CallbackDataPrefix("g_"))
	private.HandleCallbackQuery(handlers.ChooseEnd, th.CallbackDataPrefix("m_"))

	private.HandleMessage(handlers.ChoosePhoto, th.TextEqual("Фото модели"))
	private.HandleMessage(handlers.ChoosePhoto, th.CommandEqual("photo"))
	private.HandleCallbackQuery(handlers.ChoosePhotoModel, th.CallbackDataPrefix("i_"))
	private.HandleMessage(handlers.GenPhoto, th.CommandEqual("gen"))

	private.HandleMessage(handlers.ChooseRole, th.TextEqual("Роль"))
	private.HandleMessage(handlers.ChooseRole, th.CommandEqual("role"))
	private.HandleCallbackQuery(handlers.ChooseRoleCallback, th.CallbackDataPrefix("r_"))

	private.HandleMessage(handlers.ChooseEditModel, th.TextEqual("Редактирование фото"))
	private.HandleMessage(handlers.ChooseEditModel, th.CommandEqual("edit"))
	private.HandleCallbackQuery(handlers.ChooseEditCallback, th.CallbackDataPrefix("e_"))

	private.HandleMessage(handlers.EditPhoto, predicate.OnlyPhotoEdit)

	// Clear history
	private.HandleMessage(handlers.ClearHandler, th.CommandEqual("clear"))

	// Image messages
	private.HandleMessage(handlers.ImageHandler, predicate.OnlyPhoto)

	// Handle dummy button
	private.HandleCallbackQuery(handlers.DummyButton, th.CallbackDataEqual("dummy"))

	// Generate message
	private.HandleMessage(handlers.MessageHandler, th.AnyMessageWithText())
	private.HandleMessage(handlers.MessageHandler, predicate.OnlyVoice)
}

func AddAdminHandlers(bh *th.BotHandler, handlers *handlers.Handler, adminID int64) {
	admin := bh.Group(predicate.OnlyAdmin(adminID))

	// Add user
	admin.HandleMessage(handlers.AddUserHandler, th.CommandEqual("add"))

	// Del user
	admin.HandleMessage(handlers.DelUserHandler, th.CommandEqual("del"))
}

func LongPoll(ctx context.Context, bot *telego.Bot) *th.BotHandler {
	updates, err := bot.UpdatesViaLongPolling(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	bh, err := th.NewBotHandler(bot, updates)
	if err != nil {
		log.Fatal(err)
	}

	return bh
}

func WebHook(ctx context.Context, bot *telego.Bot, domain, port string) (*th.BotHandler, *fasthttp.Server) {
	srv := &fasthttp.Server{}

	// Start server FIRST
	go func() {
		fmt.Println("Starting server on :" + port)
		if err := srv.ListenAndServe(":" + port); err != nil {
			log.Printf("Server error: %v", err)
		}
	}()

	webhook := telego.WebhookFastHTTP(srv, "/bot", bot.SecretToken())
	updates, err := bot.UpdatesViaWebhook(ctx, webhook, telego.WithWebhookBuffer(128), telego.WithWebhookSet(ctx, &telego.SetWebhookParams{
		URL:         "https://" + domain + "/bot",
		SecretToken: bot.SecretToken(),
	}))
	if err != nil {
		log.Fatal(err)
	}

	bh, err := th.NewBotHandler(bot, updates)
	if err != nil {
		log.Fatal(err)
	}

	return bh, srv
}

func SetCommands(ctx context.Context, bot *telego.Bot) error {
	return bot.SetMyCommands(ctx, &telego.SetMyCommandsParams{
		Commands: []telego.BotCommand{
			telego.BotCommand{Command: "start", Description: "Приветствие и помощь"},
			telego.BotCommand{Command: "text", Description: "Текстовые модели"},
			telego.BotCommand{Command: "photo", Description: "Генерация фото"},
			telego.BotCommand{Command: "role", Description: "Выбор роли для ИИ"},
			telego.BotCommand{Command: "clear", Description: "Очистить историю"},
			telego.BotCommand{Command: "gen", Description: "Генерация фото"},
			telego.BotCommand{Command: "edit", Description: "Редактирование фото"},
		},
	})
}
