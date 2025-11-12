package common

import (
    "context"

    openai "github.com/cloudwego/eino-ext/components/model/openai"
    "github.com/cloudwego/eino/components/model"
)

func NewChatModel(cfg *Config) model.ToolCallingChatModel {
    modelName := cfg.OpenAI.Model
    if modelName == "" {
        modelName = "gpt-4o-mini"
    }
    cm, _ := openai.NewChatModel(context.Background(), &openai.ChatModelConfig{
        APIKey:  cfg.OpenAI.APIKey,
        BaseURL: cfg.OpenAI.BaseURL,
        Model:   modelName,
    })
    return cm
}
