package models

import (
	"time"

	"github.com/cyberkillua/dailyread/internal/database"
	"github.com/google/uuid"
)

type Webpage struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	Url       string    `json:"url"`
	Type      string    `json:"type"`
}

func DatabaseWebpageToWebpage(dbWebpage database.Webpage) Webpage {
	return Webpage{
		ID:        dbWebpage.ID,
		CreatedAt: dbWebpage.CreatedAt,
		UpdatedAt: dbWebpage.UpdatedAt,
		Name:      dbWebpage.Name,
		Url:       dbWebpage.Url,
		Type:      dbWebpage.Type,
	}
}
