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

func NewRefundAgent(ctx context.Context, m model.ToolCallingChatModel, store *common.Store) adk.Agent {
    r := t.NewApplyRefundTool(store)
    q := t.NewClarifyTool()
    a, _ := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
        Name:        "refund_agent",
        Description: "支付/退款",
        Instruction: common.RefundPrompt(),
        Model:       m,
        ToolsConfig: adk.ToolsConfig{ToolsNodeConfig: compose.ToolsNodeConfig{Tools: []tool.BaseTool{r, q}}},
        Exit:        &adk.ExitTool{},
    })
    return a
}
