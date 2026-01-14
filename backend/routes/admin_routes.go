package routes

import (
	"database/sql"

	"github.com/adamdgit/gotest/backend/handlers/users"
	"github.com/adamdgit/gotest/backend/middleware"
	"github.com/gofiber/fiber/v2"
)

func RegisterAdminRoutes(app *fiber.App, db *sql.DB) {
	// Get user by ID
	app.Get("/api/admin/user/:id",
		middleware.AuthSessionIsValid(db),
		middleware.AuthUserHasRole(db, "admin"),
		users.GetUserDataByID(db),
	)

	// Get all user data by search query
	app.Get("/api/admin/users",
		middleware.AuthSessionIsValid(db),
		middleware.AuthUserHasRole(db, "admin"),
		users.GetUserByQuery(db),
	)
}
