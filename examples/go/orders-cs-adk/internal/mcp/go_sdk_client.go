package mcp

import (
	"context"
	"fmt"
	"os/exec"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func InvokeUsingCommand(ctx context.Context, command string, args []string, name string, input map[string]any) (string, error) {
	cli := mcp.NewClient(&mcp.Implementation{Name: "orders-cs-client", Version: "v1"}, nil)
	transport := &mcp.CommandTransport{Command: exec.Command(command, args...)}
	session, err := cli.Connect(ctx, transport, nil)
	if err != nil {
		return "", err
	}
	defer session.Close()
	res, err := session.CallTool(ctx, &mcp.CallToolParams{Name: name, Arguments: input})
	if err != nil {
		return "", err
	}
	if res.IsError {
		var msg string
		for _, c := range res.Content {
			if tc, ok := c.(*mcp.TextContent); ok {
				msg += tc.Text
			}
		}
		if msg == "" {
			msg = "MCP tool error"
		}
		return "", fmt.Errorf(msg)
	}
	var out string
	for _, c := range res.Content {
		if tc, ok := c.(*mcp.TextContent); ok {
			out += tc.Text
		}
	}
	return out, nil
}
