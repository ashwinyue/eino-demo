package tools

import (
    "context"

    "github.com/cloudwego/eino/components/tool"
    "github.com/cloudwego/eino/components/tool/utils"
)

type ClarifyInput struct{
    Question string `json:"question"`
}

func NewClarifyTool() tool.InvokableTool {
    t, _ := utils.InferOptionableTool("ask_for_clarification", "缺少必要信息时提示用户补充", func(ctx context.Context, in *ClarifyInput, opts ...tool.Option) (string, error) {
        return in.Question, nil
    })
    return t
}

