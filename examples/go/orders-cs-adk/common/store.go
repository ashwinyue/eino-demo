package common

import "time"

type Order struct {
    ID         string
    User       string
    Status     string
    Items      []string
    Amount     float64
    Created    time.Time
    TrackingNo string
    Title      string
    TaxID      string
}

type Store struct {
    orders map[string]*Order
}

func NewDefaultStore() *Store {
    s := &Store{orders: map[string]*Order{}}
    s.orders["20251112001"] = &Order{ID: "20251112001", User: "张三", Status: "已支付", Items: []string{"蓝牙耳机"}, Amount: 199.00, Created: time.Now().Add(-24 * time.Hour)}
    s.orders["20251112002"] = &Order{ID: "20251112002", User: "李四", Status: "已发货", Items: []string{"机械键盘"}, Amount: 499.00, Created: time.Now().Add(-36 * time.Hour), TrackingNo: "SF1234567890"}
    s.orders["20251111001"] = &Order{ID: "20251111001", User: "王五", Status: "待发货", Items: []string{"显示器"}, Amount: 1299.00, Created: time.Now().Add(-48 * time.Hour)}
    return s
}

func (s *Store) Get(id string) *Order {
    return s.orders[id]
}

func (s *Store) Cancel(id string) string {
    o := s.orders[id]
    if o == nil {
        return "未找到该订单"
    }
    if o.Status == "待发货" || o.Status == "已支付" {
        o.Status = "已取消"
        return "订单已取消"
    }
    return "当前状态不可取消"
}

