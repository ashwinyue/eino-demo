package tools

import (
    "context"
    "fmt"
    "strings"

    "github.com/cloudwego/eino/components/tool"
    "github.com/cloudwego/eino/components/tool/utils"
    "orders-cs-adk/common"
)

type QueryOrderInput struct{
    OrderID string `json:"order_id"`
}

func NewQueryOrderTool(store *common.Store) tool.InvokableTool {
    t, _ := utils.InferOptionableTool("query_order", "查询订单详情", func(ctx context.Context, in *QueryOrderInput, opts ...tool.Option) (string, error) {
        o := store.Get(in.OrderID)
        if o == nil { return "未找到该订单", nil }
        return fmt.Sprintf("订单号：%s\n用户：%s\n状态：%s\n商品：%s\n金额：%.2f\n下单时间：%s", o.ID, o.User, o.Status, strings.Join(o.Items, ","), o.Amount, o.Created.Format("2006-01-02 15:04")), nil
    })
    return t
}

type CancelOrderInput struct{
    OrderID string `json:"order_id"`
}

func NewCancelOrderTool(store *common.Store) tool.InvokableTool {
    t, _ := utils.InferOptionableTool("cancel_order", "取消订单", func(ctx context.Context, in *CancelOrderInput, opts ...tool.Option) (string, error) {
        return store.Cancel(in.OrderID), nil
    })
    return t
}

