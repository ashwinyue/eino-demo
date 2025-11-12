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

func NewOrderQueryAgent(ctx context.Context, m model.ToolCallingChatModel, store *common.Store) adk.Agent {
    q := t.NewQueryOrderTool(store)
    c := t.NewClarifyTool()
    a, _ := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
        Name:        "order_query_agent",
        Description: "订单查询",
        Instruction: common.OrderQueryPrompt(),
        Model:       m,
        ToolsConfig: adk.ToolsConfig{ToolsNodeConfig: compose.ToolsNodeConfig{Tools: []tool.BaseTool{q, c}}},
        Exit:        &adk.ExitTool{},
    })
    return a
}
