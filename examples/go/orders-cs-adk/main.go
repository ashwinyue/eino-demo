package main

import (
    "bufio"
    "context"
    "fmt"
    "os"
    "strings"

    "github.com/cloudwego/eino/adk"
    "github.com/cloudwego/eino/callbacks"
    "orders-cs-adk/agents"
    "orders-cs-adk/common"
)

func main() {
    ctx := context.Background()
    cfg, _ := common.LoadConfig(".")
    m := common.NewChatModel(cfg)
    store := common.NewDefaultStore()

    callbacks.AppendGlobalHandlers([]callbacks.Handler{&common.LoggerCallbacks{}}...)

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

    r := adk.NewRunner(ctx, adk.RunnerConfig{Agent: sup, EnableStreaming: true})
    fmt.Println("orders-cs-adk 交互模式，输入问题，输入 exit 退出。")
    scanner := bufio.NewScanner(os.Stdin)
    for {
        fmt.Print("> ")
        if !scanner.Scan() {
            break
        }
        q := strings.TrimSpace(scanner.Text())
        if q == "" {
            continue
        }
        if strings.EqualFold(q, "exit") {
            break
        }
        it := r.Query(ctx, q)
        var last string
        for {
            e, ok := it.Next()
            if !ok {
                break
            }
            if e.Err != nil {
                fmt.Println("错误:", e.Err)
                break
            }
            if e.Output != nil {
                msg, _, _ := adk.GetMessage(e)
                last = msg.Content
                if msg.Content != "" {
                    fmt.Println(msg.Content)
                }
            }
        }
        if last == "" {
            fmt.Println("(无输出)")
        }
    }
}
