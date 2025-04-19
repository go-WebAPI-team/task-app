package config

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var DB *sql.DB

func ConnectDB() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading  .env file")
	}

	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_POST")
	database := os.Getenv("DB_NAME")

	if user == "" || password == "" || host == "" || port == "" || database == "" {
		log.Fatal("Database configuration is missing isn .env file")
	}
	dsn := user + ":" + password + "@tcp(" + host + ":" + port + ")/" + database + "?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("fFailed to connect to database: ", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal("Failed to Ping database: ", err)
	}

	DB = db
	log.Println("Database connection established")

}
