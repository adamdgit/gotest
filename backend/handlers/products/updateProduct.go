package products

import (
	"database/sql"
	"reflect"

	"github.com/gofiber/fiber/v2"
)

// JSON format from login body request
type ProductEdit struct {
	Name        string `json:"name"`
	Brand       string `json:"brand"`
	Description string `json:"description"`
	Price       string `json:"price"`
}

// Get post by provided id
func UpdateProduct(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req ProductEdit
		id := c.Params("id")

		err := c.BodyParser(&req)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request params",
			})
		}

		// Builds dynamic update query based on body params provided
		// We may only want to update name or price instead of all columns
		stmt := "UPDATE products SET "
		args := []interface{}{}
		i := 0

		value := reflect.ValueOf(req)
		key := reflect.TypeOf(req)

		for j := 0; j < value.NumField(); j++ {
			field := value.Field(j)
			fieldName := key.Field(j).Tag.Get("json") // Get JSON tag as column name

			if fieldName == "" {
				continue
			}

			// Only include non-empty values
			if field.Kind() == reflect.String && field.String() != "" {
				if i > 0 {
					stmt += ", "
				}
				stmt += fieldName + " = ?"
				args = append(args, field.String())
				i++
			}
		}

		stmt += " WHERE id = ?"
		args = append(args, id)

		// execute query with dynamic params via args
		_, err = db.Exec(stmt, args...)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Error updating product",
			})
		}

		return c.SendStatus(fiber.StatusOK)
	}
}
