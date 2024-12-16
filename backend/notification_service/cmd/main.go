package main

import (
	"log"
	"net/http"
	"notifications_service/internal/database"
	"notifications_service/internal/handlers"
	"os"

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

	// Создаем маршрутизатор
	r := mux.NewRouter()

	// Маршруты для уведомлений
	r.HandleFunc("/notifications", handlers.CreateNotification(db)).Methods("POST")
	r.HandleFunc("/notifications", handlers.FetchNotifications(db)).Methods("GET")
	r.HandleFunc("/notifications", handlers.DeleteNotification(db)).Methods("DELETE")
	r.HandleFunc("/notifications/read", handlers.MarkNotificationAsRead(db)).Methods("PATCH")
	r.HandleFunc("/notifications/{userId}/clear", handlers.ClearNotifications(db)).Methods("DELETE")

	// Применяем CORS middleware
	corsHandler := enableCORS(r)

	// Чтение порта из переменных окружения
	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}

	// Запуск сервера
	log.Printf("Notifications Service running on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, corsHandler))
}
