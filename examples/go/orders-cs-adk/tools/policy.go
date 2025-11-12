package tools

import (
    "context"
    "strings"

    "github.com/cloudwego/eino/components/tool"
    "github.com/cloudwego/eino/components/tool/utils"
)

type SearchPolicyInput struct{
    Query string `json:"query"`
}

func NewSearchPolicyTool() tool.InvokableTool {
    policies := []string{
        "售后政策：7天无理由退货，30天质量问题可换货，一年质保。",
        "发票规则：支持电子普通发票，开票抬头与税号需完整填写。",
        "配送说明：支付后48小时内发货，默认顺丰/京东快递。",
        "退款流程：待发货或已支付可申请原路退款，已发货需先走退货流程。",
    }
    t, _ := utils.InferOptionableTool("search_policy", "检索售后与FAQ政策", func(ctx context.Context, in *SearchPolicyInput, opts ...tool.Option) (string, error) {
        q := strings.ToLower(in.Query)
        var hits []string
        for _, p := range policies {
            lp := strings.ToLower(p)
            if strings.Contains(q, "售后") && strings.Contains(lp, "售后") { hits = append(hits, p) }
            if (strings.Contains(q, "发票") || strings.Contains(q, "开票")) && strings.Contains(lp, "发票") { hits = append(hits, p) }
            if (strings.Contains(q, "配送") || strings.Contains(q, "快递")) && strings.Contains(lp, "配送") { hits = append(hits, p) }
            if strings.Contains(q, "退款") && strings.Contains(lp, "退款") { hits = append(hits, p) }
        }
        if len(hits) == 0 { return "未检索到相关政策，请完善问题表述。", nil }
        return strings.Join(hits, "\n"), nil
    })
    return t
}

