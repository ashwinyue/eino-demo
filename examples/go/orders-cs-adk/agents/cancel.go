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

func NewOrderCancelAgent(ctx context.Context, m model.ToolCallingChatModel, store *common.Store) adk.Agent {
    c := t.NewCancelOrderTool(store)
    q := t.NewClarifyTool()
    a, _ := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
        Name:        "order_cancel_agent",
        Description: "订单取消",
        Instruction: common.OrderCancelPrompt(),
        Model:       m,
        ToolsConfig: adk.ToolsConfig{ToolsNodeConfig: compose.ToolsNodeConfig{Tools: []tool.BaseTool{c, q}}},
        Exit:        &adk.ExitTool{},
    })
    return a
}
