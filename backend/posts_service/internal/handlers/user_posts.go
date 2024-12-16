package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"os"
	"posts_service/internal/database"
	"strconv"

	"github.com/gorilla/mux"
)

func FetchUserPosts(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		username := vars["username"]

		// Получаем userID по username через Users Service
		userID, err := fetchUserIDByUsername(username)
		if err != nil {
			http.Error(w, "Failed to find user by username", http.StatusNotFound)
			return
		}

		// Получаем посты пользователя
		posts, err := database.FetchUserPosts(db, userID)
		if err != nil {
			http.Error(w, "Failed to fetch posts", http.StatusInternalServerError)
			return
		}

		// Отправляем ответ
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(posts)
	}
}

// fetchUserIDByUsername запрашивает userID из Users Service по username.
func fetchUserIDByUsername(username string) (int, error) {
	userServiceURL := os.Getenv("USERS_SERVICE_URL")
	if userServiceURL == "" {
		return 0, errNoUserService
	}

	url := userServiceURL + "/users/by_username?username=" + username
	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, errUserNotFound
	}

	var user struct {
		ID int `json:"id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return 0, err
	}
	return user.ID, nil
}

// Ошибки
var (
	errNoUserService = strconv.ErrSyntax
	errUserNotFound  = strconv.ErrRange
)
