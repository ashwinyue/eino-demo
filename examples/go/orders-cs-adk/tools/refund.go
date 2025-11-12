package tools

import (
    "context"

    "github.com/cloudwego/eino/components/tool"
    "github.com/cloudwego/eino/components/tool/utils"
    "orders-cs-adk/common"
)

type ApplyRefundInput struct{
    OrderID string `json:"order_id"`
    Reason string `json:"reason"`
}

func NewApplyRefundTool(store *common.Store) tool.InvokableTool {
    t, _ := utils.InferOptionableTool("apply_refund", "申请退款", func(ctx context.Context, in *ApplyRefundInput, opts ...tool.Option) (string, error) {
        if in.OrderID == "" { return "缺少订单号", nil }
        o := store.Get(in.OrderID)
        if o == nil { return "未找到该订单", nil }
        if o.Status == "待发货" || o.Status == "已支付" { return "退款申请已受理，预计3-5个工作日原路退回", nil }
        if o.Status == "已发货" { return "订单已发货，请走售后退货流程", nil }
        return "当前状态暂不支持退款", nil
    })
    return t
}

