package main

import (
	"log"
	"net/http"
	"os"

	"speedkit-service/internal/handlers"
	"github.com/gorilla/mux"
)

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE, PATCH")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/recognize", handlers.Recognize()).Methods("POST")

	corsHandler := enableCORS(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8085"
	}

	log.Printf("Auth Service running on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, corsHandler))
}
