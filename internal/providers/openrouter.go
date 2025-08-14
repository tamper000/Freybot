package providers

import (
	"time"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

const openrouterBaseURL = "https://openrouter.ai/api/v1/"

func NewOpenRouterClient(token string, timeout time.Duration) Client {
	client := openai.NewClient(option.WithAPIKey(token), option.WithBaseURL(openrouterBaseURL))

	return &OpenaiClient{
		client:  &client,
		timeout: timeout,
	}
}
