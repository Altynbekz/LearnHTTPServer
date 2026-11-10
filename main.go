package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"httpServerWithGo/internal/database"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file loaded; using the process environment: %v", err)
	}

	dbURL := os.Getenv("DB_URL")
	platform := os.Getenv("PLATFORM")
	secretKey := os.Getenv("SECRET_KEY")
	if dbURL == "" {
		log.Fatal("DB_URL is not set")
	}
	if secretKey == "" {
		log.Fatal("SECRET_KEY is not set")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("open database: %v", err)
	}
	defer db.Close()

	dbQueries := database.New(db)
	apiCfg := &apiConfig{db: dbQueries, platform: platform, secretKey: secretKey}

	server := http.Server{
		Addr:    ":8080",
		Handler: newMuxWithConfig(apiCfg),
	}
	log.Println("Server is ready to use")
	log.Fatal(server.ListenAndServe())
}
