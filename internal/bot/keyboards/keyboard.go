package keyboards

import (
	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
	"github.com/tamper000/freybot/internal/config"
)

func init() {
	var groupButtons []telego.InlineKeyboardButton
	for _, group := range config.ModelGroupOrder {
		stringGroup := string(group)
		groupButtons = append(groupButtons, tu.InlineKeyboardButton(stringGroup).WithCallbackData("g_"+stringGroup))
	}

	GroupKeyboard = generateTelegramKeyboard(groupButtons, 2)

	var imageButtons []telego.InlineKeyboardButton
	for _, model := range config.PhotoModels {
		imageButtons = append(imageButtons, tu.InlineKeyboardButton(model.Title).WithCallbackData("i_"+model.ApiName))
	}

	ImageKeyboard = generateTelegramKeyboard(imageButtons, 2)

	var editButtons []telego.InlineKeyboardButton
	for _, model := range config.EditModels {
		editButtons = append(editButtons, tu.InlineKeyboardButton(model.Title).WithCallbackData("e_"+model.ApiName))
	}

	EditKeyboard = generateTelegramKeyboard(editButtons, 2)
}

var GroupKeyboard *telego.InlineKeyboardMarkup
var ImageKeyboard *telego.InlineKeyboardMarkup
var EditKeyboard *telego.InlineKeyboardMarkup

var MainKeyboard = tu.Keyboard(
	tu.KeyboardRow(
		tu.KeyboardButton("Текстовые модели"),
	),
	tu.KeyboardRow(
		tu.KeyboardButton("Фото модели"),
	),
	tu.KeyboardRow(
		tu.KeyboardButton("Редактирование фото"),
	),
	tu.KeyboardRow(
		tu.KeyboardButton("Роль"),
	),
).WithResizeKeyboard().WithInputFieldPlaceholder("Выбери модель или напиши вопрос.")

var RolesKeyboard = tu.InlineKeyboard(
	tu.InlineKeyboardRow(
		tu.InlineKeyboardButton("Обычный").WithCallbackData("r_default"),
		tu.InlineKeyboardButton("Няшка").WithCallbackData("r_nyasha"),
		tu.InlineKeyboardButton("Умный").WithCallbackData("r_smart"),
	),
	tu.InlineKeyboardRow(
		tu.InlineKeyboardButton("Злой матершиник").WithCallbackData("r_evil"),
	),
)

func GenerateModelsKeyboard(info []config.AIModel) *telego.InlineKeyboardMarkup {
	var modelButtons []telego.InlineKeyboardButton
	for _, modelInfo := range info {
		var suffix string
		if modelInfo.Image {
			suffix += " 🌆"
		}

		button := tu.InlineKeyboardButton(modelInfo.Title + suffix).WithCallbackData("m_" + modelInfo.CallbackData)
		modelButtons = append(modelButtons, button)
	}

	button := tu.InlineKeyboardButton("🔙Назад").WithCallbackData("g_back")
	modelButtons = append(modelButtons, button)

	return generateTelegramKeyboard(modelButtons, 2)
}

func GenerateDummyButton(text string) *telego.InlineKeyboardMarkup {
	button := tu.InlineKeyboardButton(text).WithCallbackData("dummy")
	row := tu.InlineKeyboardRow(button)
	return tu.InlineKeyboard(row)
}

func generateTelegramKeyboard(buttons []telego.InlineKeyboardButton, maxButtonsPerRow int) *telego.InlineKeyboardMarkup {
	if len(buttons) == 0 {
		return &telego.InlineKeyboardMarkup{InlineKeyboard: [][]telego.InlineKeyboardButton{}}
	}

	var keyboard [][]telego.InlineKeyboardButton
	var currentRow []telego.InlineKeyboardButton

	for _, button := range buttons {
		currentRow = append(currentRow, button)

		if len(currentRow) >= maxButtonsPerRow {
			keyboard = append(keyboard, currentRow)
			currentRow = nil
		}
	}

	// Добавляем последнюю неполную строку
	if len(currentRow) > 0 {
		keyboard = append(keyboard, currentRow)
	}

	return tu.InlineKeyboard(keyboard...)
}
