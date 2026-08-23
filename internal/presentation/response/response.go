package response

import (
	"encoding/json"
	"net/http"
	"project/internal/presentation/schemas"
)

// JSON - функция, необходимая для формирования HTTP ответа в формате JSON.
//
// В качестве параметров получает:
// - w, куда записывается ответ
// - status, HTTP статус, с которым возвращается ответ
// - payload, тело ответа, приходит ввиде любого тип данных
func JSON(w http.ResponseWriter, status int, payload any, options ...Option) {
	w.Header().Set("Content-Type", "application/json")
	for _, opt := range options {
		opt(w)
	}
	w.WriteHeader(status)

	buf, _ := json.Marshal(payload)
	w.Write(buf)
}

// Error - функция, необходимая для отправки ошибки в формате JSON по HTTP.
// Параметры получает те же, что и функция JSON
func Error(w http.ResponseWriter, status int, message schemas.ErrorSchema) {
	JSON(w, status, message)
}
