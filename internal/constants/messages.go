package constants

const (
	MsgDeliveryNotFound        = "Entrega no encontrada"
	MsgDeliveryCreated         = "Entrega creada exitosamente"
	MsgDeliveryUpdated         = "Entrega actualizada exitosamente"
	MsgDeliveryDeleted         = "Entrega eliminada exitosamente"
	MsgDispenserCreated        = "Dispenser creado exitosamente"
	MsgDispenserUpdated        = "Dispenser actualizado exitosamente"
	MsgDispenserDeleted        = "Dispenser eliminado exitosamente"
	MsgInternalServerError     = "Error interno del servidor"
	MsgInvalidInput            = "Entrada inválida"
	MsgDatabaseConnectionError = "Error de conexión a la base de datos"
	MsgDatabaseConnected       = "Conexión a la base de datos exitosa"
	MsgInvalidID               = "Número de ID inválido"
	MsgPDFGenerationError      = "Error generando PDF"

	// Mensajes de sistema
	MsgErrorLoadingConfig  = "Error cargando configuración"
	MsgErrorConnectingDB   = "Error conectando a la base de datos"
	MsgDBConnectedSuccess  = "Base de datos conectada y tablas migradas correctamente"
	MsgServerRunning       = "Servidor corriendo en http://localhost:%s"
	MsgErrorStartingServer = "Error iniciando el servidor"

	MsgTextAcepted = "Declaro haber leído y aceptado integramente el Contrato de Alquiler de Equipo Frío Calor (el \"Contrato\") que se encuentra disponible en el siguiente link: www.somoselagua.com.ar/tycfriocalor y que ha sido protocolizado por el escribano Juan Franciso Iribarren , titular del Registro Notarial No 60 de La Matanza mediante escritura No 17 Folio 52 de fecha 12 de Febrero de 2021, conforme he comprobado en el siguiente link: www.somoselagua.com.ar/certfriocalor. Con arreglo a lo previsto en el Art. 4 de la ley 24.240, acepto que los términos y condiciones del Contrato y de uso del Equipo me sean suministrados por el medio antes descripto, en reemplazo del soporte físico. Usted tiene derecho a revocar la aceptación de este contrato dentro de los diez días computados a partir de la fecha de la presente Orden de Trabajo."

	// Textos del PDF Orden de Trabajo
	PDFHeaderTitle         = "SERVICIO TECNICO"
	PDFOrderTitle          = "ORDEN DE TRABAJO"
	PDFSectionService      = "INFORMACION DEL SERVICIO"
	PDFSectionClient       = "DATOS DEL CLIENTE"
	PDFSectionEquipment    = "EQUIPOS INSTALADOS"
	PDFSectionTask         = "TAREA REALIZADA"
	PDFSectionTerms        = "TERMINOS Y CONDICIONES"
	PDFLabelDate           = "Fecha:"
	PDFLabelActionType     = "Tipo de Accion:"
	PDFLabelAccountNumber  = "Nro. Cuenta:"
	PDFLabelDeliveryNumber = "Nro. Reparto:"
	PDFLabelName           = "Nombre:"
	PDFLabelAddress        = "Direccion:"
	PDFLabelLocality       = "Localidad:"
	PDFTableItem           = "Item"
	PDFTableBrand          = "Marca"
	PDFTableSerialNumber   = "Numero de Serie"
	PDFSectionAcceptance   = "ACEPTACION DIGITAL"
	PDFLabelAccepted       = "Estado:"
	PDFLabelAcceptedValue  = "ACEPTADO DIGITALMENTE"
	PDFLabelDateTime       = "Fecha y Hora:"
	PDFLabelToken          = "Token de Verificacion:"
	PDFFooterImportant     = "IMPORTANTE: No realizar la devolucion del equipo sin su correspondiente comprobante, el cual es entregado en el momento por nuestro representante."
	PDFAcceptanceNote      = "El cliente fue informado sobre los terminos y condiciones del servicio y acepto digitalmente mediante el token de verificacion."

	// Descripciones de tareas
	TaskInstallation = "Se realizo la instalacion del Dispenser Frio Calor"
	TaskRemoval      = "Se realizo el retiro del Dispenser Frio Calor"
	TaskReplacement  = "Se realizo el recambio del Dispenser Frio Calor"

	// Mensajes de validación
	ValidationRequired    = "%s es requerido"
	ValidationMinLength   = "%s debe tener al menos %s caracteres"
	ValidationMaxLength   = "%s debe tener máximo %s caracteres"
	ValidationExactLength = "%s debe tener exactamente %s caracteres"
	ValidationOneOf       = "%s debe ser uno de: %s"
	ValidationNumeric     = "%s debe ser numérico"
	ValidationGreaterThan = "%s debe ser mayor que %s"
	ValidationInvalid     = "%s no es válido"

	// Store error messages
	ErrFindAllDeliveries        = "error al buscar todas las entregas: %w"
	ErrFindDeliveryByID         = "error al buscar entrega con id %d: %w"
	ErrFindDeliveriesFilters    = "error al buscar entregas con filtros: %w"
	ErrFindDeliveriesByRto      = "error al buscar entregas por RTO: %w"
	ErrCreateDelivery           = "error al crear entrega: %w"
	ErrUpdateDelivery           = "error al actualizar entrega: %w"
	ErrDeleteDelivery           = "error al eliminar entrega con id %d: %w"
	ErrFindAllDispensers        = "error al buscar todos los dispensers: %w"
	ErrFindDispenserByID        = "error al buscar dispenser con id %d: %w"
	ErrFindDispensersByDelivery = "error al buscar dispensers de la entrega %d: %w"
	ErrCreateDispenser          = "error al crear dispenser: %w"
	ErrUpdateDispenser          = "error al actualizar dispenser: %w"
	ErrDeleteDispenser          = "error al eliminar dispenser con id %d: %w"
	ErrCreateWorkOrder          = "error al crear orden de trabajo: %w"
	ErrCountWorkOrders          = "error al contar órdenes de trabajo: %w"

	// Terms Session Messages
	MsgTermsAlreadyAccepted = "Términos ya fueron aceptados previamente"
	MsgTermsAcceptedSuccess = "Términos aceptados exitosamente"
	MsgTermsAlreadyRejected = "Términos ya fueron rechazados previamente"
	MsgTermsRejected        = "Términos rechazados"
	MsgSessionExpired       = "el token ha expirado"
	MsgSessionNotAvailable  = "el token no está disponible para esta acción (estado: %s)"

	// Terms Session Errors
	ErrVerifyingExistingSession = "error verificando sesión existente: %w"
	ErrGeneratingToken          = "error generando token: %w"
	ErrCreatingSession          = "error creando sesión: %w"
	ErrUpdatingSession          = "error actualizando sesión: %w"

	// Terms Session Logs
	LogSessionFoundReusing       = "Sesión existente encontrada, reutilizando token"
	LogSessionCreated            = "Sesión de términos creada exitosamente"
	LogSessionMarkedExpired      = "Error marcando sesión como expirada"
	LogSessionAlreadyAccepted    = "Sesión ya estaba aceptada, respuesta idempotente"
	LogTermsAccepted             = "Términos aceptados, iniciando notificación a Infobip"
	LogSessionAlreadyRejected    = "Sesión ya estaba rechazada, respuesta idempotente"
	LogTermsRejected             = "Términos rechazados, iniciando notificación a Infobip"
	LogRetryingInfobip           = "Reintentando notificación a Infobip"
	LogErrorUpdatingNotifyStatus = "Error actualizando estado de notificación exitosa"
	LogInfobipSuccess            = "Notificación a Infobip exitosa"
	LogInfobipFailed             = "Fallo en notificación a Infobip"
	LogInfobipFailedAll          = "Notificación a Infobip falló después de todos los reintentos"
	LogErrorUpdatingNotifyFailed = "Error actualizando estado de notificación fallida"

	// Terms Session Events
	EventTermsAccepted = "TERMS_ACCEPTED"
	EventTermsRejected = "TERMS_REJECTED"

	// Delivery with Terms Messages
	MsgInvalidData               = "Datos inv\u00e1lidos"
	MsgValidationFailed          = "Validaci\u00f3n fallida"
	MsgAtLeastOneDispenser       = "Debe incluir al menos un dispenser"
	MsgDispenserQuantityMismatch = "La cantidad de dispensers no coincide con el campo 'cantidad'"
	MsgServerError               = "Error del servidor"
	MsgCouldNotInitiateDelivery  = "No se pudo iniciar la entrega"
	MsgDeliveryInitiatedSuccess  = "Entrega iniciada exitosamente"
	MsgParameterMissing          = "Par\u00e1metro faltante"

	MsgCouldNotCompleteDelivery  = "No se pudo completar la entrega"
	MsgDeliveryCompletedSuccess  = "Entrega completada exitosamente"
	MsgDeliveryCreatedAfterTerms = "Entrega creada exitosamente despu\u00e9s de aceptar t\u00e9rminos"
	MsgUseTermsStatusEndpoint    = "Use el endpoint /api/v1/terms/status/:token para verificar el estado de los t\u00e9rminos"

	// Delivery with Terms Errors
	ErrTermsSessionNotFound     = "sesi\u00f3n de t\u00e9rminos no encontrada"
	ErrTermsNotAcceptedPending  = "los t\u00e9rminos no han sido aceptados (estado: PENDING)"
	ErrTermsNotAcceptedRejected = "los t\u00e9rminos no han sido aceptados (estado: REJECTED)"
	ErrTermsSessionExpired      = "la sesi\u00f3n de t\u00e9rminos ha expirado"

	// Delivery with Terms Logs
	LogValidationFailedInitiate = "Validaci\u00f3n fall\u00f3 en InitiateDelivery"
	LogErrorInitiatingDelivery  = "Error iniciando entrega"
	LogErrorCompletingDelivery  = "Error completando entrega"

	// Terms Session Handler Messages
	MsgInvalidRequest            = "Solicitud inv\u00e1lida"
	MsgErrorCreatingTermsSession = "Error creando sesi\u00f3n de t\u00e9rminos"
	MsgTokenRequired             = "Token requerido"
	MsgSessionIDRequired         = "SessionID requerido"
	MsgSessionNotFound           = "Sesi\u00f3n no encontrada"
	MsgAcceptingTerms            = "Aceptando t\u00e9rminos y condiciones"
	MsgRejectingTerms            = "Rechazando t\u00e9rminos y condiciones"

	// Terms Session Handler Logs
	LogErrorValidatingInfobip  = "Error validando request de Infobip"
	LogCreatingTermsSession    = "Creando sesi\u00f3n de t\u00e9rminos para Infobip"
	LogErrorCreatingSession    = "Error creando sesi\u00f3n de t\u00e9rminos"
	LogQueryingTermsStatus     = "Consultando estado de t\u00e9rminos"
	LogErrorGettingTermsStatus = "Error obteniendo estado de t\u00e9rminos"
	LogQueryingBySessionID     = "Consultando sesi\u00f3n por sessionID"
	LogErrorGettingBySessionID = "Error obteniendo sesi\u00f3n por sessionID"
	LogErrorAcceptingTerms     = "Error aceptando t\u00e9rminos"
	LogErrorRejectingTerms     = "Error rechazando t\u00e9rminos"
)
