package agents

import (
	"context"
	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
)

func NewResearchAgent(ctx context.Context, cm model.ToolCallingChatModel, instruction string, searchTool tool.BaseTool) (adk.Agent, error) {
	return adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:        "ResearchAgent",
		Description: "使用Web搜索收集信息并生成结构化研究报告",
		Instruction: instruction,
		Model:       cm,
		ToolsConfig: adk.ToolsConfig{
			ToolsNodeConfig: compose.ToolsNodeConfig{
				Tools: []tool.BaseTool{searchTool},
			},
		},
		MaxIterations: 6,
	})
}
