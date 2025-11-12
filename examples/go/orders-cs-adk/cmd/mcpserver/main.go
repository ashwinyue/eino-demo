package main

import (
	"context"
	"log"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type PolicyArgs struct {
	Query string `json:"query" jsonschema:"required,description=政策查询关键词"`
}
type SearchArgs struct {
	Query string `json:"query" jsonschema:"required,description=搜索关键词"`
}

func policyTool(ctx context.Context, req *mcp.CallToolRequest, in PolicyArgs) (*mcp.CallToolResult, string, error) {
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
		return mcp.NewToolResultText("未检索到相关政策"), "", nil
	}
	return mcp.NewToolResultText(strings.Join(hits, "\n")), "", nil
}

func searchTool(ctx context.Context, req *mcp.CallToolRequest, in SearchArgs) (*mcp.CallToolResult, string, error) {
	return mcp.NewToolResultText("搜索结果：" + in.Query + "\n1) 示例结果A\n2) 示例结果B"), "", nil
}

func main() {
	srv := mcp.NewServer(&mcp.Implementation{Name: "orders-cs-mcp", Version: "v1"}, nil)
	mcp.AddTool(srv, &mcp.Tool{Name: "policy", Description: "检索售后与FAQ政策"}, mcp.NewStructuredToolHandler(policyTool))
	mcp.AddTool(srv, &mcp.Tool{Name: "search", Description: "外部搜索"}, mcp.NewStructuredToolHandler(searchTool))
	if err := srv.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatal(err)
	}
}
