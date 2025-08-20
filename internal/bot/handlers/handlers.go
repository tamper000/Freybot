package handlers

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/mymmrac/telego"
	"github.com/tamper000/freybot/internal/bot/keyboards"
	"github.com/tamper000/freybot/internal/providers"
	"github.com/tamper000/freybot/internal/repository"
	md "github.com/zavitkov/tg-markdown"

	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

func NewHandlers(ioClient providers.Client, pClient providers.ClientPollinations, opClient providers.Client,
	userRepo repository.UserRepository, dialogRepo repository.DialogRepository) *Handler {

	return &Handler{
		ioClient:   ioClient,
		pClient:    pClient,
		opClient:   opClient,
		userRepo:   userRepo,
		dialogRepo: dialogRepo,
	}
}

func (h *Handler) StartHandler(ctx *th.Context, message telego.Message) error {
	chatID := tu.ID(message.From.ID)
	msg := tu.Message(chatID, startMessage).WithParseMode(telego.ModeHTML).WithReplyMarkup(keyboards.MainKeyboard)
	_, err := ctx.Bot().SendMessage(ctx, msg)
	return err
}

func (h *Handler) MessageHandler(ctx *th.Context, message telego.Message) error {
	chatID := tu.ID(message.From.ID)

	ctx.Bot().SendChatAction(ctx, tu.ChatAction(chatID, telego.ChatActionTyping))

	user, err := h.userRepo.GetUser(message.From.ID)
	if err != nil {
		msg := tu.Message(chatID, "Произошла неизвестная ошибка.")
		ctx.Bot().SendMessage(ctx, msg)
		return err
	}

	if user.Model == "" {
		msg := tu.Message(chatID, "Для начала выберите одну из моделей!")
		_, err := ctx.Bot().SendMessage(ctx, msg)
		return err
	}
	info := GetModelByUser(user)

	var client providers.Client
	provider := GetProviderByUser(user)
	client = h.GetClientByProvider(provider)

	var answer string
	if message.Text != "" {
		answer = message.Text
	} else if message.Voice != nil {
		text, err := h.TranscribeAudio(ctx, ctx.Bot(), message.Voice)
		if err != nil {
			msg := tu.Message(chatID, "К сожалению не удалось обработать ваше голосовое сообщение.")
			ctx.Bot().SendMessage(ctx, msg)
			return err
		}

		answer = text
	}

	if err := h.dialogRepo.AddMessage(message.From.ID, "user", answer); err != nil {
		return err
	}

	history, err := h.dialogRepo.GetHistory(message.From.ID)
	if err != nil {
		h.dialogRepo.DeleteLastMessage(message.From.ID)
		return err
	}

	msg := tu.Message(chatID, answer).
		WithReplyMarkup(keyboards.GenerateDummyButton("Генерируется " + info.Title))
	sended, err := ctx.Bot().SendMessage(ctx, msg)
	if err != nil {
		h.dialogRepo.DeleteLastMessage(message.From.ID)
		return err
	}

	resp, err := client.NewMessage(history, user.Model, user.Role)
	if err != nil {
		h.dialogRepo.DeleteLastMessage(message.From.ID)
		messageEdit := tu.EditMessageText(chatID, sended.MessageID, "Не удалось получить ответ...")
		ctx.Bot().EditMessageText(ctx, messageEdit)
		return err
	}

	if strings.Contains(resp, "<think>") {
		respSplitted := strings.Split(resp, "</think>")
		resp = respSplitted[1]
	}

	h.dialogRepo.AddMessage(message.From.ID, "assistant", resp)
	ctx.Bot().DeleteMessage(ctx, tu.Delete(chatID, sended.MessageID))

	resp = strings.ReplaceAll(resp, "<br>", "\n")
	resp = strings.ReplaceAll(resp, `\n`, "\n")

	result, err := SplitHTML(resp)
	if err != nil {
		fmt.Println(resp)
		msg := tu.Message(chatID, err.Error()).
			WithReplyMarkup(keyboards.GenerateDummyButton("Сгенерировано " + info.Title))
		ctx.Bot().SendMessage(ctx, msg)
		return err
	}

	for i := range len(result) {
		text := result[i]
		msg := tu.Message(chatID, text).WithParseMode(telego.ModeHTML)
		if i == len(result)-1 {
			msg = msg.WithReplyMarkup(keyboards.GenerateDummyButton("Сгенерировано " + info.Title))
		}
		ctx.Bot().SendMessage(ctx, msg)
	}
	return nil
}

func (h *Handler) GenPhoto(ctx *th.Context, message telego.Message) error {
	chatID := tu.ID(message.From.ID)

	ctx.Bot().SendChatAction(ctx, tu.ChatAction(chatID, telego.ChatActionTyping))

	splitted := strings.SplitN(message.Text, " ", 2)
	if len(splitted) < 2 {
		msg := tu.Message(chatID, "Напишите корректный запрос.").WithParseMode(telego.ModeMarkdown)
		_, err := ctx.Bot().SendMessage(ctx, msg)
		return err
	}

	prompt := splitted[1]

	user, err := h.userRepo.GetUser(message.From.ID)
	if err != nil {
		msg := tu.Message(chatID, "Произошла неизвестная ошибка.")
		ctx.Bot().SendMessage(ctx, msg)
		return err
	}
	if user.Photo == "" {
		msg := tu.Message(chatID, "Для начала выберите одну из моделей!").WithParseMode(telego.ModeMarkdown)
		_, err := ctx.Bot().SendMessage(ctx, msg)
		return err
	}
	info := GetPhotoModelByApiName(user.Photo)

	msg := tu.Message(chatID, "_Генерируем изображение..._").WithParseMode(telego.ModeMarkdown).WithParseMode(telego.ModeMarkdown).
		WithReplyMarkup(keyboards.GenerateDummyButton("Генерируется " + info.Title))
	sended, err := ctx.Bot().SendMessage(ctx, msg)
	if err != nil {
		return err
	}

	start := time.Now()
	photoBytes, err := h.pClient.GeneratePhoto(prompt, user.Photo)
	if err != nil {
		msg := tu.EditMessageText(chatID, sended.MessageID, "К сожалению не удалось сгенерировать фото.")
		_, err := ctx.Bot().EditMessageText(ctx, msg)
		return err
	}

	ctx.Bot().DeleteMessage(ctx, tu.Delete(chatID, sended.MessageID))
	elapsed := time.Since(start)

	photo := tu.Photo(chatID, tu.FileFromBytes(photoBytes, "ai_image.jpg")).WithCaption(fmt.Sprintf("На генерацию было затрачено _%.2f сек_.", elapsed.Seconds())).WithParseMode(telego.ModeMarkdown).
		WithReplyMarkup(keyboards.GenerateDummyButton("Сгенерировано " + info.Title))
	_, err = ctx.Bot().SendPhoto(ctx, photo)
	return err
}

func (h *Handler) ImageHandler(ctx *th.Context, message telego.Message) error {
	chatID := tu.ID(message.From.ID)

	ctx.Bot().SendChatAction(ctx, tu.ChatAction(chatID, telego.ChatActionTyping))

	user, err := h.userRepo.GetUser(message.From.ID)
	if err != nil {
		msg := tu.Message(chatID, "Произошла неизвестная ошибка.")
		ctx.Bot().SendMessage(ctx, msg)
		return err
	}
	if user.Model == "" {
		msg := tu.Message(chatID, "Для начала выберите одну из моделей!")
		_, err := ctx.Bot().SendMessage(ctx, msg)
		return err
	}
	info := GetModelByUser(user)

	if !IsSupportPhotoByUser(user) {
		msg := tu.Message(chatID, "Текущая модель не поддерживает изображения!").WithReplyMarkup(keyboards.MainKeyboard)
		_, err := ctx.Bot().SendMessage(ctx, msg)
		return err
	}

	var client providers.Client
	provider := GetProviderByUser(user)
	client = h.GetClientByProvider(provider)

	msg := tu.Message(chatID, "_Думаем..._").WithParseMode(telego.ModeMarkdown).
		WithReplyMarkup(keyboards.GenerateDummyButton("Генерируется " + info.Title))
	sended, _ := ctx.Bot().SendMessage(ctx, msg)

	file, err := ctx.Bot().GetFile(context.Background(), &telego.GetFileParams{
		FileID: message.Photo[len(message.Photo)-1].FileID,
	})
	if err != nil {
		msg := tu.Message(chatID, "К сожалению не удалось обработать ваше фото.")
		_, err := ctx.Bot().SendMessage(ctx, msg)
		return err
	}

	url := ctx.Bot().FileDownloadURL(file.FilePath)
	fileData, err := tu.DownloadFile(url)
	if err != nil {
		msg := tu.Message(chatID, "К сожалению не удалось обработать ваше фото.")
		_, err := ctx.Bot().SendMessage(ctx, msg)
		return err
	}

	resp, err := client.NewMessageWithPhoto(message.Caption, user.Model, fileData)
	if err != nil {
		messageEdit := tu.EditMessageText(chatID, sended.MessageID, "Не удалось получить ответ...")
		ctx.Bot().EditMessageText(ctx, messageEdit)
		return err
	}

	text := md.ConvertMarkdownToTelegramMarkdownV2(resp)
	messageEdit := tu.EditMessageText(chatID, sended.MessageID, text).
		WithReplyMarkup(keyboards.GenerateDummyButton("Сгенерировано " + info.Title)).
		WithParseMode(telego.ModeMarkdownV2)
	if _, err = ctx.Bot().EditMessageText(ctx, messageEdit); err != nil {
		messageEdit.ParseMode = ""
		ctx.Bot().EditMessageText(ctx, messageEdit)
	}
	return err
}

func (h *Handler) ClearHandler(ctx *th.Context, message telego.Message) error {
	chatID := tu.ID(message.From.ID)

	if err := h.dialogRepo.ClearHistory(message.From.ID); err != nil {
		msg := tu.Message(chatID, "Не удалось очистить историю.")
		ctx.Bot().SendMessage(ctx, msg)
		return err
	}

	msg := tu.Message(chatID, "История успешно очищена!")
	_, err := ctx.Bot().SendMessage(ctx, msg)
	return err
}

func (h *Handler) DummyButton(ctx *th.Context, query telego.CallbackQuery) error {
	ctx.Bot().AnswerCallbackQuery(ctx, tu.CallbackQuery(query.ID))
	return nil
}
