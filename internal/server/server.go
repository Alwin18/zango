package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/Alwin18/zango/internal/service"
)

type ServerConfig struct {
	Host string
	Port int
}

func Start(cfg ServerConfig) error {
	mux := http.NewServeMux()

	// cek folder public
	publicDir := "public"
	if _, err := os.Stat(publicDir); os.IsNotExist(err) {
		return fmt.Errorf("folder %s tidak ditemukan", publicDir)
	}

	// server static file
	fs := http.FileServer(http.Dir(publicDir))
	mux.Handle("/", fs)

	// start
	mux.HandleFunc("/api/service/start", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		name := r.URL.Query().Get("name")
		err := service.Start(name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		fmt.Fprintf(w, "Service %s dijalankan", name)
	})

	// stop
	mux.HandleFunc("/api/service/stop", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		name := r.URL.Query().Get("name")
		err := service.Stop(name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		fmt.Fprintf(w, "Service %s dihentikan", name)
	})

	// status
	mux.HandleFunc("/api/service/status", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		name := r.URL.Query().Get("name")
		status, err := service.Status(name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		resp := map[string]any{
			"name":   name,
			"active": status,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})

	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	fmt.Println("Menjalankan HTTP Server di http://" + addr)
	fmt.Println("Serving folder:", filepath.Join(publicDir))
	return http.ListenAndServe(addr, mux)
}
