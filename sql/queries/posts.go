package sqlite

import (
	"database/sql"

	"github.com/adamdgit/gotest/backend/models"
)

type PostModel struct {
	DB *sql.DB
}

func (m *PostModel) GetAllPosts() ([]models.Post, error) {
	stmt := `SELECT id, title, content, createdat FROM posts`
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	posts := []models.Post{}

	for rows.Next() {
		r := models.Post{}
		err := rows.Scan(&r.ID, &r.Title, &r.Content, &r.Created_At)
		if err != nil {
			return nil, err
		}
		posts = append(posts, r)
	}

	return posts, nil
}
