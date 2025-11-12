package server

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "orders-cs-adk/common"
    "orders-cs-adk/internal/app"
)

type Server struct {
    eng *gin.Engine
    app *app.App
    cfg *common.Config
}

func New(cfg *common.Config) *Server {
    g := gin.Default()
    a := app.New(cfg)
    s := &Server{eng: g, app: a, cfg: cfg}
    g.GET("/health", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"status": "ok"}) })
    g.POST("/chat", s.chat)
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

func (s *Server) Run() error {
    port := s.cfg.Server.Port
    if port == "" {
        port = "8080"
    }
    return s.eng.Run(":" + port)
}

