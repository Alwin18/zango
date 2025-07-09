package server

import (
	"fmt"
	"log"
	"net/http"
)

type ServerConfig struct {
	Host string
	Port int
}

func Start(cfg ServerConfig) error {
	mux := http.NewServeMux()

	// handler untuk halaman utama
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "<h1>Selamat datang di Zanggo</h1>")
	})

	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	log.Println("Server berhasil dijalankan di port: ", cfg.Port)
	return http.ListenAndServe(addr, mux)
}
