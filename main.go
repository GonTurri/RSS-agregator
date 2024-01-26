package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/GonTurri/RSS-agregator/internal/database"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {
	godotenv.Load()
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("port not found")
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("dbURL not found")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("error when connecting to databse: ", err)
	}
	dbQueries := database.New(db)

	apCfg := apiConfig{DB: dbQueries}

	go startScraping(dbQueries, 10, time.Minute)

	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))
	v1Router := chi.NewRouter()

	//GET
	v1Router.Get("/healthz", handlerReadiness)
	v1Router.Get("/err", handlerErr)
	v1Router.Get("/users", apCfg.middlewareAuth(apCfg.getUserHnadler))
	v1Router.Get("/feeds", apCfg.getFeedsHandler)
	v1Router.Get("/feed_follows", apCfg.middlewareAuth(apCfg.getFeedFollowsForUser))

	//POST

	v1Router.Post("/users", apCfg.createUserHandler)
	v1Router.Post("/feeds", apCfg.middlewareAuth(apCfg.createFeedHandler))
	v1Router.Post("/feed_follows", apCfg.middlewareAuth(apCfg.createFeedFollowHandler))

	// DELETE

	v1Router.Delete("/feed_follows/{feedFollowID}", apCfg.middlewareAuth(apCfg.deleteFeedFollowHandler))

	router.Mount("/v1", v1Router)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	log.Printf("server running on port: %s", port)
	log.Fatal(server.ListenAndServe())

}
