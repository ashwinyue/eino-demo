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

func NewInvoiceAgent(ctx context.Context, m model.ToolCallingChatModel) adk.Agent {
    i := t.NewIssueInvoiceTool()
    q := t.NewClarifyTool()
    a, _ := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
        Name:        "invoice_agent",
        Description: "发票服务",
        Instruction: common.InvoicePrompt(),
        Model:       m,
        ToolsConfig: adk.ToolsConfig{ToolsNodeConfig: compose.ToolsNodeConfig{Tools: []tool.BaseTool{i, q}}},
        Exit:        &adk.ExitTool{},
    })
    return a
}
