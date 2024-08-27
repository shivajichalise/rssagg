package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/shivajichalise/rssagg/internal/database"
)

func (apiConf *apiConfig) handlerCreateFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters = struct {
		FeedID uuid.UUID `json:"feed_id"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}

	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("ERROR: Could not parse JSON: %v", err))
		return
	}

	feedFollow, err := apiConf.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		UserID:    user.ID,
		FeedID:    params.FeedID,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("ERROR: Could not create feed follow: %v", err))
		return
	}

	respondWithJson(w, 201, databaseFeedFollowToFeedFollow(feedFollow))
}

func (apiConf *apiConfig) handlerGetFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	feedFollows, err := apiConf.DB.GetFeedFollows(r.Context(), user.ID)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("ERROR: Could not get feed follows: %v", err))
		return
	}

	respondWithJson(w, 201, databaseFeedFollowsToFeedFollows(feedFollows))
}

func (apiConf *apiConfig) handlerDeleteFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	id := chi.URLParam(r, "id")
	feedFollowId, err := uuid.Parse(id)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("ERROR: Could not parse feed follow id: %v", err))
		return
	}

	err = apiConf.DB.DeleteFeedFollows(r.Context(), database.DeleteFeedFollowsParams{
		ID:     feedFollowId,
		UserID: user.ID,
	})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("ERROR: Could not create feed follow: %v", err))
		return
	}

	respondWithJson(w, 200, struct{}{})
}
