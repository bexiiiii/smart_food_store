package handlers

import (
	"fmt"
	"net/http"
)

func CartHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Cart endpoint")
}
