package middleware

// Constantes para mensajes de error de autenticación
const (
	ErrTokenNotProvided       = "Token no proporcionado"
	ErrTokenRequiredDetail    = "Se requiere el header Authorization con un Bearer token"
	ErrInvalidTokenFormat     = "Formato de token inválido"
	ErrInvalidFormatDetail    = "El header Authorization debe tener el formato: Bearer <token>"
	ErrInvalidToken           = "Token inválido"
	ErrTokenExpired           = "Token inválido o expirado"
	ErrAuthServiceUnavailable = "Servicio de autenticación no disponible"
	ErrInternalValidation     = "Error interno al validar token"
	ErrProcessingResponse     = "Error procesando respuesta de validación"
)

// Mensajes de log
const (
	LogRequestWithoutAuth    = "Request sin header Authorization"
	LogInvalidAuthFormat     = "Formato de Authorization inválido"
	LogTokenInvalidOrExpired = "Token inválido o expirado"
	LogTokenValidated        = "Token validado exitosamente"
	LogErrorCreatingRequest  = "Error creando request de validación"
	LogErrorCallingService   = "Error llamando al servicio de validación"
	LogErrorReadingResponse  = "Error leyendo respuesta de validación"
	LogServiceReturnedError  = "Servicio de validación retornó error"
	LogErrorParsingResponse  = "Error parseando respuesta de validación"
)
