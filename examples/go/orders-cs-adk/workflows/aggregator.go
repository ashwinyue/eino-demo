package workflows

import (
	"context"
	"fmt"
	"strings"

	"orders-cs-adk/common"

	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
)

type AggregatorInput struct {
	OrderID string
	Query   string
}

type MergeInput struct {
	Q string
	D string
	P string
}

func NewAggregatorWorkflowRunner(ctx context.Context, store *common.Store) (compose.Runnable[AggregatorInput, *schema.Message], error) {
	wf := compose.NewWorkflow[AggregatorInput, *schema.Message]()

    type SOut struct{ Output string }
    lambdaQuery := compose.InvokableLambda(func(ctx context.Context, in AggregatorInput) (SOut, error) {
        if in.OrderID == "" {
            return SOut{Output: "未提供订单号"}, nil
        }
        o := store.Get(in.OrderID)
        if o == nil {
            return SOut{Output: "未找到该订单"}, nil
        }
        return SOut{Output: fmt.Sprintf("订单号：%s\n用户：%s\n状态：%s\n商品：%s\n金额：%.2f\n下单时间：%s", o.ID, o.User, o.Status, strings.Join(o.Items, ","), o.Amount, o.Created.Format("2006-01-02 15:04"))}, nil
    })

    lambdaDelivery := compose.InvokableLambda(func(ctx context.Context, in AggregatorInput) (SOut, error) {
        if in.OrderID == "" {
            return SOut{Output: "未提供订单号"}, nil
        }
        o := store.Get(in.OrderID)
        if o == nil {
            return SOut{Output: "未找到该订单"}, nil
        }
        if o.Status == "已发货" {
            return SOut{Output: fmt.Sprintf("物流单号：%s\n预计送达：3天", o.TrackingNo)}, nil
        }
        return SOut{Output: "该订单尚未发货"}, nil
    })

    lambdaPolicy := compose.InvokableLambda(func(ctx context.Context, in AggregatorInput) (SOut, error) {
        q := strings.ToLower(in.Query)
        var hits []string
		policies := []string{
			"售后政策：7天无理由退货，30天质量问题可换货，一年质保。",
			"发票规则：支持电子普通发票，开票抬头与税号需完整填写。",
			"配送说明：支付后48小时内发货，默认顺丰/京东快递。",
			"退款流程：待发货或已支付可申请原路退款，已发货需先走退货流程。",
		}
		for _, p := range policies {
			lp := strings.ToLower(p)
			if strings.Contains(q, "售后") && strings.Contains(lp, "售后") {
				hits = append(hits, p)
			}
			if (strings.Contains(q, "发票") || strings.Contains(q, "开票")) && strings.Contains(lp, "发票") {
				hits = append(hits, p)
			}
			if (strings.Contains(q, "配送") || strings.Contains(q, "快递")) && strings.Contains(lp, "配送") {
				hits = append(hits, p)
			}
			if strings.Contains(q, "退款") && strings.Contains(lp, "退款") {
				hits = append(hits, p)
			}
		}
        if len(hits) == 0 {
            return SOut{Output: "未检索到相关政策"}, nil
        }
        return SOut{Output: strings.Join(hits, "\n")}, nil
    })

	wf.AddLambdaNode("query", lambdaQuery).AddInput(compose.START)
	wf.AddLambdaNode("delivery", lambdaDelivery).AddInput(compose.START)
	wf.AddLambdaNode("policy", lambdaPolicy).AddInput(compose.START)

	merge := compose.InvokableLambda(func(ctx context.Context, in MergeInput) (*schema.Message, error) {
		var b strings.Builder
		if in.Q != "" {
			b.WriteString("订单信息：\n" + in.Q + "\n\n")
		}
		if in.D != "" {
			b.WriteString("物流信息：\n" + in.D + "\n\n")
		}
		if in.P != "" {
			b.WriteString("相关政策：\n" + in.P + "\n")
		}
		if b.Len() == 0 {
			b.WriteString("未获得可汇总的信息")
		}
		return &schema.Message{Role: schema.Assistant, Content: b.String()}, nil
	})

	wf.AddLambdaNode("merge", merge).
		AddInput("query", compose.MapFields("Output", "Q")).
		AddInput("delivery", compose.MapFields("Output", "D")).
		AddInput("policy", compose.MapFields("Output", "P"))

	wf.End().AddInput("merge")
	r, err := wf.Compile(ctx)
	if err != nil {
		return nil, err
	}
	return r, nil
}
