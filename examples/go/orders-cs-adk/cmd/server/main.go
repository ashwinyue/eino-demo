package main

import (
	"log"
	"orders-cs-adk/common"
	"orders-cs-adk/internal/server"
	"os"
)

func main() {
	cfg, err := common.LoadConfig(".")
	if err != nil {
		log.Println("failed to load config", err)
		os.Exit(1)
	}
	s := server.New(cfg)
	if err := s.Run(); err != nil {
		log.Println("server run error", err)
		os.Exit(1)
	}
}
