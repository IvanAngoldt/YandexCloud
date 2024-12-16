package main

import (
	"log"
	"net/http"
	"os"

	"posts_service/internal/database"
	"posts_service/internal/handlers"
	"posts_service/internal/middlewares"

	"github.com/gorilla/mux"
)

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE, PATCH")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	// Подключение к базе данных
	db, err := database.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	r := mux.NewRouter()

	r.Use(middlewares.AuthMiddleware)

	// Маршруты для постов
	r.HandleFunc("/posts", handlers.CreatePost(db)).Methods("POST")
	r.HandleFunc("/posts", handlers.FetchPosts(db)).Methods("GET")
	r.HandleFunc("/posts/{id}", handlers.FetchPostById(db)).Methods("GET")
	r.HandleFunc("/posts/{id}", handlers.DeletePost(db)).Methods("DELETE")

	// Маршруты для лайков
	r.HandleFunc("/likes", handlers.ToggleLike(db)).Methods("POST", "DELETE")
	r.HandleFunc("/likes", handlers.GetLikesForPost(db)).Methods("GET")

	// Маршрут для получения постов конкретного пользователя
	r.HandleFunc("/profile/{username}/posts", handlers.FetchUserPosts(db)).Methods("GET")

	corsHandler := enableCORS(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8083" // Порт для пост-сервиса
	}

	log.Printf("Posts Service running on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, corsHandler))
}
