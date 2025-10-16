package providers

import (
	"time"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

const llm7BaseURL = "https://api.llm7.io/v1"

func NewLLM7Client(token string, timeout time.Duration) Client {
	client := openai.NewClient(option.WithAPIKey(token), option.WithBaseURL(llm7BaseURL))

	return &OpenaiClient{
		client:  &client,
		timeout: timeout,
	}
}
