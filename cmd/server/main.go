package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/YoavIsaacs/url_shortener/internal/config"
	"github.com/YoavIsaacs/url_shortener/internal/handler"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	fmt.Println("Server starting...")
	// this will exit on error, so no need to handle here...
	apiConfig := config.CreateConfig()

	mux := http.NewServeMux()

	err := godotenv.Load("internal/.env")
	if err != nil {
		fmt.Println("error: error loading .env file")
		return
	}

	port := os.Getenv("PORT")
	servAddr := "localhost:" + port
	serv := http.Server{
		Handler:     mux,
		Addr:        servAddr,
		ReadTimeout: time.Second * 30,
	}

	mux.HandleFunc("GET /api/health", handler.HealthCheck)
	mux.HandleFunc("POST /api/urls", handler.HandleAddURL(apiConfig))
	mux.HandleFunc("POST /api/urls/add-qr", handler.HandleAddQR(apiConfig))
	mux.HandleFunc("POST /admin/reset", handler.HandleDeleteAllURLs(apiConfig))
	mux.HandleFunc("POST /admin/reset-single", handler.HandleDeleteSingleURL(apiConfig))
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		shortURL := r.URL.Path[1:]
		if shortURL == "" || strings.Contains(shortURL, "/") {
			http.NotFound(w, r)
			return
		}
		handler.HandleRedirect(apiConfig, shortURL)(w, r)
	})

	err = serv.ListenAndServe()
	if err != nil {
		fmt.Println("error: server error")
	}
}
