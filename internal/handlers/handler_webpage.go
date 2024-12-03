package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/cyberkillua/dailyread/internal/database"
	"github.com/cyberkillua/dailyread/internal/models"
	"github.com/cyberkillua/dailyread/internal/utils"
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

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	webpage, err := apiConfig.DB.CreateWebpage(r.Context(), database.CreateWebpageParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
		Url:       params.URL,
		Type:      params.Type,
	})

	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error creating webpage: %v", err))
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, models.DatabaseWebpageToWebpage(webpage))

}
