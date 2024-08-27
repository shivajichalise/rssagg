package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"github.com/shivajichalise/rssagg/internal/database"

	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("ERROR: Error loading .env file")
	}

	PORT := os.Getenv("PORT")
	if PORT == "" {
		log.Fatal("ERROR: Port is not defined")
	}

	DB_URL := os.Getenv("DB_URL")
	if DB_URL == "" {
		log.Fatal("ERROR: DB URL is not defined")
	}

	conn, err := sql.Open("postgres", DB_URL)
	if err != nil {
		log.Fatal("ERROR: Could not connect to database")
	}

	queries := database.New(conn)
	apiConf := apiConfig{
		DB: queries,
	}

	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	v1Router := chi.NewRouter()
	v1Router.Get("/healthz", handlerReadiness)
	v1Router.Get("/err", handlerError)
	v1Router.Post("/users", apiConf.handlerCreateUser)
	v1Router.Get("/users", apiConf.middlewareAuth(apiConf.handlerGetUserByApiKey))
	v1Router.Post("/feeds", apiConf.middlewareAuth(apiConf.handlerCreateFeed))
	v1Router.Get("/feeds", apiConf.handlerGetFeeds)
	v1Router.Post("/feed-follows", apiConf.middlewareAuth(apiConf.handlerCreateFeedFollow))
	v1Router.Get("/feed-follows", apiConf.middlewareAuth(apiConf.handlerGetFeedFollow))
	v1Router.Delete("/feed-follows/{id}", apiConf.middlewareAuth(apiConf.handlerDeleteFeedFollow))

	router.Mount("/v1", v1Router)

	server := &http.Server{
		Handler: router,
		Addr:    ":" + PORT,
	}

	fmt.Println("INFO: Server is starting on port: ", PORT)
	server_err := server.ListenAndServe()
	if server_err != nil {
		log.Fatal(err)
	}
}
