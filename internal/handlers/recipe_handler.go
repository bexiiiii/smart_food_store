package handlers

import (
	"fmt"
	"net/http"
)

func RecipeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Recipe endpoint")
}
