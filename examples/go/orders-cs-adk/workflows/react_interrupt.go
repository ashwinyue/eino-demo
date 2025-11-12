package workflows

import (
    "context"
    "log"
    "os"

    openai "github.com/cloudwego/eino-ext/components/model/openai"
    "github.com/cloudwego/eino/components/model"
    "github.com/cloudwego/eino/components/prompt"
    "github.com/cloudwego/eino/components/tool"
    "github.com/cloudwego/eino/compose"
    "github.com/cloudwego/eino/schema"
    "github.com/cloudwego/eino/components/tool/utils"
)

type localState struct {
    history []*schema.Message
}

func newInvoiceTemplate() prompt.ChatTemplate {
    return prompt.FromMessages(schema.FString,
        schema.SystemMessage("你是订单客服助手。当用户请求开具电子发票时，使用 issue_invoice 工具进行处理。"),
        schema.UserMessage("请为订单 {order_id} 开具电子发票，抬头 {title} 税号 {tax_id}"),
    )
}

type IssueInvoiceInput struct{
    OrderID string `json:"order_id"`
    Title string `json:"title"`
    TaxID string `json:"tax_id"`
}

func newModel(ctx context.Context) model.ChatModel {
    cm, err := openai.NewChatModel(ctx, &openai.ChatModelConfig{
        APIKey:  os.Getenv("OPENAI_API_KEY"),
        BaseURL: os.Getenv("OPENAI_API_BASE"),
        Model:   os.Getenv("OPENAI_MODEL"),
    })
    if err != nil { log.Fatal(err) }
    issueTool, _ := utils.InferOptionableTool("issue_invoice", "开具电子发票", func(ctx context.Context, in *IssueInvoiceInput, opts ...tool.Option) (string, error) {
        if in.OrderID == "" || in.Title == "" || in.TaxID == "" { return "缺少开票必要信息", nil }
        return "已开具电子发票：订单=" + in.OrderID + " 抬头=" + in.Title + " 税号=" + in.TaxID, nil
    })
    tools := []tool.BaseTool{issueTool}
    var infos []*schema.ToolInfo
    for _, x := range tools { info, _ := x.Info(ctx); infos = append(infos, info) }
    _ = cm.BindTools(infos)
    return cm
}

func newToolsNode(ctx context.Context) *compose.ToolsNode {
    issueTool, _ := utils.InferOptionableTool("issue_invoice", "开具电子发票", func(ctx context.Context, in *IssueInvoiceInput, opts ...tool.Option) (string, error) {
        if in.OrderID == "" || in.Title == "" || in.TaxID == "" { return "缺少开票必要信息", nil }
        return "已开具电子发票：订单=" + in.OrderID + " 抬头=" + in.Title + " 税号=" + in.TaxID, nil
    })
    tn, err := compose.NewToolNode(ctx, &compose.ToolsNodeConfig{Tools: []tool.BaseTool{issueTool}})
    if err != nil { log.Fatal(err) }
    return tn
}

func NewInvoiceApprovalGraph(ctx context.Context, store compose.CheckPointStore) (compose.Runnable[map[string]any, *schema.Message], error) {
    tpl := newInvoiceTemplate()
    cm := newModel(ctx)
    tn := newToolsNode(ctx)

    g := compose.NewGraph[map[string]any, *schema.Message](compose.WithGenLocalState(func(ctx context.Context) *localState { return &localState{} }))
    _ = g.AddChatTemplateNode("ChatTemplate", tpl)
    _ = g.AddChatModelNode("ChatModel", cm,
        compose.WithStatePreHandler(func(ctx context.Context, in []*schema.Message, s *localState) ([]*schema.Message, error) { s.history = append(s.history, in...); return s.history, nil }),
        compose.WithStatePostHandler(func(ctx context.Context, out *schema.Message, s *localState) (*schema.Message, error) { s.history = append(s.history, out); return out, nil }),
    )
    _ = g.AddToolsNode("ToolsNode", tn, compose.WithStatePreHandler(func(ctx context.Context, in *schema.Message, s *localState) (*schema.Message, error) { return s.history[len(s.history)-1], nil }))

    _ = g.AddEdge(compose.START, "ChatTemplate")
    _ = g.AddEdge("ChatTemplate", "ChatModel")
    _ = g.AddEdge("ToolsNode", "ChatModel")
    _ = g.AddBranch("ChatModel", compose.NewGraphBranch(func(ctx context.Context, in *schema.Message) (string, error) {
        if len(in.ToolCalls) > 0 { return "ToolsNode", nil }
        return compose.END, nil
    }, map[string]bool{"ToolsNode": true, compose.END: true}))

    return g.Compile(ctx,
        compose.WithCheckPointStore(store),
        compose.WithInterruptBeforeNodes([]string{"ToolsNode"}),
    )
}
