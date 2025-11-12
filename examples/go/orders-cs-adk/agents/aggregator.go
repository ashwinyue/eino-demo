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

func NewAggregatorAgent(ctx context.Context, m model.ToolCallingChatModel, store *common.Store) adk.Agent {
    agg := t.NewAggregateTool(ctx, store)
    c := t.NewClarifyTool()
    a, _ := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
        Name:        "aggregator_agent",
        Description: "聚合查询",
        Instruction: common.AggregatorPrompt(),
        Model:       m,
        ToolsConfig: adk.ToolsConfig{ToolsNodeConfig: compose.ToolsNodeConfig{Tools: []tool.BaseTool{agg, c}}},
        Exit:        &adk.ExitTool{},
    })
    return a
}
