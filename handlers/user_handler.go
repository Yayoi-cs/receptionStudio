package handlers

import (
	"fmt"
	"net/http"
)

func UserHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello userHandler")
}
