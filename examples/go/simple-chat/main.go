package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	openai "github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/adk/prebuilt/supervisor"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
)

type Order struct {
	ID         string
	User       string
	Status     string
	Items      []string
	Amount     float64
	Created    time.Time
	TrackingNo string
}

type OrderStore struct {
	orders map[string]*Order
}

func NewOrderStore() *OrderStore {
	return &OrderStore{orders: map[string]*Order{}}
}

func (s *OrderStore) Add(o *Order) {
	s.orders[o.ID] = o
}

func (s *OrderStore) Get(id string) *Order {
	return s.orders[id]
}

func (s *OrderStore) Cancel(id string) (bool, string) {
	o := s.orders[id]
	if o == nil {
		return false, "未找到该订单"
	}
	if o.Status == "待发货" || o.Status == "已支付" {
		o.Status = "已取消"
		return true, "订单已取消"
	}
	return false, "当前状态不可取消"
}

func parseOrderID(q string) string {
	re := regexp.MustCompile(`\d{6,}`)
	m := re.FindString(q)
	return m
}

func classifyIntent(q string) string {
	l := strings.ToLower(q)
	if strings.Contains(l, "查询") || strings.Contains(l, "查看") || strings.Contains(l, "订单状态") {
		return "query"
	}
	if strings.Contains(l, "取消") || strings.Contains(l, "退订") {
		return "cancel"
	}
	if strings.Contains(l, "售后") || strings.Contains(l, "退货") || strings.Contains(l, "换货") || strings.Contains(l, "维修") {
		return "aftersales"
	}
	return "faq"
}

type RouteInfo struct {
	Intent  string
	OrderID string
}

//

func OrderQueryAgent(ctx context.Context, q string, store *OrderStore) (*schema.Message, error) {
	id := parseOrderID(q)
	if id == "" {
		return &schema.Message{Role: schema.Assistant, Content: "请提供订单号，例如：查询订单 20251112001"}, nil
	}
	o := store.Get(id)
	if o == nil {
		return &schema.Message{Role: schema.Assistant, Content: "未找到该订单"}, nil
	}
	s := fmt.Sprintf("订单号：%s\n用户：%s\n状态：%s\n商品：%s\n金额：%.2f\n下单时间：%s", o.ID, o.User, o.Status, strings.Join(o.Items, ","), o.Amount, o.Created.Format("2006-01-02 15:04"))
	return &schema.Message{Role: schema.Assistant, Content: s}, nil
}

func OrderCancelAgent(ctx context.Context, q string, store *OrderStore) (*schema.Message, error) {
	id := parseOrderID(q)
	if id == "" {
		return &schema.Message{Role: schema.Assistant, Content: "请提供要取消的订单号，例如：取消订单 20251111001"}, nil
	}
	ok, msg := store.Cancel(id)
	if ok {
		return &schema.Message{Role: schema.Assistant, Content: msg}, nil
	}
	return &schema.Message{Role: schema.Assistant, Content: msg}, nil
}

func AfterSalesAgent(ctx context.Context, q string) (*schema.Message, error) {
	s := "售后支持指引：\n1）7天无理由退货，商品完好不影响二次销售；\n2）30天内质量问题可换货；\n3）一年质保，提供维修服务；\n请准备订单号、商品照片与问题描述，提交至售后工单。"
	return &schema.Message{Role: schema.Assistant, Content: s}, nil
}

func FAQAgent(ctx context.Context, q string) (*schema.Message, error) {
	l := strings.ToLower(q)
	if strings.Contains(l, "发票") {
		return &schema.Message{Role: schema.Assistant, Content: "支持电子普通发票，开票抬头与税号可在个人中心设置。"}, nil
	}
	if strings.Contains(l, "配送") || strings.Contains(l, "快递") {
		return &schema.Message{Role: schema.Assistant, Content: "默认使用顺丰/京东快递，支付后48小时内发货。"}, nil
	}
	return &schema.Message{Role: schema.Assistant, Content: "请描述您的问题，例如：查询订单、取消订单、售后或常见问题。"}, nil
}

func RouteAndHandle(ctx context.Context, q string, store *OrderStore, ri RouteInfo) (*schema.Message, error) {
	switch ri.Intent {
	case "query":
		if ri.OrderID != "" {
			return OrderQueryAgent(ctx, ri.OrderID, store)
		}
		return OrderQueryAgent(ctx, q, store)
	case "cancel":
		if ri.OrderID != "" {
			return OrderCancelAgent(ctx, ri.OrderID, store)
		}
		return OrderCancelAgent(ctx, q, store)
	case "aftersales":
		return AfterSalesAgent(ctx, q)
	default:
		return FAQAgent(ctx, q)
	}
}

func main() {
	ctx := context.Background()
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Fatal("环境变量 OPENAI_API_KEY 未设置")
	}
	store := NewOrderStore()
	store.Add(&Order{ID: "20251112001", User: "张三", Status: "已支付", Items: []string{"蓝牙耳机"}, Amount: 199.00, Created: time.Now().Add(-24 * time.Hour)})
	store.Add(&Order{ID: "20251112002", User: "李四", Status: "已发货", Items: []string{"机械键盘"}, Amount: 499.00, Created: time.Now().Add(-36 * time.Hour), TrackingNo: "SF1234567890"})
	store.Add(&Order{ID: "20251111001", User: "王五", Status: "待发货", Items: []string{"显示器"}, Amount: 1299.00, Created: time.Now().Add(-48 * time.Hour)})

	samples := []string{"查询订单 20251112001", "帮我取消订单 20251111001", "退货流程怎么走", "是否支持开发票", "查询物流 20251112002", "申请退款 20251112001 原因 不想要了", "开票 20251112001 抬头 测试公司 税号 123456789"}
	cm, _ := openai.NewChatModel(ctx, &openai.ChatModelConfig{APIKey: apiKey, BaseURL: os.Getenv("OPENAI_API_BASE"), Model: "gpt-4o-mini", Timeout: 15 * time.Second})

	type qInput struct {
		OrderID string `json:"order_id"`
	}
	qTool, _ := utils.InferOptionableTool("query_order", "查询订单详情", func(ctx context.Context, in *qInput, opts ...tool.Option) (string, error) {
		o := store.Get(in.OrderID)
		if o == nil {
			return "未找到该订单", nil
		}
		return fmt.Sprintf("订单号：%s\n用户：%s\n状态：%s\n商品：%s\n金额：%.2f\n下单时间：%s", o.ID, o.User, o.Status, strings.Join(o.Items, ","), o.Amount, o.Created.Format("2006-01-02 15:04")), nil
	})

	type cInput struct {
		OrderID string `json:"order_id"`
	}
	cTool, _ := utils.InferOptionableTool("cancel_order", "取消订单", func(ctx context.Context, in *cInput, opts ...tool.Option) (string, error) {
		_, msg := store.Cancel(in.OrderID)
		return msg, nil
	})

	type ClarifyInput struct {
		Question string `json:"question" jsonschema_description:"缺少必要信息时向用户询问的问题"`
	}
	clarifyTool, _ := utils.InferOptionableTool("ask_for_clarification", "缺少必要信息时提示用户补充", func(ctx context.Context, in *ClarifyInput, opts ...tool.Option) (string, error) {
		return in.Question, nil
	})

	qa, _ := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:        "order_query_agent",
		Description: "订单查询",
		Instruction: "仅处理订单查询任务。优先调用 query_order 工具，输入形如 {\"order_id\": \"...\"}；若缺少订单号，请调用 ask_for_clarification 询问用户补充。",
		Model:       cm,
		ToolsConfig: adk.ToolsConfig{ToolsNodeConfig: compose.ToolsNodeConfig{Tools: []tool.BaseTool{qTool, clarifyTool}, UnknownToolsHandler: func(ctx context.Context, name, input string) (string, error) {
			return fmt.Sprintf("未知工具: %s", name), nil
		}}},
	})

	ca, _ := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:        "order_cancel_agent",
		Description: "订单取消",
		Instruction: "仅处理订单取消任务。优先调用 cancel_order 工具，输入形如 {\"order_id\": \"...\"}；若缺少订单号，请调用 ask_for_clarification 询问用户补充。",
		Model:       cm,
		ToolsConfig: adk.ToolsConfig{ToolsNodeConfig: compose.ToolsNodeConfig{Tools: []tool.BaseTool{cTool, clarifyTool}, UnknownToolsHandler: func(ctx context.Context, name, input string) (string, error) {
			return fmt.Sprintf("未知工具: %s", name), nil
		}}},
	})

	asa, _ := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:        "aftersales_agent",
		Description: "售后咨询",
		Instruction: "提供退换货与质保流程说明。",
		Model:       cm,
	})

	faqa, _ := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:        "faq_agent",
		Description: "常见问题",
		Instruction: "回答发票、配送等常见问题。",
		Model:       cm,
	})

	type tInput struct {
		OrderID    string `json:"order_id"`
		TrackingNo string `json:"tracking_no"`
	}
	tTool, _ := utils.InferOptionableTool("track_delivery", "查询物流信息", func(ctx context.Context, in *tInput, opts ...tool.Option) (string, error) {
		if in.OrderID == "" && in.TrackingNo == "" {
			return "缺少订单号或物流单号", nil
		}
		if in.TrackingNo != "" {
			return fmt.Sprintf("物流单号 %s 当前在转运中心，预计2天送达", in.TrackingNo), nil
		}
		o := store.Get(in.OrderID)
		if o == nil {
			return "未找到该订单", nil
		}
		if o.Status == "已发货" {
			return fmt.Sprintf("订单 %s 已发货，物流单号 %s，预计3天送达", o.ID, o.TrackingNo), nil
		}
		return "该订单尚未发货", nil
	})
	la, _ := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:        "logistics_agent",
		Description: "物流查询",
		Instruction: "仅处理物流查询任务。优先调用 track_delivery 工具，输入可包含 order_id 或 tracking_no；若缺少必要参数，请调用 ask_for_clarification 询问用户补充。",
		Model:       cm,
		ToolsConfig: adk.ToolsConfig{ToolsNodeConfig: compose.ToolsNodeConfig{Tools: []tool.BaseTool{tTool, clarifyTool}, UnknownToolsHandler: func(ctx context.Context, name, input string) (string, error) {
			return fmt.Sprintf("未知工具: %s", name), nil
		}}},
	})

	type iInput struct{ OrderID, Title, TaxID string }
	iTool, _ := utils.InferOptionableTool("issue_invoice", "开具电子发票", func(ctx context.Context, in *iInput, opts ...tool.Option) (string, error) {
		if in.OrderID == "" || in.Title == "" || in.TaxID == "" {
			return "缺少开票必要信息", nil
		}
		return fmt.Sprintf("订单 %s 已开具电子发票，抬头 %s 税号 %s", in.OrderID, in.Title, in.TaxID), nil
	})
	ia, _ := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:        "invoice_agent",
		Description: "发票服务",
		Instruction: "处理电子发票开具与说明。优先调用 issue_invoice 工具，输入需包含 {order_id,title,tax_id}；若缺少字段，请调用 ask_for_clarification 询问用户补充。",
		Model:       cm,
		ToolsConfig: adk.ToolsConfig{ToolsNodeConfig: compose.ToolsNodeConfig{Tools: []tool.BaseTool{iTool, clarifyTool}, UnknownToolsHandler: func(ctx context.Context, name, input string) (string, error) {
			return fmt.Sprintf("未知工具: %s", name), nil
		}}},
	})

	type rInput struct{ OrderID, Reason string }
	rTool, _ := utils.InferOptionableTool("apply_refund", "申请退款", func(ctx context.Context, in *rInput, opts ...tool.Option) (string, error) {
		if in.OrderID == "" {
			return "缺少订单号", nil
		}
		o := store.Get(in.OrderID)
		if o == nil {
			return "未找到该订单", nil
		}
		if o.Status == "待发货" || o.Status == "已支付" {
			return "退款申请已受理，预计3-5个工作日原路退回", nil
		}
		if o.Status == "已发货" {
			return "订单已发货，请走售后退货流程", nil
		}
		return "当前状态暂不支持退款", nil
	})
	ra, _ := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:        "refund_agent",
		Description: "支付/退款",
		Instruction: "处理退款申请与说明。优先调用 apply_refund 工具，输入需包含 {order_id,reason}；若缺少字段，请调用 ask_for_clarification 询问用户补充。",
		Model:       cm,
		ToolsConfig: adk.ToolsConfig{ToolsNodeConfig: compose.ToolsNodeConfig{Tools: []tool.BaseTool{rTool, clarifyTool}, UnknownToolsHandler: func(ctx context.Context, name, input string) (string, error) {
			return fmt.Sprintf("未知工具: %s", name), nil
		}}},
	})

	type sInput struct{ Query string }
	policies := []string{
		"售后政策：7天无理由退货，30天质量问题可换货，一年质保。",
		"发票规则：支持电子普通发票，开票抬头与税号需完整填写。",
		"配送说明：支付后48小时内发货，默认顺丰/京东快递。",
		"退款流程：待发货或已支付可申请原路退款，已发货需先走退货流程。",
	}
	sTool, _ := utils.InferOptionableTool("search_policy", "检索售后与FAQ政策", func(ctx context.Context, in *sInput, opts ...tool.Option) (string, error) {
		q := strings.ToLower(in.Query)
		var hits []string
		for _, p := range policies {
			lp := strings.ToLower(p)
			if strings.Contains(lp, "售后") && strings.Contains(q, "售后") {
				hits = append(hits, p)
			}
			if strings.Contains(lp, "发票") && (strings.Contains(q, "发票") || strings.Contains(q, "开票")) {
				hits = append(hits, p)
			}
			if strings.Contains(lp, "配送") && (strings.Contains(q, "配送") || strings.Contains(q, "快递")) {
				hits = append(hits, p)
			}
			if strings.Contains(lp, "退款") && strings.Contains(q, "退款") {
				hits = append(hits, p)
			}
		}
		if len(hits) == 0 {
			return "未检索到相关政策，请完善问题表述。", nil
		}
		return strings.Join(hits, "\n"), nil
	})

	asa, _ = adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:        "aftersales_agent",
		Description: "售后咨询",
		Instruction: "调用 search_policy 检索售后与FAQ政策，根据检索结果生成回复。",
		Model:       cm,
		ToolsConfig: adk.ToolsConfig{ToolsNodeConfig: compose.ToolsNodeConfig{Tools: []tool.BaseTool{sTool}, UnknownToolsHandler: func(ctx context.Context, name, input string) (string, error) {
			return fmt.Sprintf("未知工具: %s", name), nil
		}}},
	})

	faqa, _ = adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:        "faq_agent",
		Description: "常见问题",
		Instruction: "调用 search_policy 检索常见问题政策并生成回复。",
		Model:       cm,
		ToolsConfig: adk.ToolsConfig{ToolsNodeConfig: compose.ToolsNodeConfig{Tools: []tool.BaseTool{sTool}, UnknownToolsHandler: func(ctx context.Context, name, input string) (string, error) {
			return fmt.Sprintf("未知工具: %s", name), nil
		}}},
	})

	sv, _ := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:        "supervisor",
		Description: "主管路由",
		Instruction: "根据用户意图在订单查询、订单取消、售后、FAQ、物流、发票、退款之间进行分派。不自行处理具体工作。",
		Model:       cm,
		Exit:        &adk.ExitTool{},
	})

	a, _ := supervisor.New(ctx, &supervisor.Config{Supervisor: sv, SubAgents: []adk.Agent{qa, ca, asa, faqa, la, ia, ra}})
	r := adk.NewRunner(ctx, adk.RunnerConfig{Agent: a})
	for _, q := range samples {
		it := r.Query(ctx, q)
		for {
			e, ok := it.Next()
			if !ok {
				break
			}
			if e.Output != nil {
				msg, _, _ := adk.GetMessage(e)
				fmt.Println("问题：", q)
				fmt.Println("回复：", msg.Content)
			}
		}
	}
}
