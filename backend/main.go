package main

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type contextKey string

const SessionKey contextKey = "X-Session-Id"

func main() {
	r := chi.NewRouter()

	r.With(AuthMiddleware).Post("/process", processFile)

	// r.Route("/jobs", func(r chi.Router) {
	// 	r.Get("/", listJobs)
	// 	r.Get("/{id}", getJob)
	// 	r.Delete("/{id}", cancelJob)
	// })

	http.ListenAndServe(":8080", r)
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionID := r.Header.Get("X-Session-Id")
		if sessionID == "" {
			sessionID := uuid.NewString()
			// TODO: check if UUID already exists, use DB for that or something
			// like that

			w.Header().Set("X-Session-Id", sessionID)
		}

		ctx := context.WithValue(r.Context(), SessionKey, sessionID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func processFile(w http.ResponseWriter, r *http.Request) {

}
