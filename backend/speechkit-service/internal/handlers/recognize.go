package handlers

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/sirupsen/logrus"
)

func Recognize() http.HandlerFunc {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})

	// Получаем IAM-токен и Folder ID из переменных окружения
	iamToken := os.Getenv("YANDEX_API_KEY")
	folderID := os.Getenv("YANDEX_FOLDER_ID")
	if iamToken == "" || folderID == "" {
		logger.Fatal("YANDEX_API_KEY or YANDEX_FOLDER_ID are not set")
	}

	return func(w http.ResponseWriter, r *http.Request) {
		// Проверка Content-Type
		if r.Header.Get("Content-Type") != "audio/ogg" {
			http.Error(w, "Unsupported Content-Type, expected 'audio/ogg'", http.StatusBadRequest)
			return
		}

		// Чтение аудио данных из тела запроса
		audioData, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusInternalServerError)
			return
		}

		// Создание запроса к Yandex SpeechKit
		yandexURL := "https://stt.api.cloud.yandex.net/speech/v1/stt:recognize"
		req, err := http.NewRequest("POST", yandexURL, bytes.NewReader(audioData))
		if err != nil {
			http.Error(w, "Failed to create request to Yandex SpeechKit", http.StatusInternalServerError)
			return
		}

		// Установка заголовков
		req.Header.Set("Authorization", "Api-Key "+iamToken)
		req.Header.Set("Content-Type", "audio/ogg")

		// Добавление параметров запроса
		q := req.URL.Query()
		q.Add("folderId", folderID)
		q.Add("lang", "ru-RU")
		req.URL.RawQuery = q.Encode()

		// Отправка запроса к Yandex SpeechKit
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			http.Error(w, "Failed to send request to Yandex SpeechKit", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		// Чтение ответа от Yandex SpeechKit
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, "Failed to read response from Yandex SpeechKit", http.StatusInternalServerError)
			return
		}

		if resp.StatusCode != http.StatusOK {
			http.Error(w, "Yandex SpeechKit returned an error: "+string(body), http.StatusBadGateway)
			return
		}

		// Декодируем результат и возвращаем его
		var yandexResponse map[string]interface{}
		if err := json.Unmarshal(body, &yandexResponse); err != nil {
			http.Error(w, "Failed to decode response from Yandex SpeechKit", http.StatusInternalServerError)
			return
		}

		// Логирование успешного результата
		logger.WithField("recognized_text", yandexResponse["result"]).Info("Text recognized successfully and sent to the client")

		// Возвращаем результат распознавания
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(yandexResponse); err != nil {
			http.Error(w, "Failed to write response", http.StatusInternalServerError)
			return
		}
	}
}
