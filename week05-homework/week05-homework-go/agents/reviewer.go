package agents

import (
	"context"
	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/components/model"
)

func NewReviewerAgent(ctx context.Context, cm model.ToolCallingChatModel, instruction string) (adk.Agent, error) {
	return adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:          "ReviewerAgent",
		Description:   "审核初稿并给出修改建议",
		Instruction:   instruction,
		Model:         cm,
		MaxIterations: 1,
	})
}
