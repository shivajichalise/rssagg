package main

import (
	"fmt"
	"net/http"

	"github.com/shivajichalise/rssagg/internal/auth"
	"github.com/shivajichalise/rssagg/internal/database"
)

type authHandler func(http.ResponseWriter, *http.Request, database.User)

func (apiConf *apiConfig) middlewareAuth(handler authHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := auth.GetApiKey(r.Header)
		if err != nil {
			respondWithError(w, 403, fmt.Sprintf("ERROR: Auth error: %v", err))
			return
		}

		user, err := apiConf.DB.GetUserByAPIKey(r.Context(), apiKey)
		if err != nil {
			respondWithError(w, 404, fmt.Sprintf("ERROR: User not found: %v", err))
			return
		}

		handler(w, r, user)
	}
}
