package providers

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/openai/openai-go"
	"github.com/tamper000/freybot/internal/config"
	"github.com/tamper000/freybot/internal/models"
)

type Client interface {
	NewMessage(history []models.Message, model, role string) (string, error)
	NewMessageWithPhoto(message, model string, photo []byte) (string, error)
}

type OpenaiClient struct {
	client  *openai.Client
	timeout time.Duration
}

func (c *OpenaiClient) NewMessage(history []models.Message, model, role string) (string, error) {
	prompt := GetRole(role)
	openaiHistory := GenerateHistory(prompt, history)

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	resp, err := c.client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model:    openai.ChatModel(model),
		Messages: openaiHistory,
	})

	fmt.Println(resp, err, model)

	if err != nil {
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}

func (c *OpenaiClient) NewMessageWithPhoto(message, model string, photo []byte) (string, error) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	photoBase64 := base64.StdEncoding.EncodeToString([]byte(photo))
	imageURL := fmt.Sprintf("data:image/jpeg;base64,%s", photoBase64)

	resp, err := c.client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model: openai.ChatModel(model),
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(message),
			openai.UserMessage([]openai.ChatCompletionContentPartUnionParam{
				openai.ImageContentPart(openai.ChatCompletionContentPartImageImageURLParam{
					URL: imageURL,
				}),
			}),
		},
	})

	if err != nil {
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}

func CreateClients(cfg *config.Config) (Client, ClientPollinations, Client, Client) {
	return NewIoNewClient(cfg.Models.IoNetToken, cfg.Models.Timeout),
		NewPollinationsClient(cfg.Models.PollinationsToken, cfg.Models.Timeout),
		NewOpenRouterClient(cfg.Models.OpenRouterToken, cfg.Models.Timeout),
		NewLLM7Client(cfg.Models.LLM7Token, cfg.Models.Timeout)
}

func GetRole(role string) string {
	switch role {
	case "default":
		return defaultPrompt
	case "nyasha":
		return nyashaPrompt
	case "smart":
		return smartPrompt
	case "evil":
		return evilPrompt
	}

	return defaultPrompt
}

func GenerateHistory(prompt string, history []models.Message) (data []openai.ChatCompletionMessageParamUnion) {
	data = append(data, openai.SystemMessage(prompt+textFormatPrompt))

	for _, msg := range history {
		switch msg.Role {
		case "assistant":
			data = append(data, openai.AssistantMessage(msg.Content))
		case "user":
			data = append(data, openai.UserMessage(msg.Content))
		}
	}

	return
}
