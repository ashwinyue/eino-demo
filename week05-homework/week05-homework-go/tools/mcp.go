package tools

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	adktool "github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func FetchPromptUsingCommand(ctx context.Context, command string, args []string, agentName string, style string, length int) (string, error) {
	if base := os.Getenv("MCP_BASE_URL"); base != "" {
		name := "get_prompt"
		input := map[string]any{"agent_name": agentName}
		if agentName == "write" {
			name = "get_writing_prompt"
			input = map[string]any{"style": style, "length": length}
		}
		return invokeHTTP(ctx, base, name, input)
	}
	cli := mcp.NewClient(&mcp.Implementation{Name: "writer-client", Version: "v1"}, nil)
	transport := &mcp.CommandTransport{Command: exec.Command(command, args...)}
	session, err := cli.Connect(ctx, transport, nil)
	if err != nil {
		return "", err
	}
	defer session.Close()
	var toolName string
	var arguments map[string]any
	if agentName == "write" {
		toolName = "get_writing_prompt"
		arguments = map[string]any{"style": style, "length": length}
	} else {
		toolName = "get_prompt"
		arguments = map[string]any{"agent_name": agentName}
	}
	res, err := session.CallTool(ctx, &mcp.CallToolParams{Name: toolName, Arguments: arguments})
	if err != nil {
		return "", err
	}
	var out string
	for _, c := range res.Content {
		if tc, ok := c.(*mcp.TextContent); ok {
			out += tc.Text
		}
	}
	return out, nil
}

// MCPSearchInput is the input schema for MCP search tool
type MCPSearchInput struct {
	Topic string `json:"topic" jsonschema_description:"搜索主题"`
}

// NewMCPSearchTool creates an ADK tool that proxies to MCP search

func NewMCPSearchTool(ctx context.Context, command string, args []string) (adktool.InvokableTool, error) {
	return utils.InferTool(
		"search",
		"根据主题进行网络搜索（通过MCP）",
		func(ctx context.Context, input *MCPSearchInput) (string, error) {
			if base := os.Getenv("MCP_BASE_URL"); base != "" {
				return invokeHTTP(ctx, base, "search", map[string]any{"topic": input.Topic})
			}
			cli := mcp.NewClient(&mcp.Implementation{Name: "writer-client", Version: "v1"}, nil)
			transport := &mcp.CommandTransport{Command: exec.Command(command, args...)}
			session, err := cli.Connect(ctx, transport, nil)
			if err != nil {
				return "", err
			}
			defer session.Close()
			res, err := session.CallTool(ctx, &mcp.CallToolParams{Name: "search", Arguments: map[string]any{"topic": input.Topic}})
			if err != nil {
				return "", err
			}
			var out string
			for _, c := range res.Content {
				if tc, ok := c.(*mcp.TextContent); ok {
					out += tc.Text
				}
			}
			return out, nil
		},
	)
}

// DefaultMCPCommand returns command and args from environment variables.
// MCP_COMMAND: executable name or absolute path
// MCP_ARGS: space-separated arguments
func DefaultMCPCommand() (string, []string) {
	cmd := os.Getenv("MCP_COMMAND")
	if cmd == "" {
		cmd = "go"
	}
	raw := os.Getenv("MCP_ARGS")
	if raw == "" {
		return cmd, []string{"run", "./mcpserver"}
	}
	// simple split by spaces
	var args []string
	for _, part := range strings.Fields(raw) {
		args = append(args, part)
	}
	return cmd, args
}

func invokeHTTP(ctx context.Context, base string, name string, input map[string]any) (string, error) {
	body, _ := json.Marshal(map[string]any{"name": name, "input": input})
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, strings.TrimRight(base, "/")+"/invoke", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	hc := &http.Client{Timeout: 8 * time.Second}
	resp, err := hc.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	var out struct {
		Output string `json:"output"`
		Error  string `json:"error"`
	}
	dec := json.NewDecoder(resp.Body)
	_ = dec.Decode(&out)
	if out.Error != "" {
		return "", fmt.Errorf(out.Error)
	}
	return out.Output, nil
}
