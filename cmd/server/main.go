package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq" // make sure this import is here!
)

func main() {
	err := godotenv.Load("../../internal/.env")
	if err != nil {
		fmt.Println("error: could not open env file...")
		os.Exit(1)
	}

	dbURL := os.Getenv("DB_URL")

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Println("error: could not initialize db:", err)
		return
	}
	defer db.Close()

	// Actually test the connection
	if err := db.Ping(); err != nil {
		fmt.Println("error: could not connect to db:", err)
		return
	}

	fmt.Println("success")
}
