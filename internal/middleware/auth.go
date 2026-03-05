package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// Cliente HTTP reutilizable con timeout configurado
var validationHTTPClient = &http.Client{
	Timeout: 5 * time.Second,
}

// TokenValidationResponse representa la respuesta del servicio de validación
type TokenValidationResponse struct {
	Valid  bool   `json:"valido"`
	Detail string `json:"detail,omitempty"`
}

// AuthMiddleware valida el token JWT contra el servicio externo de autenticación
func AuthMiddleware(authServiceURL string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extraer el token del header Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			log.Warn().Msg(LogRequestWithoutAuth)
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":  ErrTokenNotProvided,
				"detail": ErrTokenRequiredDetail,
			})
			c.Abort()
			return
		}

		// Verificar que sea un Bearer token
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			log.Warn().Str("authHeader", authHeader).Msg(LogInvalidAuthFormat)
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":  ErrInvalidTokenFormat,
				"detail": ErrInvalidFormatDetail,
			})
			c.Abort()
			return
		}

		token := parts[1]

		// Validar el token contra el servicio externo
		isValid, detail := validateToken(c.Request.Context(), authServiceURL, token)
		if !isValid {
			log.Warn().
				Str("detail", detail).
				Str("ip", c.ClientIP()).
				Msg(LogTokenInvalidOrExpired)
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":  ErrInvalidToken,
				"detail": detail,
			})
			c.Abort()
			return
		}

		log.Debug().
			Str("ip", c.ClientIP()).
			Str("path", c.Request.URL.Path).
			Msg(LogTokenValidated)

		// Token válido, continuar con la request
		c.Next()
	}
}

// validateToken hace una llamada HTTP al servicio de validación de tokens
func validateToken(ctx context.Context, authServiceURL, token string) (bool, string) {
	// Construir la URL del endpoint de validación
	validationURL := fmt.Sprintf("%s/validar-token", authServiceURL)

	// Crear la request con context para propagación de cancelación
	req, err := http.NewRequestWithContext(ctx, "GET", validationURL, nil)
	if err != nil {
		log.Error().Err(err).Msg(LogErrorCreatingRequest)
		return false, ErrInternalValidation
	}

	// Agregar el Bearer token al header
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	// Hacer la request con el cliente HTTP reutilizable
	resp, err := validationHTTPClient.Do(req)
	if err != nil {
		log.Error().Err(err).Str("url", validationURL).Msg(LogErrorCallingService)
		return false, ErrAuthServiceUnavailable
	}
	defer resp.Body.Close()

	// Leer la respuesta
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error().Err(err).Msg(LogErrorReadingResponse)
		return false, ErrProcessingResponse
	}

	// Si el status code no es 200, el token es inválido
	if resp.StatusCode != http.StatusOK {
		log.Warn().
			Int("statusCode", resp.StatusCode).
			Str("response", string(body)).
			Msg(LogServiceReturnedError)
		return false, ErrTokenExpired
	}

	// Parsear la respuesta JSON
	var validationResp TokenValidationResponse
	if err := json.Unmarshal(body, &validationResp); err != nil {
		log.Error().Err(err).Str("body", string(body)).Msg(LogErrorParsingResponse)
		return false, ErrProcessingResponse
	}

	// Retornar el resultado
	if !validationResp.Valid {
		return false, validationResp.Detail
	}

	return true, ""
}
