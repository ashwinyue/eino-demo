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

func NewFAQAgent(ctx context.Context, m model.ToolCallingChatModel, cfg *common.Config) adk.Agent {
    s := t.NewHTTPPolicyTool(cfg)
    a, _ := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
        Name:        "faq_agent",
        Description: "常见问题",
        Instruction: common.FAQPrompt(),
        Model:       m,
        ToolsConfig: adk.ToolsConfig{ToolsNodeConfig: compose.ToolsNodeConfig{Tools: []tool.BaseTool{s}}},
        Exit:        &adk.ExitTool{},
    })
    return a
}
