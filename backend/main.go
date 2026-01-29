package main

import (
	"database/sql"
	"encoding/gob"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/adamdgit/gotest/backend/docs"
	"github.com/adamdgit/gotest/backend/models"
	"github.com/adamdgit/gotest/backend/routes"
	"github.com/adamdgit/gotest/backend/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
	fiberSwagger "github.com/swaggo/fiber-swagger"
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

	// Create PostgreSQL connection string
	conn := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", username, password, address, dbname)

	// Init db with config
	db, err := sql.Open("pgx", conn)
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal("unable to connect to postgres:", err)
	}
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(200)
	db.SetMaxIdleConns(100)
	defer db.Close()

	// Create database tables
	if err := utils.InitDB(db); err != nil {
		log.Fatal("DB init failed:", err)
	}

	// Init Fiber app
	app := fiber.New()

	// Initialise swagger docs
	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	// app.Use(csrf.New(csrf.ConfigDefault))

	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowOrigins:     "http://localhost:5173",
	}))

	// Run daily cron to cleanup expired sessions
	cronJob := cron.New()
	cronJob.AddFunc("0 0 * * *", func() { utils.CleanupExpiredSessions(db) })
	cronJob.Start()

	// FIX: gob encoder error when reading models.UserRole
	gob.Register(models.UserRole(""))

	// Register all API routes
	routes.RegisterAdminRoutes(app, db)
	routes.RegisterAuthRoutes(app, db)
	routes.RegisterUserRoutes(app, db)

	app.Static("/", "./public")

	log.Fatal(app.Listen(PORT))
}
