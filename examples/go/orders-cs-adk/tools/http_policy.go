package tools

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
	"orders-cs-adk/common"
	gosdk "orders-cs-adk/internal/mcp"
	mcpc "orders-cs-adk/internal/mcp"
)

type HTTPPolicyInput struct {
	Query string `json:"query"`
}

func NewHTTPPolicyTool(cfg *common.Config) tool.InvokableTool {
	if cfg.Services.MCPCommand != "" {
		cmd := cfg.Services.MCPCommand
		t, _ := utils.InferOptionableTool("http_policy", "通过 MCP(go-sdk) 调用政策检索服务", func(ctx context.Context, in *HTTPPolicyInput, opts ...tool.Option) (string, error) {
			out, err := gosdk.InvokeUsingCommand(ctx, cmd, []string{}, "policy", map[string]any{"query": in.Query})
			if err != nil {
				return "政策检索失败", nil
			}
			return out, nil
		})
		return t
	}
	if cfg.Services.MCPBaseURL != "" {
		client := mcpc.New(cfg.Services.MCPBaseURL)
		t, _ := utils.InferOptionableTool("http_policy", "通过 MCP 调用政策检索服务", func(ctx context.Context, in *HTTPPolicyInput, opts ...tool.Option) (string, error) {
			out, err := client.Invoke(ctx, "policy", map[string]any{"query": in.Query})
			if err != nil {
				return "政策检索失败", nil
			}
			return out, nil
		})
		return t
	}
	t, _ := utils.InferOptionableTool("http_policy", "通过外部HTTP政策服务检索售后与FAQ政策", func(ctx context.Context, in *HTTPPolicyInput, opts ...tool.Option) (string, error) {
		base := cfg.Services.PolicyAPIURL
		if base == "" {
			policies := []string{
				"售后政策：7天无理由退货，30天质量问题可换货，一年质保。",
				"发票规则：支持电子普通发票，开票抬头与税号需完整填写。",
				"配送说明：支付后48小时内发货，默认顺丰/京东快递。",
				"退款流程：待发货或已支付可申请原路退款，已发货需先走退货流程。",
			}
			q := strings.ToLower(in.Query)
			var hits []string
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
				return "未检索到相关政策", nil
			}
			return strings.Join(hits, "\n"), nil
		}
		u, _ := url.Parse(base)
		q := u.Query()
		q.Set("q", in.Query)
		u.RawQuery = q.Encode()
		req, _ := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
		hc := &http.Client{Timeout: 8 * time.Second}
		resp, err := hc.Do(req)
		if err != nil {
			return "政策检索失败", nil
		}
		defer resp.Body.Close()
		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			b, _ := io.ReadAll(resp.Body)
			if len(b) == 0 {
				return "政策检索失败", nil
			}
			return string(b), nil
		}
		b, _ := io.ReadAll(resp.Body)
		return string(b), nil
	})
	return t
}
