package handlers

import (
	"github.com/tamper000/freybot/internal/providers"
	"github.com/tamper000/freybot/internal/repository"
)

type Handler struct {
	ioClient   providers.Client
	pClient    providers.ClientPollinations
	opClient   providers.Client
	userRepo   repository.UserRepository
	dialogRepo repository.DialogRepository
}

var startMessage = `<b>Привет! Я — твой многостаночник с ИИ</b> 🔥
Я умею: генерировать текст, создавать изображения, распознавать голос и понимать фото (для этого подходят только поддерживаемые модели).

<b>Как пользоваться:</b>
- <b>Спрашивай что угодно:</b> <i>от учебных задач до развлечений. </i>
- <b>Выбирай модель под задачу:</b> <i>GPT, Qwen, DeepSeek, Gemini (для фоторазбора также доступны Qwen, Gemini, GPT, Mistral — отмечены 🌆).</i>
- <b>Включай роли:</b> <i>от узкого специалиста до «няшечки» или дерзкого хулигана.</i>

<b>Рекомендации:</b>
- <b>Текст:</b> <i>Qwen v3, DeepSeek R1 и его подверсии, Gemini, GPT</i> — отличные варианты.
- <b>Фото:</b> <i>SDXL TURBO или GPT</i> — быстрый и качественный результат.

<b>Как сгенерировать картинку:</b>
- <b>Используй команду:</b> <code>/gen человек паук на фоне нью йорка</code>

<i><b>Готов? Пиши запрос — я подберу лучшую модель и выдам результат 🚀</b></i>`
