package sc

import (
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

func NewOpenAIClient() openai.Client {
	return openai.NewClient(option.WithAPIKey(GlobalConfig.OpenAIApiKey))
}

func NewOpenAIChatCompletionNewParams(ms []openai.ChatCompletionMessageParamUnion) openai.ChatCompletionNewParams {
	return openai.ChatCompletionNewParams{
		Messages: ms,
		Model:    openai.ChatModelGPT4o,
		Seed:     openai.Int(0),
	}
}
