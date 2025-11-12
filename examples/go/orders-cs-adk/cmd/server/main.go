package main

import (
    "orders-cs-adk/common"
    "orders-cs-adk/internal/server"
)

func main() {
    cfg, _ := common.LoadConfig(".")
    s := server.New(cfg)
    _ = s.Run()
}

