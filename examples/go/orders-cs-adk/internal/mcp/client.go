package mcp

import (
    "bytes"
    "context"
    "encoding/json"
    "net/http"
)

type Client struct {
    BaseURL string
    HC      *http.Client
}

func New(base string) *Client {
    return &Client{BaseURL: base, HC: &http.Client{}}
}

type InvokeRequest struct{
    Name string `json:"name"`
    Input map[string]any `json:"input"`
}

type InvokeResponse struct{
    Output string `json:"output"`
}

func (c *Client) Invoke(ctx context.Context, name string, input map[string]any) (string, error) {
    reqBody, _ := json.Marshal(InvokeRequest{Name: name, Input: input})
    req, _ := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseURL+"/invoke", bytes.NewReader(reqBody))
    req.Header.Set("Content-Type", "application/json")
    resp, err := c.HC.Do(req)
    if err != nil { return "", err }
    defer resp.Body.Close()
    var out InvokeResponse
    _ = json.NewDecoder(resp.Body).Decode(&out)
    return out.Output, nil
}

