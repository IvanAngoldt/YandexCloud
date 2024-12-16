package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"

	"users_service/internal/database"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

// RegisterRequest представляет данные для запроса регистрации
type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// RegisterUser обрабатывает регистрацию нового пользователя
func RegisterUser(db *sql.DB) http.HandlerFunc {
	// Инициализация логгера
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})

	return func(w http.ResponseWriter, r *http.Request) {
		logger.Info("Users-Service: Registration request received")

		var req RegisterRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			logger.WithError(err).Warn("Users-Service: Failed to decode request payload")
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		// Очистка данных
		req.Username = strings.TrimSpace(req.Username)
		req.Email = strings.TrimSpace(req.Email)
		req.Password = strings.TrimSpace(req.Password)

		// Валидация данных
		if req.Username == "" || req.Email == "" || req.Password == "" {
			logger.Warn("Users-Service: Validation error - all fields are required")
			http.Error(w, "All fields are required", http.StatusBadRequest)
			return
		}

		// Хэширование пароля
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			logger.WithError(err).Error("Users-Service: Failed to hash password")
			http.Error(w, "Failed to hash password", http.StatusInternalServerError)
			return
		}

		// Сохранение пользователя в базе данных
		err = database.SaveUser(db, req.Username, req.Email, string(hashedPassword))
		if err != nil {
			logger.WithFields(logrus.Fields{
				"username": req.Username,
				"email":    req.Email,
				"error":    err.Error(),
			}).Error("Users-Service: Failed to register user in the database")
			http.Error(w, "Failed to register user", http.StatusInternalServerError)
			return
		}

		// Успешный ответ
		logger.WithFields(logrus.Fields{
			"username": req.Username,
			"email":    req.Email,
		}).Info("Users-Service: User successfully registered")

		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(map[string]string{"message": "Registration successful"}); err != nil {
			logger.WithError(err).Error("Users-Service: Failed to send response to client")
		}
	}
}
