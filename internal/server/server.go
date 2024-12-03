package server

import (
	"log"
	"net/http"

	"github.com/cyberkillua/dailyread/internal/config"
	"github.com/cyberkillua/dailyread/internal/database"
	"github.com/cyberkillua/dailyread/internal/handlers"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
)

type Server struct {
	config *config.Config
	db     *database.Queries
	router *chi.Mux
}

func New(cfg *config.Config, db *database.Queries) *Server {
	router := chi.NewRouter()

	// CORS middleware
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	srv := &Server{
		config: cfg,
		db:     db,
		router: router,
	}

	srv.setupRoutes()
	return srv
}

func (s *Server) setupRoutes() {
	v1Router := chi.NewRouter()

	apiConfig := &handlers.APIConfig{DB: s.db}

	v1Router.Get("/healthz", handlers.HandlerReadiness)
	v1Router.Get("/err", handlers.HandlerErr)
	v1Router.Post("/webpages", apiConfig.CreateWebpage)
	// v1Router.Post("/users", apiConfig.handlerCreateUser)
	// v1Router.Get("/users", apiConfig.middlewareAuth(apiConfig.handlerGetUserByAPIKey))
	// v1Router.Post("/feeds", apiConfig.middlewareAuth(apiConfig.handlerCreateFeed))
	// v1Router.Get("/feeds", apiConfig.getFeeds)
	// v1Router.Post("/follow", apiConfig.middlewareAuth(apiConfig.handlerFeedFollow))

	s.router.Mount("/v1", v1Router)
}

func (s *Server) Start() error {
	srv := &http.Server{
		Handler: s.router,
		Addr:    ":" + s.config.Port,
	}

	log.Printf("Server listening on port %v", s.config.Port)
	return srv.ListenAndServe()
}
