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

func NewLogisticsAgent(ctx context.Context, m model.ToolCallingChatModel, store *common.Store) adk.Agent {
    d := t.NewTrackDeliveryTool(store)
    q := t.NewClarifyTool()
    a, _ := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
        Name:        "logistics_agent",
        Description: "物流查询",
        Instruction: common.LogisticsPrompt(),
        Model:       m,
        ToolsConfig: adk.ToolsConfig{ToolsNodeConfig: compose.ToolsNodeConfig{Tools: []tool.BaseTool{d, q}}},
        Exit:        &adk.ExitTool{},
    })
    return a
}
