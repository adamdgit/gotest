package routes

import (
	"database/sql"
	"time"

	"github.com/adamdgit/gotest/backend/api"
	"github.com/adamdgit/gotest/backend/handlers"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

// Per IP rate limiter, each call has its own counter storage
func rateLimit(max int, window time.Duration) fiber.Handler {
	return limiter.New(limiter.Config{
		Max:               max,
		Expiration:        window,
		LimiterMiddleware: limiter.SlidingWindow{},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(
				api.ErrRateLimited,
			)
		},
	})
}

func RegisterAuthRoutes(app *fiber.App, db *sql.DB) {
	// Refresh session using refresh token
	app.Get("/api/auth/refresh", rateLimit(30, time.Minute), handlers.RefreshAccessToken(db))

	app.Post("/api/auth/login", rateLimit(10, time.Minute), handlers.Login(db))

	app.Get("/api/auth/logout", handlers.Logout(db))

	app.Post("/api/auth/register", rateLimit(5, 10*time.Minute), handlers.Register(db))
}
