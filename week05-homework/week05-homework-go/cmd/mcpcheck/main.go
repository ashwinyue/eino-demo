package main

import (
	"context"
	"fmt"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"os/exec"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cli := mcp.NewClient(&mcp.Implementation{Name: "mcp-check-client", Version: "v1"}, nil)
	transport := &mcp.CommandTransport{Command: exec.Command("go", "run", "./mcpserver")}
	session, err := cli.Connect(ctx, transport, nil)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	p, err := session.CallTool(ctx, &mcp.CallToolParams{Name: "get_prompt", Arguments: map[string]any{
		"agent_name": "research",
	}})
	if err != nil {
		panic(err)
	}
	var promptText string
	for _, c := range p.Content {
		if tc, ok := c.(*mcp.TextContent); ok {
			promptText += tc.Text
		}
	}
	fmt.Printf("get_prompt(research) len=%d\n", len(promptText))

	s, err := session.CallTool(ctx, &mcp.CallToolParams{Name: "search", Arguments: map[string]any{
		"topic": "AI Agent",
	}})
	if err != nil {
		panic(err)
	}
	var searchText string
	for _, c := range s.Content {
		if tc, ok := c.(*mcp.TextContent); ok {
			searchText += tc.Text
		}
	}
	fmt.Printf("search output:\n%s\n", searchText)
}
