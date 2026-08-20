package main

import (
	"log"
	"net/http"

	"dash/db"
	"dash/routes/auth"
	"dash/routes/status"
	"dash/utils"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func main() {
	utils.LoadEnv()
	db.Connect()

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"http://localhost:5173"},
		AllowedMethods: []string{
			"GET",
			"POST",
			"PUT",
			"DELETE",
			"PATCH",
			"OPTIONS",
		},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	api := humachi.New(r, huma.DefaultConfig("dash api", "0.0.0"))

	auth.AuthRoutes(api)
	status.StatusRoutes(api)

	log.Println("HTTP server running at http://localhost:3333")
	http.ListenAndServe(":3333", r)
}
