package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/YoavIsaacs/url_shortener/internal/config"
	"github.com/YoavIsaacs/url_shortener/internal/handler"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
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
		Handler: mux,
		Addr:    servAddr,
	}

	mux.HandleFunc("GET /api/health", handler.HealthCheck)
	mux.HandleFunc("POST /urls", handler.HandleAddURL(apiConfig))
	mux.HandleFunc("POST /admin/reset", handler.HandleDeleteAllURLs(apiConfig))
	mux.HandleFunc("POST /admin/reset-single", handler.HandleDeleteSingleURL(apiConfig))

	err = serv.ListenAndServe()
	if err != nil {
		fmt.Println("error: server error")
	}
}
