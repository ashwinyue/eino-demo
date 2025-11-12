package agents

import (
	"context"
	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/components/model"
)

func NewPolisherAgent(ctx context.Context, cm model.ToolCallingChatModel, instruction string) (adk.Agent, error) {
	return adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:          "PolisherAgent",
		Description:   "根据审核建议润色并输出终稿",
		Instruction:   instruction,
		Model:         cm,
		MaxIterations: 1,
	})
}
