package workflows

import (
    "context"
    "strings"
    "testing"
    "orders-cs-adk/common"
)

func TestAggregatorWorkflow(t *testing.T) {
    store := common.NewDefaultStore()
    r, err := NewAggregatorWorkflowRunner(context.Background(), store)
    if err != nil { t.Fatalf("compile workflow failed: %v", err) }
    out, err := r.Invoke(context.Background(), AggregatorInput{OrderID: "20251112002", Query: "售后"})
    if err != nil { t.Fatalf("invoke workflow failed: %v", err) }
    if !strings.Contains(out.Content, "订单信息：") || !strings.Contains(out.Content, "物流信息：") {
        t.Fatalf("merged output missing sections: %s", out.Content)
    }
}

