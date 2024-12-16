package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

// RegisterRequest представляет данные для регистрации пользователя
type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// RegisterUser обрабатывает регистрацию пользователя
func RegisterUser() http.HandlerFunc {
	// Инициализация логгера
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(os.Stdout)

	return func(w http.ResponseWriter, r *http.Request) {
		var req RegisterRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			logger.WithError(err).Warn("Auth-Service: Invalid request payload")
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		// Очистка входных данных
		req.Username = strings.TrimSpace(req.Username)
		req.Email = strings.TrimSpace(req.Email)
		req.Password = strings.TrimSpace(req.Password)

		// Проверка обязательных полей
		if req.Username == "" || req.Email == "" || req.Password == "" {
			logger.Warn("Auth-Service: Missing required fields")
			http.Error(w, "All fields are required", http.StatusBadRequest)
			return
		}

		// Получение URL сервиса пользователей
		userServiceURL := os.Getenv("USERS_SERVICE_URL")
		if userServiceURL == "" {
			logger.Error("Auth-Service: Users service URL is not configured")
			http.Error(w, "Users service URL is not configured", http.StatusInternalServerError)
			return
		}

		// Подготовка запроса
		reqBody, err := json.Marshal(req)
		if err != nil {
			logger.WithError(err).Error("Auth-Service: Failed to encode request payload")
			http.Error(w, "Failed to encode request", http.StatusInternalServerError)
			return
		}

		// Выполнение запроса к сервису пользователей
		resp, err := http.Post(userServiceURL+"/users/register", "application/json", bytes.NewReader(reqBody))
		if err != nil {
			logger.WithFields(logrus.Fields{
				"service_url": userServiceURL,
				"error":       err.Error(),
			}).Error("Auth-Service: Failed to register user")
			http.Error(w, "Failed to register user", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		// Обработка ответа
		if resp.StatusCode != http.StatusCreated {
			logger.WithField("status_code", resp.StatusCode).Warn("Auth-Service: Failed to register user")

			var errResp map[string]interface{}
			if err := json.NewDecoder(resp.Body).Decode(&errResp); err == nil {
				http.Error(w, errResp["message"].(string), resp.StatusCode)
			} else {
				http.Error(w, "Failed to register user", resp.StatusCode)
			}
			return
		}

		logger.Info("Auth-Service: User successfully registered")

		w.WriteHeader(http.StatusCreated)
		var successResp map[string]string
		if err := json.NewDecoder(resp.Body).Decode(&successResp); err != nil {
			logger.WithError(err).Error("Auth-Service: Failed to parse response from users service")
			http.Error(w, "Failed to parse response from users service", http.StatusInternalServerError)
			return
		}

		if err := json.NewEncoder(w).Encode(successResp); err != nil {
			logger.WithError(err).Error("Auth-Service: Failed to send response to client")
		}
	}
}
