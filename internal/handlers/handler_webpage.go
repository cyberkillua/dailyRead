package handlers

import (
	"net/http"

	"github.com/cyberkillua/dailyread/internal/database"
)

type APIConfig struct {
	DB *database.Queries
}

func (apiConfig *APIConfig) CreateWebpage(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string `json:"name"`
		URL  string `json:"url"`
		Type string `json:"type"`
	}
}
