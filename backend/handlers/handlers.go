package handlers

import (
	"database/sql"
	"log"
	"time"

	"github.com/adamdgit/gotest/backend/models"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

var db *sql.DB

// get all posts
func GetAllPosts(c *fiber.Ctx) error {
	stmt := "SELECT * FROM posts"

	rows, err := db.Query(stmt)
	if err != nil {
		log.Printf("Error: %s", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error retrieving from database",
		})
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var post models.Post

		err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.Created_At, &post.Updated_At)
		if err != nil {
			log.Printf("Error: %s", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Error retrieving from database",
			})
		}
		posts = append(posts, post)
	}

	return c.JSON(posts)
}

// Get post by provided id
func GetPostById(c *fiber.Ctx) error {
	id := c.Params("id")

	stmt := "SELECT * FROM posts WHERE post.id = ?"

	row, err := db.Query(stmt, id)
	if err != nil {
		log.Printf("Error: %s", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error retrieving from database",
		})
	}
	defer row.Close()

	var post models.Post

	err = row.Scan(&post.ID, &post.Title, &post.Content, &post.Created_At, &post.Updated_At)
	if err != nil {
		log.Printf("Error: %s", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error retrieving from database",
		})
	}

	return c.JSON(post)
}

func Login(c *fiber.Ctx) error {
	user := c.FormValue("user")
	pass := c.FormValue("pass")

	// Throws Unauthorized error
	if user != "john" || pass != "doe" {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	// Create the Claims
	claims := jwt.MapClaims{
		"name":  "John Doe",
		"admin": true,
		"exp":   time.Now().Add(time.Hour * 72).Unix(),
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{"token": t})
}

func Register(c *fiber.Ctx) error {
	return c.JSON("")
}
