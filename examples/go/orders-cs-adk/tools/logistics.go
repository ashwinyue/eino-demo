package tools

import (
    "context"
    "fmt"

    "github.com/cloudwego/eino/components/tool"
    "github.com/cloudwego/eino/components/tool/utils"
    "orders-cs-adk/common"
)

type TrackDeliveryInput struct{
    OrderID string `json:"order_id"`
    TrackingNo string `json:"tracking_no"`
}

func NewTrackDeliveryTool(store *common.Store) tool.InvokableTool {
    t, _ := utils.InferOptionableTool("track_delivery", "查询物流信息", func(ctx context.Context, in *TrackDeliveryInput, opts ...tool.Option) (string, error) {
        if in.OrderID == "" && in.TrackingNo == "" { return "缺少订单号或物流单号", nil }
        if in.TrackingNo != "" { return fmt.Sprintf("物流单号 %s 当前在转运中心，预计2天送达", in.TrackingNo), nil }
        o := store.Get(in.OrderID)
        if o == nil { return "未找到该订单", nil }
        if o.Status == "已发货" { return fmt.Sprintf("订单 %s 已发货，物流单号 %s，预计3天送达", o.ID, o.TrackingNo), nil }
        return "该订单尚未发货", nil
    })
    return t
}

