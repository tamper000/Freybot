package handlers

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
	"github.com/tamper000/freybot/internal/bot/keyboards"
)

func (h *Handler) AddUserHandler(ctx *th.Context, message telego.Message) error {
	chatID := tu.ID(message.From.ID)

	splitted := strings.SplitN(message.Text, " ", 2)
	if len(splitted) < 2 {
		msg := tu.Message(chatID, "Напишите корректный запрос.").WithParseMode(telego.ModeMarkdown)
		_, err := ctx.Bot().SendMessage(ctx, msg)
		return err
	}
	userIDint, err := strconv.Atoi(splitted[1])

	if err != nil {
		msg := tu.Message(chatID, "*Введи правильный* _ID_ *пользователя!*").WithParseMode(telego.ModeMarkdown).WithReplyMarkup(keyboards.MainKeyboard)
		_, err := ctx.Bot().SendMessage(ctx, msg)
		return err
	}
	err = h.userRepo.AddUser(int64(userIDint))
	if err != nil {
		msg := tu.Message(chatID, fmt.Sprintf("*Не удалось добавить пользователя в БД* \n`%s`", err.Error())).WithParseMode(telego.ModeMarkdown).WithReplyMarkup(keyboards.MainKeyboard)
		_, err := ctx.Bot().SendMessage(ctx, msg)
		return err
	}

	msg := tu.Message(chatID, fmt.Sprintf("*Успешно добавили пользователя* `%d`", userIDint)).WithParseMode(telego.ModeMarkdown).WithReplyMarkup(keyboards.MainKeyboard)
	_, err = ctx.Bot().SendMessage(ctx, msg)
	return err
}

func (h *Handler) DelUserHandler(ctx *th.Context, message telego.Message) error {
	chatID := tu.ID(message.From.ID)

	splitted := strings.SplitN(message.Text, " ", 2)
	if len(splitted) < 2 {
		msg := tu.Message(chatID, "Напишите корректный запрос.").WithParseMode(telego.ModeMarkdown)
		_, err := ctx.Bot().SendMessage(ctx, msg)
		return err
	}
	userIDint, err := strconv.Atoi(splitted[1])

	if err != nil {
		msg := tu.Message(chatID, "*Введи правильный* _ID_ *пользователя!*").WithParseMode(telego.ModeMarkdown).WithReplyMarkup(keyboards.MainKeyboard)
		_, err := ctx.Bot().SendMessage(ctx, msg)
		return err
	}
	err = h.userRepo.DelUser(int64(userIDint))
	if err != nil {
		msg := tu.Message(chatID, fmt.Sprintf("*Не удалось удалить пользователя в БД* \n`%s`", err.Error())).WithParseMode(telego.ModeMarkdown).WithReplyMarkup(keyboards.MainKeyboard)
		_, err := ctx.Bot().SendMessage(ctx, msg)
		return err
	}

	msg := tu.Message(chatID, fmt.Sprintf("*Успешно удалили пользователя* `%d`", userIDint)).WithParseMode(telego.ModeMarkdown).WithReplyMarkup(keyboards.MainKeyboard)
	_, err = ctx.Bot().SendMessage(ctx, msg)
	return err
}
