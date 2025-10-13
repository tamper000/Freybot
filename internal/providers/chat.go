package providers

import (
	"time"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

const chAtURL = "https://ch.at/v1/"

func NewChAtClient(timeout time.Duration) Client {
	client := openai.NewClient(option.WithBaseURL(chAtURL))

	return &OpenaiClient{
		client:  &client,
		timeout: timeout,
	}
}
