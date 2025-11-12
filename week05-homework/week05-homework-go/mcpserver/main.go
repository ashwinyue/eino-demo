package main

import (
	"context"
	"encoding/json"
	"github.com/PuerkitoBio/goquery"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"net/http"
	"strings"
	"week05-homework-go/prompts"
)

type SearchArgs struct {
	Topic string `json:"topic" jsonschema:"搜索主题"`
}
type GetPromptArgs struct {
	AgentName string `json:"agent_name" jsonschema:"代理名称"`
}
type GetWritingPromptArgs struct {
	Style  string `json:"style" jsonschema:"写作风格"`
	Length int    `json:"length" jsonschema:"写作长度"`
}

func searchJSON(topic string) string {
	t := strings.TrimSpace(topic)
	if t == "" {
		return "[]"
	}
	// Try real DuckDuckGo HTML parsing
	url := "https://duckduckgo.com/html/?q=" + strings.ReplaceAll(t, " ", "+")
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0 Safari/537.36")
	hc := &http.Client{Timeout: 8 * 1e9}
	resp, err := hc.Do(req)
	if err == nil && resp.StatusCode >= 200 && resp.StatusCode < 300 {
		defer resp.Body.Close()
		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err == nil {
			var results []map[string]string
			doc.Find(".result__body").Each(func(i int, s *goquery.Selection) {
				if len(results) >= 5 {
					return
				}
				a := s.Find("a.result__a")
				href, _ := a.Attr("href")
				title := strings.TrimSpace(a.Text())
				body := strings.TrimSpace(s.Find(".result__snippet").Text())
				if title != "" && href != "" {
					results = append(results, map[string]string{"title": title, "body": body, "href": href})
				}
			})
			if len(results) > 0 {
				b, _ := json.Marshal(results)
				return string(b)
			}
		}
	}
	// Fallback to static sample
	samples := []map[string]string{
		{"title": t + " - 示例结果A", "body": "示例摘要A", "href": "https://example.com/a"},
		{"title": t + " - 示例结果B", "body": "示例摘要B", "href": "https://example.com/b"},
	}
	b, _ := json.Marshal(samples)
	return string(b)
}

func searchTool(ctx context.Context, req *mcp.CallToolRequest, in SearchArgs) (*mcp.CallToolResult, any, error) {
	t := strings.TrimSpace(in.Topic)
	out := searchJSON(t)
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: out}}}, nil, nil
}

func promptByName(name string) string {
	switch strings.TrimSpace(name) {
	case "research":
		return prompts.ResearchPrompt
	case "write":
		return prompts.WritingPrompt("通俗易懂", 1000)
	case "review":
		return prompts.ReviewPrompt
	case "polish":
		return prompts.PolishingPrompt
	default:
		return "Error: Prompt not found."
	}
}

func getPromptTool(ctx context.Context, req *mcp.CallToolRequest, in GetPromptArgs) (*mcp.CallToolResult, any, error) {
	p := promptByName(in.AgentName)
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: p}}}, nil, nil
}

func getWritingPromptTool(ctx context.Context, req *mcp.CallToolRequest, in GetWritingPromptArgs) (*mcp.CallToolResult, any, error) {
	sty := strings.TrimSpace(in.Style)
	if sty == "" {
		sty = "通俗易懂"
	}
	ln := in.Length
	if ln <= 0 {
		ln = 1000
	}
	p := prompts.WritingPrompt(sty, ln)
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: p}}}, nil, nil
}

type InvokeRequest struct {
	Name  string         `json:"name"`
	Input map[string]any `json:"input"`
}

type InvokeResponse struct {
	Output string `json:"output"`
	Error  string `json:"error,omitempty"`
}

func runHTTP(addr string) error {
	handler := func(w http.ResponseWriter, r *http.Request) {
		var req InvokeRequest
		dec := json.NewDecoder(r.Body)
		_ = dec.Decode(&req)
		var out string
		switch req.Name {
		case "search":
			if v, ok := req.Input["topic"].(string); ok {
				out = searchJSON(v)
			} else {
				out = "[]"
			}
		case "get_prompt":
			if v, ok := req.Input["agent_name"].(string); ok {
				out = promptByName(v)
			} else {
				out = "Error"
			}
		case "get_writing_prompt":
			sty := "通俗易懂"
			ln := 1000
			if v, ok := req.Input["style"].(string); ok {
				sty = v
			}
			if v, ok := req.Input["length"].(float64); ok {
				ln = int(v)
			}
			out = prompts.WritingPrompt(sty, ln)
		default:
			b, _ := json.Marshal(InvokeResponse{Error: "unknown tool"})
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write(b)
			return
		}
		b, _ := json.Marshal(InvokeResponse{Output: out})
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(b)
	}
	http.HandleFunc("/invoke", handler)
	http.HandleFunc("/mcp/invoke", handler)
	return http.ListenAndServe(addr, nil)
}

func main() {
	srv := mcp.NewServer(&mcp.Implementation{Name: "writer-mcp", Version: "v1"}, nil)
	mcp.AddTool(srv, &mcp.Tool{Name: "search", Description: "根据主题进行网络搜索"}, searchTool)
	mcp.AddTool(srv, &mcp.Tool{Name: "get_prompt", Description: "根据代理名称获取提示词"}, getPromptTool)
	mcp.AddTool(srv, &mcp.Tool{Name: "get_writing_prompt", Description: "获取带风格和长度的写作提示词"}, getWritingPromptTool)

	handler := mcp.NewStreamableHTTPHandler(
		func(_ *http.Request) *mcp.Server { return srv },
		&mcp.StreamableHTTPOptions{},
	)
	http.HandleFunc("/mcp", handler.ServeHTTP)
	go func() { _ = http.ListenAndServe(":8000", nil) }()

	_ = srv.Run(context.Background(), &mcp.StdioTransport{})
}
