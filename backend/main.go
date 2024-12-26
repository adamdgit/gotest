package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	"github.com/adamdgit/gotest/backend/routes"
)

const PORT = ":8081"

func main() {
	username := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")
	address := os.Getenv("DB_ADDRESS")
	dbname := os.Getenv("DB_NAME")

	// Initialise MySQL connection string
	conn := fmt.Sprintf("%s:%s@tcp(%s)/%s", username, password, address, dbname)

	db, err := sql.Open("mysql", conn)
	if err != nil {
		log.Fatal(err)
	}
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	defer db.Close()

	// Initialise Fiber
	app := fiber.New()

	// Handle CORS
	app.Use(cors.New())

	// Handle static folder
	app.Static("/", "./public")

	// Setup all API routes
	routes.RegisterAPIRoutes(app, db)

	// Listen on port
	log.Fatal(app.Listen(PORT))
	println("Listening on ", PORT)
}
