package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

// LoginRequest представляет данные для запроса входа
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Login обрабатывает вход пользователя
func Login() http.HandlerFunc {
    logger := logrus.New()
    logger.SetFormatter(&logrus.JSONFormatter{})
    logger.SetOutput(os.Stdout)

	logger.Printf("Auth-Service: Request to login")

    return func(w http.ResponseWriter, r *http.Request) {
        logger.Info("Auth-Service: Login request received")

        var req LoginRequest
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            logger.WithError(err).Warn("Auth-Service: Invalid request payload")
            http.Error(w, "Invalid request payload", http.StatusBadRequest)
            return
        }

        if req.Email == "" || req.Password == "" {
            logger.Warn("Auth-Service: Email and password are required")
            http.Error(w, "Email and password are required", http.StatusBadRequest)
            return
        }

        userServiceURL := os.Getenv("USERS_SERVICE_URL")
        if userServiceURL == "" {
            logger.Error("Auth-Service: Users service URL is not configured")
            http.Error(w, "Users service URL is not configured", http.StatusInternalServerError)
            return
        }

        resp, err := http.Get(fmt.Sprintf("%s/users/by_email?email=%s", userServiceURL, req.Email))
        if err != nil {
            logger.WithFields(logrus.Fields{
                "email":       req.Email,
                "service_url": userServiceURL,
                "error":       err.Error(),
            }).Error("Auth-Service: Error during user fetch")
            http.Error(w, "Invalid email or password", http.StatusForbidden)
            return
        }
        defer resp.Body.Close()

        if resp.StatusCode == http.StatusNotFound {
            logger.WithField("email", req.Email).Warn("Auth-Service: User not found")
            http.Error(w, "Invalid email or password", http.StatusForbidden)
            return
        }

        if resp.StatusCode != http.StatusOK {
            logger.WithField("status_code", resp.StatusCode).Warn("Auth-Service: Unexpected response from user service")
            http.Error(w, "Invalid email or password", http.StatusForbidden)
            return
        }

        var user struct {
            ID           int    `json:"id"`
            Username     string `json:"username"`
            Email        string `json:"email"`
            PasswordHash string `json:"password_hash"`
        }
        if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
            logger.WithError(err).Error("Auth-Service: Failed to parse user data")
            http.Error(w, "Failed to parse user data", http.StatusInternalServerError)
            return
        }

        if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
            logger.WithField("email", req.Email).Warn("Auth-Service: Invalid email or password")
            http.Error(w, "Invalid email or password", http.StatusForbidden)
            return
        }

        token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
            "user_id": user.ID,
            "email":   user.Email,
            "exp":     time.Now().Add(72 * time.Hour).Unix(),
        })

        secret := os.Getenv("JWT_SECRET")
        if secret == "" {
            logger.Error("Auth-Service: JWT_SECRET is not configured")
            http.Error(w, "JWT_SECRET is not configured", http.StatusInternalServerError)
            return
        }

        tokenStr, err := token.SignedString([]byte(secret))
        if err != nil {
            logger.WithError(err).Error("Auth-Service: Failed to generate token")
            http.Error(w, "Failed to generate token", http.StatusInternalServerError)
            return
        }

        logger.WithFields(logrus.Fields{
            "user_id":   user.ID,
            "username":  user.Username,
            "email":     user.Email,
            "token_exp": time.Now().Add(72 * time.Hour).Format(time.RFC3339),
        }).Info("Auth-Service: Login successful")

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(map[string]interface{}{
            "message": "Login successful",
            "user": map[string]interface{}{
                "id":       user.ID,
                "username": user.Username,
                "email":    user.Email,
            },
            "token": tokenStr,
        })
    }
}
