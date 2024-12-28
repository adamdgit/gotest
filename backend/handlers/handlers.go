package handlers

import (
	"database/sql"
	"log"
	"time"

	"github.com/adamdgit/gotest/backend/models"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
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

func CheckPasswordHash(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func Login(c *fiber.Ctx) error {
	// Get form data
	username := c.FormValue("username")
	password := c.FormValue("password")

	// Get username and password from DB
	stmt := "SELECT username, password FROM users WHERE users.username = ?"

	row, err := db.Query(stmt, username)
	if err != nil {
		log.Printf("Error: %s", err)
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	defer row.Close()

	var user models.User

	err = row.Scan(&user.Username, &user.Password)
	if err != nil {
		log.Printf("Error: %s", err)
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	// Check password matches the hash
	hash := user.Password
	match := CheckPasswordHash(password, hash)

	if !match {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	// Create the Claims
	claims := jwt.MapClaims{
		"name":  username,
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

	// password := "secret"
	// hash, _ := HashPassword(password)

	// func HashPassword(password string) (string, error) {
	//   bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	//   return string(bytes), err
	// }

	return c.JSON("")
}
