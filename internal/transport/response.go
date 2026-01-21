package transport

import (
	"GoFrioCalor/internal/constants"
	"encoding/json"
	"net/http"
	"strings"
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

// GetHTTPStatusFromError mapea errores a códigos HTTP
func GetHTTPStatusFromError(err error) int {
	if err == nil {
		return http.StatusOK
	}
	errMsg := err.Error()
	// Errores 404 - Not Found
	if strings.Contains(errMsg, constants.ErrTermsSessionNotFound) ||
		strings.Contains(errMsg, "sesión de términos no encontrada") ||
		strings.Contains(errMsg, "sesión no encontrada") {
		return http.StatusNotFound
	}
	// Errores 410 - Gone (recurso expirado)
	if strings.Contains(errMsg, constants.MsgSessionExpired) ||
		strings.Contains(errMsg, "el token ha expirado") {
		return http.StatusGone
	}
	// Errores 400 - Bad Request
	if strings.Contains(errMsg, constants.ErrTermsNotAcceptedPending) ||
		strings.Contains(errMsg, constants.ErrTermsNotAcceptedRejected) ||
		strings.Contains(errMsg, constants.ErrTermsSessionExpired) ||
		strings.Contains(errMsg, "el token") {
		return http.StatusBadRequest
	}
	// Por defecto 500 - Internal Server Error
	return http.StatusInternalServerError
}
