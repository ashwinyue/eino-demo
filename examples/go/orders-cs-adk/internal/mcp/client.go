package mcp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Client struct {
	BaseURL string
	HC      *http.Client
}

func New(base string) *Client {
	return &Client{BaseURL: base, HC: &http.Client{Timeout: 8 * time.Second}}
}

type InvokeRequest struct {
	Name  string         `json:"name"`
	Input map[string]any `json:"input"`
}

type InvokeResponse struct {
	Output string `json:"output"`
}

func (c *Client) Invoke(ctx context.Context, name string, input map[string]any) (string, error) {
	reqBody, _ := json.Marshal(InvokeRequest{Name: name, Input: input})
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseURL+"/invoke", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.HC.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var msg struct {
			Error string `json:"error"`
		}
		dec := json.NewDecoder(resp.Body)
		if err := dec.Decode(&msg); err != nil {
			return "", fmt.Errorf("mcp http status %d", resp.StatusCode)
		}
		if msg.Error == "" {
			return "", fmt.Errorf("mcp http status %d", resp.StatusCode)
		}
		return "", fmt.Errorf(msg.Error)
	}
	var out InvokeResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return "", err
	}
	return out.Output, nil
}
