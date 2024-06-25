package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"watermelon_king-api/handlers"
	"watermelon_king-api/pkg/turso"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

func main() {
	Init()
	r := chi.NewRouter()
	port := os.Getenv("PORT")
	r.Use(middleware.Logger)
	r.Post("/", handlers.NewHandlerFunc(handlers.CreateUserHandler))
	r.Get("/", handlers.NewHandlerFunc(handlers.GetUserHandler))

	fmt.Println("Server listening on port", port)
	http.ListenAndServe(port, r)
}

func Init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
	turso.Init()
}
