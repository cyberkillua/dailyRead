// package main

// import (
// 	"fmt"
// 	"net/http"

// 	"github.com/cyberkillua/go-api-template/internal/auth"
// 	"github.com/cyberkillua/go-api-template/internal/utils"
// 	"github.com/cyberkillua/go-api-template/internal/database"
// )

// type authHandler func(http.ResponseWriter, *http.Request, database.User)

// func (apiConfig apiConfig) middlewareAuth(handler authHandler) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		apikey, err := auth.GetAPIKey(r.Header)

// 		if err != nil {
// 			utils.RespondWithError(w, http.StatusUnauthorized, fmt.Sprint("Error getting API key: %v", err))
// 			return
// 		}

// 		user, err := apiConfig.DB.GetUserByAPIKey(r.Context(), apikey)

// 		if err != nil {
// 			utils.RespondWithError(w, http.StatusNotFound, fmt.Sprint("Error getting user: %v", err))
// 			return
// 		}

// 		handler(w, r, user)
// 	}
// }

package middleware
