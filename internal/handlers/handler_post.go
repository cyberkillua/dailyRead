package handlers

import (
	"fmt"
	"net/http"

	"github.com/cyberkillua/dailyread/internal/models"
	"github.com/cyberkillua/dailyread/internal/utils"
)

func (apiConfig *APIConfig) GetPost(w http.ResponseWriter, r *http.Request) {
	posts, err := apiConfig.DB.GetPosts(r.Context())
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprint("Error getting posts: %v", err))
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, models.DatabasePostsToPosts(posts))
}
