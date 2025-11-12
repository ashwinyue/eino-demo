package agents

import (
    "context"

    "github.com/cloudwego/eino/adk"
    "github.com/cloudwego/eino/adk/prebuilt/supervisor"
    "github.com/cloudwego/eino/components/model"
    "orders-cs-adk/common"
)

func NewSupervisor(ctx context.Context, m model.ToolCallingChatModel, subs []adk.Agent) adk.Agent {
    sv, _ := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
        Name:        "supervisor",
        Description: "主管路由",
        Instruction: common.SupervisorPrompt(),
        Model:       m,
        Exit:        &adk.ExitTool{},
    })
    a, _ := supervisor.New(ctx, &supervisor.Config{Supervisor: sv, SubAgents: subs})
    return a
}
