package main

// TODO: Add table names as constants

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"elo-rating/controllers"
	"elo-rating/database"

	"github.com/gorilla/mux"
)

type MiddlewareFunc func(http.Handler) http.Handler

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		levelStr := os.Getenv("LOG_LEVEL")
		if levelStr == "debug" {
			buf, err := io.ReadAll(r.Body)
			if err != nil {
				log.Printf("[ERROR] Error reading request body: %v", err.Error())
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if string(buf) == "" {
				log.Printf("[DEBUG] %v %v", r.Method, r.RequestURI)
			} else {
				log.Printf("[DEBUG] %v %v | Request body: %v", r.Method, r.RequestURI, string(buf))
			}
			reader := io.NopCloser(bytes.NewBuffer(buf))
			r.Body = reader
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	db := database.InitDatabase("./app.db")
	db.Migrate()
	defer database.Connection.Close()

	router := mux.NewRouter()
	controllers.SetRoutes(router)
	router.Use(loggingMiddleware)

	// Start the server
	log.Println("[INFO] Starting server on :8080")

	srv := &http.Server{
		Handler: router,
		Addr:    "127.0.0.1:8080",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
