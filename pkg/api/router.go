package api

import (
	"encoding/json"
	"net/http"
	"strings"
)

// APIServer manages routing and security middlewares for the API endpoints.
type APIServer struct {
	SecretToken string
}

// Response is a generic API JSON response helper.
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// AuthMiddleware verifies the custom authorization token in the request headers.
func (s *APIServer) AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			s.writeJSON(w, http.StatusUnauthorized, Response{Success: false, Message: "Missing auth header"})
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token != s.SecretToken {
			s.writeJSON(w, http.StatusForbidden, Response{Success: false, Message: "Invalid authentication token"})
			return
		}

		next(w, r)
	}
}

// SetupRoutes registers handlers for endpoints in a standard pattern.
func (s *APIServer) SetupRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/v1/status", s.AuthMiddleware(s.handleStatus))
}

func (s *APIServer) handleStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		s.writeJSON(w, http.StatusMethodNotAllowed, Response{Success: false, Message: "Method not allowed"})
		return
	}

	data := map[string]string{
		"status":  "healthy",
		"version": "1.0.0-preview",
	}

	s.writeJSON(w, http.StatusOK, Response{Success: true, Data: data})
}

func (s *APIServer) writeJSON(w http.ResponseWriter, status int, val Response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(val)
}
