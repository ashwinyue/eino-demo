package tools

import (
    "context"

    "github.com/cloudwego/eino/components/tool"
    "github.com/cloudwego/eino/components/tool/utils"
    "orders-cs-adk/common"
    "orders-cs-adk/workflows"
)

type AggregateInput struct{
    OrderID string `json:"order_id"`
    Query   string `json:"query"`
}

func NewAggregateTool(ctx context.Context, store *common.Store) tool.InvokableTool {
    r, _ := workflows.NewAggregatorWorkflowRunner(ctx, store)
    t, _ := utils.InferOptionableTool("aggregate_info", "聚合获取订单详情、物流状态与相关政策并格式化输出", func(ctx context.Context, in *AggregateInput, opts ...tool.Option) (string, error) {
        out, err := r.Invoke(ctx, workflows.AggregatorInput{OrderID: in.OrderID, Query: in.Query})
        if err != nil { return "聚合失败", nil }
        return out.Content, nil
    })
    return t
}
