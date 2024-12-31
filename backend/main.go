package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/adamdgit/gotest/backend/routes"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

const PORT = ":8081"

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	username := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")
	address := os.Getenv("DB_ADDRESS")
	dbname := os.Getenv("DB_NAME")

	// Init MySQL connection string
	conn := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", username, password, address, dbname)

	// Init db with config
	db, err := sql.Open("mysql", conn)
	if err != nil {
		log.Fatal(err)
	}
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	defer db.Close()

	// Init Fiber app
	app := fiber.New()

	// Setup CORS config
	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowOrigins:     "http://localhost:5173",
	}))

	// Setup all API routes
	routes.RegisterAPIRoutes(app, db)

	// Handle static folder
	app.Static("/", "./public")

	// Listen on port
	log.Fatal(app.Listen(PORT))
}
