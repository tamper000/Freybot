package handlers

import (
	"fmt"
	"strings"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
	"github.com/tamper000/freybot/internal/bot/keyboards"
)

func (h *Handler) ChooseGroup(ctx *th.Context, message telego.Message) error {
	chatID := tu.ID(message.From.ID)
	msg := tu.Message(chatID, "_Выбери группу модели:_").WithParseMode(telego.ModeMarkdown).WithReplyMarkup(keyboards.GroupKeyboard)
	_, err := ctx.Bot().SendMessage(ctx, msg)
	return err

}

func (h *Handler) ChooseModel(ctx *th.Context, query telego.CallbackQuery) error {
	ctx.Bot().AnswerCallbackQuery(ctx.Context(), tu.CallbackQuery(query.ID))
	group := strings.Replace(query.Data, "g_", "", 1)

	if group == "back" {
		chatID := tu.ID(query.From.ID)
		msg := tu.EditMessageText(chatID, query.Message.GetMessageID(), "_Выбери группу модели:_").WithParseMode(telego.ModeMarkdown).WithReplyMarkup(keyboards.GroupKeyboard)
		_, err := ctx.Bot().EditMessageText(ctx, msg)
		return err
	}

	h.userRepo.UpdateGroup(query.From.ID, group)
	info := GetModelsByGroup(group)

	msg := tu.EditMessageText(tu.ID(query.From.ID), query.Message.GetMessageID(), "_Выбери одну из моделей_").WithParseMode(telego.ModeMarkdown).WithReplyMarkup(keyboards.GenerateModelsKeyboard(info))
	ctx.Bot().EditMessageText(ctx.Context(), msg)
	return nil
}

func (h *Handler) ChooseEnd(ctx *th.Context, query telego.CallbackQuery) error {
	ctx.Bot().AnswerCallbackQuery(ctx.Context(), tu.CallbackQuery(query.ID))
	CallbackData := strings.Replace(query.Data, "m_", "", 1)

	user, _ := h.userRepo.GetUser(query.From.ID)

	info := GetModelByCallback(user.Group, CallbackData)

	h.userRepo.UpdateProvider(query.From.ID, string(info.Provider))
	h.userRepo.UpdateTextModel(query.From.ID, info.ApiName)

	msg := tu.EditMessageText(tu.ID(query.From.ID), query.Message.GetMessageID(), fmt.Sprintf("*Выбор завершен!*\n\n*Модель*: _%s_", info.Title)).WithParseMode(telego.ModeMarkdown)
	ctx.Bot().EditMessageText(ctx.Context(), msg)

	return nil
}

func (h *Handler) ChoosePhoto(ctx *th.Context, message telego.Message) error {
	chatID := tu.ID(message.From.ID)
	msg := tu.Message(chatID, "Выбери одну из моделей").WithParseMode(telego.ModeMarkdown).WithReplyMarkup(keyboards.ImageKeyboard)
	_, err := ctx.Bot().SendMessage(ctx, msg)
	return err
}

func (h *Handler) ChoosePhotoModel(ctx *th.Context, query telego.CallbackQuery) error {
	ctx.Bot().AnswerCallbackQuery(ctx.Context(), tu.CallbackQuery(query.ID))
	model := strings.Replace(query.Data, "i_", "", 1)
	h.userRepo.UpdatePhotoModel(query.From.ID, model)

	info := GetPhotoModelByApiName(model)

	msg := tu.EditMessageText(tu.ID(query.From.ID), query.Message.GetMessageID(), fmt.Sprintf("*Выбор завершен!*\n\n*Модель*: _%s_", info.Title)).WithParseMode(telego.ModeMarkdown)
	ctx.Bot().EditMessageText(ctx.Context(), msg)
	return nil
}

func (h *Handler) ChooseRole(ctx *th.Context, message telego.Message) error {
	chatID := tu.ID(message.From.ID)
	msg := tu.Message(chatID, "Выбери одну из ролей.").WithParseMode(telego.ModeMarkdown).WithReplyMarkup(keyboards.RolesKeyboard)
	_, err := ctx.Bot().SendMessage(ctx, msg)
	return err
}

func (h *Handler) ChooseRoleCallback(ctx *th.Context, query telego.CallbackQuery) error {
	ctx.Bot().AnswerCallbackQuery(ctx.Context(), tu.CallbackQuery(query.ID))
	role := strings.Replace(query.Data, "r_", "", 1)
	h.userRepo.UpdateRole(query.From.ID, role)

	msg := tu.EditMessageText(tu.ID(query.From.ID), query.Message.GetMessageID(), "*Выбор завершен!*").WithParseMode(telego.ModeMarkdown)
	ctx.Bot().EditMessageText(ctx.Context(), msg)
	return nil
}
