package queries

import (
	"database/sql"

	"github.com/adamdgit/gotest/backend/models"
	"github.com/gofiber/fiber/v2"
)

// get all posts
func GetAllPosts(app *fiber.App, db *sql.DB) ([]models.Post, error) {
	stmt := "SELECT * FROM posts"

	rows, err := db.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var post models.Post

		err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.Created_At, &post.Updated_At)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}

// Get post by provided id
func getPostByID() {
	return
}
