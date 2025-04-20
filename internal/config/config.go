package config

import (
	"database/sql"
	"fmt"
	"os"
	"runtime"

	"github.com/YoavIsaacs/url_shortener/internal/internal/db/sqlc"
	"github.com/joho/godotenv"
)

type ApiConfig struct {
	Database *sqlc.Queries
}

func CreateConfig() ApiConfig {
	err := godotenv.Load("internal/.env")
	if err != nil {
		_, file, line, _ := runtime.Caller(0)
		fmt.Println(err)
		fmt.Printf("error: error loading .env file at %s, line: %d\n", file, line)
		os.Exit(1)
	}

	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Println("error: server error")
		os.Exit(1)
	}

	dbQueries := sqlc.New(db)
	return ApiConfig{
		Database: dbQueries,
	}
}
