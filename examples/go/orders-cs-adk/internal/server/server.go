package server

import (
    "net/http"
    "context"

    "github.com/gin-gonic/gin"
    "github.com/cloudwego/eino/compose"
    "orders-cs-adk/common"
    "orders-cs-adk/internal/app"
    "orders-cs-adk/workflows"
)

type Server struct {
    eng *gin.Engine
    app *app.App
    cfg *common.Config
    stores map[string]*common.MemCheckPointStore
}

func New(cfg *common.Config) *Server {
    g := gin.Default()
    a := app.New(cfg)
    s := &Server{eng: g, app: a, cfg: cfg, stores: make(map[string]*common.MemCheckPointStore)}
    g.GET("/health", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"status": "ok"}) })
    g.POST("/chat", s.chat)
    g.POST("/approval/invoice/start", s.invoiceStart)
    g.POST("/approval/invoice/resume", s.invoiceResume)
    return s
}

func (s *Server) chat(c *gin.Context) {
    var req struct{ Query string `json:"query"` }
    if err := c.ShouldBindJSON(&req); err != nil || req.Query == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
        return
    }
    out, err := s.app.Query(c.Request.Context(), req.Query)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"query": req.Query, "output": out})
}

func (s *Server) getStore(id string) *common.MemCheckPointStore {
    if st, ok := s.stores[id]; ok {
        return st
    }
    st := common.NewMemCheckPointStore()
    s.stores[id] = st
    return st
}

func (s *Server) invoiceStart(c *gin.Context) {
    var req struct {
        ID      string `json:"id"`
        OrderID string `json:"order_id"`
        Title   string `json:"title"`
        TaxID   string `json:"tax_id"`
    }
    if err := c.ShouldBindJSON(&req); err != nil || req.ID == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
        return
    }
    store := s.getStore(req.ID)
    ctx := context.Background()
    r, err := workflows.NewInvoiceApprovalGraph(ctx, store)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    _, invErr := r.Invoke(ctx, map[string]any{"order_id": req.OrderID, "title": req.Title, "tax_id": req.TaxID}, compose.WithCheckPointID(req.ID))
    if invErr != nil {
        if _, ok := compose.ExtractInterruptInfo(invErr); ok {
            c.JSON(http.StatusOK, gin.H{"status": "interrupt", "id": req.ID})
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"error": invErr.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"status": "done"})
}

func (s *Server) invoiceResume(c *gin.Context) {
    var req struct {
        ID      string `json:"id"`
        OrderID string `json:"order_id"`
        Title   string `json:"title"`
        TaxID   string `json:"tax_id"`
    }
    if err := c.ShouldBindJSON(&req); err != nil || req.ID == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
        return
    }
    store := s.getStore(req.ID)
    ctx := context.Background()
    r, err := workflows.NewInvoiceApprovalGraph(ctx, store)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    out, invErr := r.Invoke(ctx, map[string]any{"order_id": req.OrderID, "title": req.Title, "tax_id": req.TaxID}, compose.WithCheckPointID(req.ID))
    if invErr != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": invErr.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"status": "done", "output": out.Content})
}

func (s *Server) Run() error {
    port := s.cfg.Server.Port
    if port == "" {
        port = "8080"
    }
    return s.eng.Run(":" + port)
}
