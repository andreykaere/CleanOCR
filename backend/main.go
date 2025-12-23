package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/google/uuid"
)

type contextKey string

const SessionKey contextKey = "session_id"

func main() {
	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"http://0.0.0.0:8000"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "Set-Cookie"},
		ExposedHeaders:   []string{"Link", "Set-Cookie"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	r.With(AuthMiddleware).Post("/process", processFile)

	// r.Route("/jobs", func(r chi.Router) {
	// 	r.Get("/", listJobs)
	// 	r.Get("/{id}", getJob)
	// 	r.Delete("/{id}", cancelJob)
	// })

	http.ListenAndServe(":20080", r)
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var sessionID string
		cookie, err := r.Cookie("session_id")

		fmt.Println("quuz")

		if err == nil {
			sessionID = cookie.Value
		} else if err == http.ErrNoCookie {
			// TODO: check if UUID already exists, use DB for that or something
			// like that
			sessionID = uuid.NewString()
			cookie := &http.Cookie{
				Name:     "session_id",
				Value:    sessionID,
				Path:     "/",
				HttpOnly: true,
				// Secure:   true,
				// SameSite: http.SameSiteLaxMode,
				SameSite: http.SameSiteNoneMode,
			}
			http.SetCookie(w, cookie)
		} else {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		fmt.Println("foo")

		ctx := context.WithValue(r.Context(), SessionKey, sessionID)
		next.ServeHTTP(w, r.WithContext(ctx))

	})
}

func processFile(w http.ResponseWriter, r *http.Request) {
	sessionID, ok := r.Context().Value(SessionKey).(string)
	if !ok {
		http.Error(w, "user not found in context", http.StatusUnauthorized)
		return
	}
	// cookie, err := r.Cookie("session_id")

	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// sessionID := cookie.Value
	time.Sleep(1 * time.Second)
	fmt.Println("bar")

	w.Write([]byte(fmt.Sprintf("Hello from server, your cookie is %s", sessionID)))
	// http.ServeFile(w, r, )
}

// func saveFile(session_id string, file File) {
// 	detectedFileType := http.DetectContentType(fileBytes)
// 	switch detectedFileType {
// 	case "application/pdf":
// 		// TODO
// 	default:
// 		http.Error(w, "INVALID_FILE_TYPE", http.StatusBadRequest)
// 		return
// 	}
// }
