package tools

import (
    "context"
    "io"
    "net/http"
    "net/url"
    "time"

    "github.com/cloudwego/eino/components/tool"
    "github.com/cloudwego/eino/components/tool/utils"
    "orders-cs-adk/common"
    mcpc "orders-cs-adk/internal/mcp"
)

type HTTPSearchInput struct{
    Query string `json:"query"`
}

func NewHTTPSearchTool(cfg *common.Config) tool.InvokableTool {
    // Prefer MCP invoke if configured
    if cfg.Services.MCPBaseURL != "" {
        client := mcpc.New(cfg.Services.MCPBaseURL)
        t, _ := utils.InferOptionableTool("http_search", "通过 MCP 调用外部检索服务", func(ctx context.Context, in *HTTPSearchInput, opts ...tool.Option) (string, error) {
            out, err := client.Invoke(ctx, "search", map[string]any{"query": in.Query})
            if err != nil { return "检索失败", nil }
            return out, nil
        })
        return t
    }
    // Fallback to direct HTTP
    t, _ := utils.InferOptionableTool("http_search", "通过外部HTTP检索服务搜索信息", func(ctx context.Context, in *HTTPSearchInput, opts ...tool.Option) (string, error) {
        base := cfg.Services.SearchAPIURL
        if base == "" { return "检索服务未配置", nil }
        u, _ := url.Parse(base)
        q := u.Query(); q.Set("q", in.Query); u.RawQuery = q.Encode()
        req, _ := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
        hc := &http.Client{Timeout: 8 * time.Second}
        resp, err := hc.Do(req); if err != nil { return "检索失败", nil }
        defer resp.Body.Close()
        b, _ := io.ReadAll(resp.Body)
        return string(b), nil
    })
    return t
}
