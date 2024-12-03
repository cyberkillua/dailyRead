package models

import (
	"time"

	"github.com/cyberkillua/dailyread/internal/database"
	"github.com/google/uuid"
)

type Post struct {
	ID          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Url         string    `json:"url"`
	PublishedAt time.Time `json:"published_at"`
}

func DatabasePostToPost(dbPost database.Post) Post {
	return Post{
		ID:          dbPost.ID,
		CreatedAt:   dbPost.CreatedAt,
		UpdatedAt:   dbPost.UpdatedAt,
		Title:       dbPost.Title,
		Description: dbPost.Description.String,
		Url:         dbPost.Url,
		PublishedAt: dbPost.PublishedAt.Time,
	}
}

func DatabasePostsToPosts(dbPost []database.Post) []Post {
	var posts []Post
	for _, dbPost := range dbPost {
		posts = append(posts, DatabasePostToPost(dbPost))
	}
	return posts
}
