package agents

import (
	"context"
	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/components/model"
)

func NewWritingAgent(ctx context.Context, cm model.ToolCallingChatModel, instruction string) (adk.Agent, error) {
	return adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:          "WritingAgent",
		Description:   "基于研究报告生成文章初稿",
		Instruction:   instruction,
		Model:         cm,
		MaxIterations: 1,
	})
}
