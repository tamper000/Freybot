package providers

import (
	"time"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

const ioNetBaseUrl = "https://api.intelligence.io.solutions/api/v1/"

func NewIoNewClient(token string, timeout time.Duration) Client {
	client := openai.NewClient(option.WithAPIKey(token), option.WithBaseURL(ioNetBaseUrl))

	return &OpenaiClient{
		client:  &client,
		timeout: timeout,
	}
}
