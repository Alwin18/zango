package main

import (
	"log"

	"github.com/Alwin18/zango/internal/config"
	"github.com/Alwin18/zango/internal/gui"
	"github.com/Alwin18/zango/internal/server"
	"github.com/Alwin18/zango/internal/storage"
)

func main() {
	cfg, err := config.Load("config/config.yaml")
	if err != nil {
		log.Fatalf("Gagal load config: %v", err)
	}

	// init Db
	if err := storage.InitDB("internal/storage/zango.db"); err != nil {
		log.Fatalf("Gagal init DB: %v", err)
	}

	// server
	go func() {
		err := server.Start(server.ServerConfig{
			Host: cfg.Server.Host,
			Port: cfg.Server.Port,
		})
		if err != nil {
			log.Println("Server error:", err)
		}
	}()

	gui.Start()
}
