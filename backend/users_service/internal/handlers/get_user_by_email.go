package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"users_service/internal/database"

	"github.com/sirupsen/logrus"
)

func GetUserByEmail(db *sql.DB) http.HandlerFunc {
	// Инициализация логгера
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})

	return func(w http.ResponseWriter, r *http.Request) {
		email := r.URL.Query().Get("email")
		if email == "" {
			logger.Warn("Request missing 'email' parameter")
			http.Error(w, "Email is required", http.StatusBadRequest)
			return
		}

		logger.WithField("email", email).Info("Fetching user by email")

		user, err := database.GetUserByEmail(db, email)
		if err != nil {
			logger.WithError(err).Error("Database query failed")
			http.Error(w, "Failed to fetch user", http.StatusInternalServerError)
			return
		}

		if user == nil {
			logger.WithField("email", email).Warn("User not found")
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		logger.WithFields(logrus.Fields{
			"user_id":  user.ID,
			"username": user.Username,
			"email":    user.Email,
		}).Info("User found")

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(user); err != nil {
			logger.WithError(err).Error("Failed to encode response")
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
	}
}
