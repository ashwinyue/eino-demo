package agents

import (
    "context"

    "github.com/cloudwego/eino/adk"
    "github.com/cloudwego/eino/components/model"
    "github.com/cloudwego/eino/compose"
    "github.com/cloudwego/eino/components/tool"
    "orders-cs-adk/common"
    t "orders-cs-adk/tools"
)

func NewAfterSalesAgent(ctx context.Context, m model.ToolCallingChatModel, cfg *common.Config) adk.Agent {
    s := t.NewHTTPPolicyTool(cfg)
    a, _ := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
        Name:        "aftersales_agent",
        Description: "售后咨询",
        Instruction: common.AfterSalesPrompt(),
        Model:       m,
        ToolsConfig: adk.ToolsConfig{ToolsNodeConfig: compose.ToolsNodeConfig{Tools: []tool.BaseTool{s}}},
        Exit:        &adk.ExitTool{},
    })
    return a
}
