package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
)

func fetchUsernameFromUsersService(userID int) (string, error) {
	userServiceURL := os.Getenv("USERS_SERVICE_URL")
	if userServiceURL == "" {
		return "", fmt.Errorf("USERS_SERVICE_URL not set")
	}

	url := fmt.Sprintf("%s/users/%d", userServiceURL, userID)
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to fetch user: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch user: status %d", resp.StatusCode)
	}

	var user struct {
		ID       int    `json:"id"`
		Username string `json:"username"`
		Email    string `json:"email"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return "", fmt.Errorf("failed to decode user data: %w", err)
	}

	return user.Username, nil
}

// fetchUsersFromUsersService принимает список userID и возвращает список мап с данными пользователей (id, username)
func fetchUsersFromUsersService(userIDs []int) ([]map[string]interface{}, error) {
	users := make([]map[string]interface{}, 0, len(userIDs))
	for _, uid := range userIDs {
		username, err := fetchUsernameFromUsersService(uid)
		if err != nil {
			return nil, err
		}
		users = append(users, map[string]interface{}{
			"id":       uid,
			"username": username,
		})
	}
	return users, nil
}

// helper для конвертации string->int с обработкой ошибки
func atoiParam(param string) (int, error) {
	return strconv.Atoi(param)
}
