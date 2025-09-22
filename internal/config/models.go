package config

import (
	"github.com/google/uuid"
)

// Generate random callback data for models
func init() {
	for group, modelList := range Models {
		for i, model := range modelList {
			if model.CallbackData == "" {
				Models[group][i].CallbackData = uuid.New().String()[:8]
			}
		}
	}
}

type AIModel struct {
	Title        string
	CallbackData string
	ApiName      string
	Image        bool
	Provider     ProviderModel
}

type ProviderModel string
type ModelGroup string

const (
	OpenRouter   ProviderModel = "OpenRouter"
	IoNet        ProviderModel = "IoNet"
	Pollinations ProviderModel = "Pollinations"
)

const (
	GPTGroup      ModelGroup = "GPT"
	QwenGroup     ModelGroup = "Qwen"
	DeepSeekGroup ModelGroup = "DeepSeek"
	MistralGroup  ModelGroup = "Mistral"
	GeminiGroup   ModelGroup = "Gemini"
	OtherGroup    ModelGroup = "Other"
)

var gptModels = []AIModel{
	// {Title: "GPT-4.1", ApiName: "openai-large", Image: true, Provider: Pollinations},
	{Title: "GPT-4.1 Nano", ApiName: "openai", Image: true, Provider: Pollinations},
	{Title: "GPT-5 Nano", ApiName: "gpt-5-nano", Image: true, Provider: Pollinations},
	{Title: "GPT-OSS", ApiName: "openai/gpt-oss-120b", Provider: IoNet},
	{Title: "GPT-OSS-20b", ApiName: "openai/gpt-oss-20b", Provider: IoNet},
	{Title: "GPT-OSS-20b 2", ApiName: "openai/gpt-oss-20b:free", Provider: OpenRouter},
	{Title: "o4-mini", ApiName: "openai-reasoning", Provider: Pollinations},
}

var qwenModels = []AIModel{
	{Title: "Qwen3-Next", ApiName: "Qwen/Qwen3-Next-80B-A3B-Instruct", Provider: IoNet},
	{Title: "Qwen3-235B", ApiName: "qwen/qwen3-235b-a22b:free", Provider: OpenRouter},
	{Title: "Qwen3-Coder", ApiName: "Intel/Qwen3-Coder-480B-A35B-Instruct-int4-mixed-ar", Provider: IoNet},
	{Title: "Qwen3-Coder 2", ApiName: "qwen/qwen3-coder:free", Provider: OpenRouter},
	{Title: "Qwen3-235B-Thinking", ApiName: "Qwen/Qwen3-235B-A22B-Thinking-2507", Provider: IoNet},
	{Title: "Qwen3-30B", ApiName: "qwen/qwen3-30b-a3b:free", Provider: OpenRouter},
	{Title: "Qwen2.5 VL", ApiName: "Qwen/Qwen2.5-VL-32B-Instruct", Image: true, Provider: IoNet},
	{Title: "Qwen2.5-FAST", ApiName: "featherless/qwerky-72b:free", Provider: OpenRouter},
}

var deepSeekModels = []AIModel{
	{Title: "Deepseek-R1", ApiName: "deepseek-ai/DeepSeek-R1-0528", Provider: IoNet},
	{Title: "Deepseek-R1 2", ApiName: "deepseek/deepseek-r1:free", Provider: OpenRouter},
	{Title: "Deepseek-R1 3", ApiName: "deepseek-reasoning", Provider: Pollinations},
	{Title: "Deepseek V3.1", ApiName: "deepseek/deepseek-chat-v3.1:free", Provider: OpenRouter},
	{Title: "Deepseek V3", ApiName: "deepseek/deepseek-chat-v3-0324:free", Provider: OpenRouter},
	{Title: "Deepseek-Chimera", ApiName: "tngtech/deepseek-r1t2-chimera:free", Provider: OpenRouter},
	{Title: "Deepseek R1-Qwen3", ApiName: "deepseek/deepseek-r1-0528-qwen3-8b:free", Provider: OpenRouter},
}

var mistralModels = []AIModel{
	{Title: "Mistral Large", ApiName: "mistralai/Mistral-Large-Instruct-2411", Provider: IoNet},
	{Title: "Mistral Small 3.2", ApiName: "mistralai/mistral-small-3.2-24b-instruct:free", Image: true, Provider: OpenRouter},
}

var geminiModels = []AIModel{
	{Title: "Gemma 3", ApiName: "google/gemma-3-27b-it:free", Image: true, Provider: OpenRouter},
	{Title: "Gemini 2.5 Lite", ApiName: "gemini", Provider: Pollinations},
}

var otherModels = []AIModel{
	{Title: "Llama 4 Maverick", ApiName: "meta-llama/Llama-4-Maverick-17B-128E-Instruct-FP8", Provider: IoNet, Image: true},
	{Title: "GLM 4.5 Air", ApiName: "z-ai/glm-4.5-air:free", Provider: OpenRouter},
	{Title: "Kimi K2", ApiName: "moonshotai/kimi-k2:free", Provider: OpenRouter},
	{Title: "Grok 4 fast", ApiName: "x-ai/grok-4-fast:free", Provider: OpenRouter},
}

var PhotoModels = []AIModel{
	{Title: "Flux", ApiName: "flux"},
	{Title: "Flux Pro", ApiName: "flux-pro"},
	{Title: "Flux Dev", ApiName: "flux-dev"},
	{Title: "Kontext", ApiName: "kontext"},
	{Title: "Turbo 18+", ApiName: "turbo"},
	{Title: "GPTImage", ApiName: "gptimage"},
	{Title: "SDXL Turbo", ApiName: "sdxl-turbo"},
}

var ModelGroupOrder = []ModelGroup{
	GPTGroup,
	QwenGroup,
	DeepSeekGroup,
	MistralGroup,
	GeminiGroup,
	OtherGroup,
}

var Models = map[ModelGroup][]AIModel{
	GPTGroup:      gptModels,
	QwenGroup:     qwenModels,
	DeepSeekGroup: deepSeekModels,
	MistralGroup:  mistralModels,
	OtherGroup:    otherModels,
	GeminiGroup:   geminiModels,
}
