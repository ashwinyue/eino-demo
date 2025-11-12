package app

import (
    "context"

    "github.com/cloudwego/eino/adk"
    "orders-cs-adk/agents"
    "orders-cs-adk/common"
)

type App struct {
    runner *adk.Runner
}

func New(cfg *common.Config) *App {
    ctx := context.Background()
    m := common.NewChatModel(cfg)
    store := common.NewDefaultStore()
    qa := agents.NewOrderQueryAgent(ctx, m, store)
    ca := agents.NewOrderCancelAgent(ctx, m, store)
    la := agents.NewLogisticsAgent(ctx, m, store)
    ia := agents.NewInvoiceAgent(ctx, m)
    ra := agents.NewRefundAgent(ctx, m, store)
    asa := agents.NewAfterSalesAgent(ctx, m, cfg)
    faqa := agents.NewFAQAgent(ctx, m, cfg)
    sa := agents.NewSearchAgent(ctx, m, cfg)
    agg := agents.NewAggregatorAgent(ctx, m, store)
    sup := agents.NewSupervisor(ctx, m, []adk.Agent{qa, ca, la, ia, ra, asa, faqa, sa, agg})
    r := adk.NewRunner(ctx, adk.RunnerConfig{Agent: sup})
    return &App{runner: r}
}

func (a *App) Query(ctx context.Context, q string) (string, error) {
    it := a.runner.Query(ctx, q)
    var last string
    for {
        e, ok := it.Next()
        if !ok {
            break
        }
        if e.Err != nil {
            return "", e.Err
        }
        if e.Output != nil {
            msg, _, _ := adk.GetMessage(e)
            last = msg.Content
        }
    }
    return last, nil
}
