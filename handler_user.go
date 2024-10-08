package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/shivajichalise/rssagg/internal/database"
)

// handlerCreateUser() is now a method of apiConfig struct
func (apiConf *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters = struct {
		Name string `json:"name"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}

	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("ERROR: Could not parse JSON: %v", err))
		return
	}

	user, err := apiConf.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		Name:      params.Name,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("ERROR: Could not create user: %v", err))
		return
	}

	respondWithJson(w, 201, databaseUserToUser(user))
}

// handlerCreateUser() is now a method of apiConfig struct
func (apiConf *apiConfig) handlerGetUserByApiKey(w http.ResponseWriter, r *http.Request, user database.User) {
	respondWithJson(w, 200, databaseUserToUser(user))
}

func (apiConf *apiConfig) handlerGetPostsForUser(w http.ResponseWriter, r *http.Request, user database.User) {
	posts, err := apiConf.DB.GetPostsForUser(r.Context(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  10,
	})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("ERROR: Could not fetch posts: %v", err))
		return
	}

	respondWithJson(w, 200, databasePostsToPosts(posts))
}
