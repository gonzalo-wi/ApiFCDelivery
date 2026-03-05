package transport

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// Cliente HTTP reutilizable para llamadas al servicio de auth
var authHTTPClient = &http.Client{
	Timeout: 10 * time.Second,
}

type AuthHandler struct {
	authServiceURL string
}

func NewAuthHandler(authServiceURL string) *AuthHandler {
	return &AuthHandler{
		authServiceURL: authServiceURL,
	}
}

// GenerateToken hace de proxy hacia el servicio de autenticación interno
func (h *AuthHandler) GenerateToken(c *gin.Context) {
	// Obtener el x-api-key del header
	apiKey := c.GetHeader("x-api-key")
	if apiKey == "" {
		log.Warn().Msg(LogRequestWithoutAPIKey)
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":  ErrAPIKeyNotProvided,
			"detail": ErrAPIKeyRequiredDetail,
		})
		return
	}

	url := fmt.Sprintf("%s/generar-token", h.authServiceURL)

	// Crear la request con context para propagación de cancelación
	req, err := http.NewRequestWithContext(c.Request.Context(), "GET", url, nil)
	if err != nil {
		log.Error().Err(err).Msg(LogErrorCreatingAuthRequest)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": ErrInternalGeneratingToken,
		})
		return
	}

	// Pasar el x-api-key al servicio interno
	req.Header.Set("x-api-key", apiKey)

	// Ejecutar la request con el cliente HTTP reutilizable
	resp, err := authHTTPClient.Do(req)
	if err != nil {
		log.Error().Err(err).Str("url", url).Msg(LogErrorCallingAuthService)
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error":  ErrAuthServiceUnavailable,
			"detail": ErrAuthServiceUnavailableDetail,
		})
		return
	}
	defer resp.Body.Close()

	// Leer la respuesta
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error().Err(err).Msg(LogErrorReadingAuthResponse)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": ErrProcessingAuthResponse,
		})
		return
	}

	// Si el servicio de auth retorna error, pasar ese error al cliente
	if resp.StatusCode != http.StatusOK {
		log.Warn().
			Int("statusCode", resp.StatusCode).
			Str("response", string(body)).
			Str("apiKey", maskAPIKey(apiKey)).
			Msg(LogAuthServiceReturnedError)

		var errorResponse map[string]interface{}
		if err := json.Unmarshal(body, &errorResponse); err == nil {
			c.JSON(resp.StatusCode, errorResponse)
		} else {
			c.JSON(resp.StatusCode, gin.H{
				"error":  ErrGeneratingToken,
				"detail": string(body),
			})
		}
		return
	}

	// Parsear y retornar la respuesta exitosa
	var tokenResponse map[string]interface{}
	if err := json.Unmarshal(body, &tokenResponse); err != nil {
		log.Error().Err(err).Str("body", string(body)).Msg(LogErrorParsingAuthResponse)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": ErrProcessingGeneratedToken,
		})
		return
	}

	log.Info().
		Str("proveedor", fmt.Sprintf("%v", tokenResponse["proveedor"])).
		Str("apiKey", maskAPIKey(apiKey)).
		Str("ip", c.ClientIP()).
		Msg(LogTokenGeneratedSuccess)

	c.JSON(http.StatusOK, tokenResponse)
}

// maskAPIKey enmascara la API key para los logs (protección de datos sensibles)
func maskAPIKey(apiKey string) string {
	if len(apiKey) <= 8 {
		return "***"
	}
	return apiKey[:4] + "..." + apiKey[len(apiKey)-4:]
}
