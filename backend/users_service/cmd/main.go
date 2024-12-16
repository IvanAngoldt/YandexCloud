package main

import (
	"log"
	"net/http"
	"os"

	"users_service/internal/database"
	"users_service/internal/handlers"

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
	db, err := database.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	r := mux.NewRouter()

	r.HandleFunc("/users/register", handlers.RegisterUser(db)).Methods("POST")

	r.HandleFunc("/users/by_email", handlers.GetUserByEmail(db)).Methods("GET")

	r.HandleFunc("/users/by_username", handlers.GetUserByUsername(db)).Methods("GET")

	r.HandleFunc("/users/{id:[0-9]+}", handlers.GetUserByID(db)).Methods("GET")
	r.HandleFunc("/users/{id:[0-9]+}", handlers.UpdateUser(db)).Methods("PATCH")
	r.HandleFunc("/users/{id:[0-9]+}/password", handlers.UpdateUserPassword(db)).Methods("PATCH")
	r.HandleFunc("/users/{id:[0-9]+}", handlers.DeleteUser(db)).Methods("DELETE")

	r.HandleFunc("/users", handlers.ListUsers(db)).Methods("GET")

	corsHandler := enableCORS(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8084"
	}
	log.Printf("Users Service running on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, corsHandler))
}
