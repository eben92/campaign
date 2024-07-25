package server

import (
	"campaign/internal/handlers/auth"
	"campaign/internal/utils/jwt"
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.CleanPath)
	r.Use(middleware.AllowContentType("application/json", "application/x-www-form-urlencoded", "multipart/form-data", "text/plain", "text/html"))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"https://*", "http://*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "Accept-Encoding"},
	}))

	r.Use(func(next http.Handler) http.Handler {
		hfn := func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")

			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(hfn)
	})

	r.Route("/api", func(api chi.Router) {
		api.Get("/health", s.healthHandler)

		api.Route("/", s.authController)

		api.Group(func(sapi chi.Router) {
			sapi.Use(jwt.Authenticator())
			sapi.Get("/", s.HelloWorldHandler)

		})

	})

	r.Get("/", s.HelloWorldHandler)
	return r
}

func (s *Server) authController(r chi.Router) {
	client := s.db.Database()
	handler := auth.NewAuthHandler(client)

	r.Post("/signin", handler.Signin)
	r.Post("/create-account", handler.Signup)
}

func (s *Server) HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("error handling JSON marshal. Err: %v", err)
	}

	_, _ = w.Write(jsonResp)
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	jsonResp, _ := json.Marshal(s.db.Health())
	_, _ = w.Write(jsonResp)
}
