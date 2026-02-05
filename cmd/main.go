package main

import (
	"log"
	"net/http"

	"github.com/bexiiiii/smart_food_store/internal/handlers"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/users", handlers.UserHandler)
	mux.HandleFunc("/recipes", handlers.RecipeHandler)
	mux.HandleFunc("/carts", handlers.CartHandler)

	log.Println("Server started on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
