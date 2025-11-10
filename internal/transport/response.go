package transport

import (
	"encoding/json"
	"net/http"
)

// RespondWithJSON escribe una respuesta JSON con el código de estado y datos proporcionados
func RespondWithJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// RespondWithError escribe una respuesta de error en formato JSON
func RespondWithError(w http.ResponseWriter, statusCode int, message string) {
	RespondWithJSON(w, statusCode, map[string]string{"error": message})
}
