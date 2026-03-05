package transport

// Constantes para mensajes de error del handler de autenticación
const (
	ErrAPIKeyNotProvided            = "API Key no proporcionada"
	ErrAPIKeyRequiredDetail         = "Se requiere el header x-api-key"
	ErrInternalGeneratingToken      = "Error interno al generar token"
	ErrAuthServiceUnavailable       = "Servicio de autenticación no disponible"
	ErrAuthServiceUnavailableDetail = "No se pudo conectar con el servicio de tokens"
	ErrProcessingAuthResponse       = "Error procesando respuesta del servicio de auth"
	ErrGeneratingToken              = "Error generando token"
	ErrProcessingGeneratedToken     = "Error procesando token generado"
)

// Mensajes de log del auth handler
const (
	LogRequestWithoutAPIKey     = "Request sin x-api-key header"
	LogErrorCreatingAuthRequest = "Error creando request al servicio de auth"
	LogErrorCallingAuthService  = "Error llamando al servicio de auth"
	LogErrorReadingAuthResponse = "Error leyendo respuesta del servicio de auth"
	LogAuthServiceReturnedError = "Servicio de auth retornó error"
	LogErrorParsingAuthResponse = "Error parseando respuesta del servicio de auth"
	LogTokenGeneratedSuccess    = "Token generado exitosamente"
)
