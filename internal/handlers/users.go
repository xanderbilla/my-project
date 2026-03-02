// Package handlers contains HTTP route handlers.
package handlers

import "net/http"

// RegisterRoutes registers all application routes.
func GetAllUsers(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/users", usersHandler)
}

func usersHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Welcome!"))
}