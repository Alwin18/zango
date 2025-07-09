package main

import (
	"log"

	"github.com/Alwin18/zango/internal/config"
	"github.com/Alwin18/zango/internal/server"
)

func main() {
	cfg, err := config.Load("config/config.yaml")
	if err != nil {
		log.Fatalf("Gagal load config: %v", err)
	}

	// server
	err = server.Start(server.ServerConfig{
		Host: cfg.Server.Host,
		Port: cfg.Server.Port,
	})
	if err != nil {
		log.Fatalf("Gagal start server: %v", err)
	}
}
