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
	LLM7         ProviderModel = "LLM7"
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
	{Title: "GPT-4.1 Nano", ApiName: "openai-fast", Image: true, Provider: Pollinations},
	{Title: "GPT-5 Nano", ApiName: "openai", Image: true, Provider: Pollinations},
	{Title: "GPT-4.1", ApiName: "openai-large", Image: true, Provider: Pollinations},
	{Title: "GPT-5 2", ApiName: "gpt-5-chat", Image: true, Provider: LLM7},
	{Title: "GPT-OSS", ApiName: "openai/gpt-oss-120b", Provider: IoNet},
	{Title: "GPT-OSS-20b", ApiName: "openai/gpt-oss-20b", Provider: IoNet},
	{Title: "GPT-OSS-20b 2", ApiName: "openai/gpt-oss-20b:free", Provider: OpenRouter},
	{Title: "o4-mini", ApiName: "openai-reasoning", Provider: Pollinations, Image: true},
}

var qwenModels = []AIModel{
	{Title: "Qwen3-Next", ApiName: "Qwen/Qwen3-Next-80B-A3B-Instruct", Provider: IoNet},
	{Title: "Qwen3-235B", ApiName: "qwen/qwen3-235b-a22b:free", Provider: OpenRouter},
	{Title: "Qwen3-Coder", ApiName: "Intel/Qwen3-Coder-480B-A35B-Instruct-int4-mixed-ar", Provider: IoNet},
	{Title: "Qwen3-Coder 2", ApiName: "qwen/qwen3-coder:free", Provider: OpenRouter},
	{Title: "Qwen3-235B-Thinking", ApiName: "Qwen/Qwen3-235B-A22B-Thinking-2507", Provider: IoNet},
	{Title: "Qwen3-30B", ApiName: "qwen/qwen3-30b-a3b:free", Provider: OpenRouter},
	{Title: "Qwen2.5 VL", ApiName: "Qwen/Qwen2.5-VL-32B-Instruct", Image: true, Provider: IoNet},
}

var deepSeekModels = []AIModel{
	// {Title: "Deepseek V3.1", ApiName: "deepseek/deepseek-chat-v3.1:free", Provider: OpenRouter},
	{Title: "Deepseek-R1", ApiName: "deepseek-ai/DeepSeek-R1-0528", Provider: IoNet},
	{Title: "Deepseek-R1 2", ApiName: "deepseek/deepseek-r1:free", Provider: OpenRouter},
	{Title: "Deepseek V3", ApiName: "deepseek/deepseek-chat-v3-0324:free", Provider: OpenRouter},
	{Title: "Deepseek-Chimera", ApiName: "tngtech/deepseek-r1t2-chimera:free", Provider: OpenRouter},
}

var mistralModels = []AIModel{
	{Title: "Mistral Large", ApiName: "mistralai/Mistral-Large-Instruct-2411", Provider: IoNet},
	{Title: "Mistral Small 3.2", ApiName: "mistralai/mistral-small-3.2-24b-instruct:free", Image: true, Provider: OpenRouter},
	{Title: "Mistral Small 3.1", ApiName: "mistral-small-3.1-24b-instruct-2503", Image: true, Provider: LLM7},
}

var geminiModels = []AIModel{
	{Title: "Gemma 3", ApiName: "google/gemma-3-27b-it:free", Image: true, Provider: OpenRouter},
	// {Title: "Gemini 2.5 Pro", ApiName: "gemini/gemini-2.5-pro", Provider: LLM7},
	{Title: "Gemini 2.5 Lite", ApiName: "gemini", Provider: Pollinations, Image: true},
	{Title: "Gemini 2.5 Search", ApiName: "gemini-search", Provider: Pollinations, Image: true},
}

var otherModels = []AIModel{
	{Title: "Llama 4 Maverick", ApiName: "meta-llama/Llama-4-Maverick-17B-128E-Instruct-FP8", Provider: IoNet, Image: true},
	{Title: "Kimi K2", ApiName: "moonshotai/Kimi-K2-Instruct-0905", Provider: IoNet},
	{Title: "GLM 4.6", ApiName: "zai-org/GLM-4.6", Provider: IoNet},
	{Title: "GLM 4.5 Air", ApiName: "z-ai/glm-4.5-air:free", Provider: OpenRouter},
	{Title: "GLM 4.5 Flash", ApiName: "glm-4.5-flash", Provider: LLM7},
	{Title: "Minimax M2", ApiName: "minimax/minimax-m2:free", Provider: OpenRouter},
}

var PhotoModels = []AIModel{
	{Title: "Flux", ApiName: "flux"},
	{Title: "Kontext", ApiName: "kontext"},
	{Title: "Turbo 18+", ApiName: "turbo"},
	{Title: "Seedream", ApiName: "seedream"},
	{Title: "NanoBanana", ApiName: "nanobanana"},
	{Title: "GPT", ApiName: "gptimage"},
}

var EditModels = []AIModel{
	{Title: "Qwen", ApiName: "qwen"},
	{Title: "Gemini", ApiName: "gemini"},
	{Title: "Kontext", ApiName: "kontext"},
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
