package main

import (
	"database/sql"
	"encoding/gob"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/adamdgit/gotest/backend/models"
	"github.com/adamdgit/gotest/backend/routes"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/middleware/session"
	mysqlStorage "github.com/gofiber/storage/mysql"
	"github.com/joho/godotenv"
)

const PORT = ":8081"

func main() {
	// Load .env db info
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	username := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")
	address := os.Getenv("DB_ADDRESS")
	dbname := os.Getenv("DB_NAME")

	// Create MySQL connection string
	conn := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", username, password, address, dbname)

	// Init db with config
	db, err := sql.Open("mysql", conn)
	if err != nil {
		log.Fatal(err)
	}
	db.SetConnMaxLifetime(time.Minute * 3)
	defer db.Close()

	// Init Fiber app
	app := fiber.New()

	// Create a new session store using MySQL storage
	storage := mysqlStorage.New(mysqlStorage.Config{
		Host:       "127.0.0.1",
		Port:       3306,
		Username:   username,
		Password:   password,
		Database:   dbname,
		Table:      "sessions",
		GCInterval: 10 * time.Minute,
	})

	// Save session store with default config
	store := session.New(session.Config{
		CookieHTTPOnly: true,
		CookieSecure:   false, // true in PROD
		CookieSameSite: "Lax",
		Storage:        storage,
		Expiration:     12 * time.Hour,
	})

	app.Use(csrf.New(csrf.ConfigDefault))

	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowOrigins:     "http://localhost:5173",
	}))

	// FIX: gob encoder error when reading models.UserRole
	gob.Register(models.UserRole(""))

	// Setup all API routes
	routes.RegisterAPIRoutes(app, db, store)

	app.Static("/", "./public")

	log.Fatal(app.Listen(PORT))
}
