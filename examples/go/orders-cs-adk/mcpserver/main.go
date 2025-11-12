package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    "strings"
)

func policyHandler(w http.ResponseWriter, r *http.Request) {
    q := strings.ToLower(r.URL.Query().Get("q"))
    var hits []string
    policies := []string{
        "售后政策：7天无理由退货，30天质量问题可换货，一年质保。",
        "发票规则：支持电子普通发票，开票抬头与税号需完整填写。",
        "配送说明：支付后48小时内发货，默认顺丰/京东快递。",
        "退款流程：待发货或已支付可申请原路退款，已发货需先走退货流程。",
    }
    for _, p := range policies {
        lp := strings.ToLower(p)
        if strings.Contains(q, "售后") && strings.Contains(lp, "售后") { hits = append(hits, p) }
        if (strings.Contains(q, "发票") || strings.Contains(q, "开票")) && strings.Contains(lp, "发票") { hits = append(hits, p) }
        if (strings.Contains(q, "配送") || strings.Contains(q, "快递")) && strings.Contains(lp, "配送") { hits = append(hits, p) }
        if strings.Contains(q, "退款") && strings.Contains(lp, "退款") { hits = append(hits, p) }
    }
    if len(hits) == 0 { fmt.Fprint(w, "未检索到相关政策"); return }
    fmt.Fprint(w, strings.Join(hits, "\n"))
}

type SearchItem struct{
    Title string `json:"title"`
    URL   string `json:"url"`
    Snippet string `json:"snippet"`
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
    q := r.URL.Query().Get("q")
    w.Header().Set("Content-Type", "application/json")
    items := []SearchItem{
        {Title: "退货政策总览", URL: "https://mock.local/policy/returns", Snippet: "支持7天无理由退货与质量问题换货。"},
        {Title: "质保政策", URL: "https://mock.local/policy/warranty", Snippet: "一年质保，提供维修服务。"},
        {Title: "发票开具指南", URL: "https://mock.local/policy/invoice", Snippet: "电子普通发票，需要抬头与税号。"},
    }
    enc := json.NewEncoder(w)
    _ = enc.Encode(map[string]any{"query": q, "results": items})
}

func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("/policy", policyHandler)
    mux.HandleFunc("/search", searchHandler)
    http.ListenAndServe(":8000", mux)
}
