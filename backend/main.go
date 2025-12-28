package main

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type contextKey string

const SessionKey contextKey = "session_id"

func main() {
	r := chi.NewRouter()

	r.With(AuthMiddleware).Post("/process", processFile)

	// r.Route("/jobs", func(r chi.Router) {
	// 	r.Get("/", listJobs)
	// 	r.Get("/{id}", getJob)
	// 	r.Delete("/{id}", cancelJob)
	// })
	certFile := "/etc/letsencrypt/live/rozetka.hopto.org/fullchain.pem"
	keyFile := "/etc/letsencrypt/live/rozetka.hopto.org/privkey.pem"

	http.ListenAndServeTLS(":5000", certFile, keyFile, r)
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
				Secure:   true,
				SameSite: http.SameSiteLaxMode,
				//SameSite: http.SameSiteNoneMode,
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
	w.Write([]byte(fmt.Sprintf("Hello from server, your cookie is %s", sessionID)))

	file, fh, err := r.FormFile("file")
	if err != nil {
		// TODO: handle this
		panic(err)
	}
	defer file.Close()

	saveFile(sessionID, file, fh)
}

func saveFile(sessionID string, file multipart.File, fh *multipart.FileHeader) error {
	file, err := fh.Open()
	if err != nil {
		return err
	}
	defer file.Close()

	if err = os.MkdirAll("/files", 0755); err != nil {
		return err
	}

	dst, err := os.Create("/files/" + fh.Filename)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)

	return err
}

