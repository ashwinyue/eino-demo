package tools

import (
    "context"
    "fmt"

    "github.com/cloudwego/eino/components/tool"
    "github.com/cloudwego/eino/components/tool/utils"
)

type IssueInvoiceInput struct{
    OrderID string `json:"order_id"`
    Title string `json:"title"`
    TaxID string `json:"tax_id"`
}

func NewIssueInvoiceTool() tool.InvokableTool {
    t, _ := utils.InferOptionableTool("issue_invoice", "开具电子发票", func(ctx context.Context, in *IssueInvoiceInput, opts ...tool.Option) (string, error) {
        if in.OrderID == "" || in.Title == "" || in.TaxID == "" { return "缺少开票必要信息", nil }
        return fmt.Sprintf("订单 %s 已开具电子发票，抬头 %s 税号 %s", in.OrderID, in.Title, in.TaxID), nil
    })
    return t
}

