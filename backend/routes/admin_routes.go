package routes

import (
	"database/sql"

	"github.com/adamdgit/gotest/backend/handlers"
	"github.com/adamdgit/gotest/backend/middleware"
	"github.com/gofiber/fiber/v2"
)

func RegisterAdminRoutes(app *fiber.App, db *sql.DB) {
	// Get user by ID
	app.Get("/api/admin/user/:id",
		middleware.AuthSessionIsValid(db),
		middleware.AuthUserHasRole(db, "admin"),
		handlers.GetUserDataByID(db),
	)

	// Get all user data by search query
	app.Get("/api/admin/users",
		middleware.AuthSessionIsValid(db),
		middleware.AuthUserHasRole(db, "admin"),
		handlers.GetUserByQuery(db),
	)
}
