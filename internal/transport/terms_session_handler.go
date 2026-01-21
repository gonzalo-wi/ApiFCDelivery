package transport

import (
	"GoFrioCalor/internal/constants"
	"GoFrioCalor/internal/dto"
	"GoFrioCalor/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type TermsSessionHandler struct {
	service    service.TermsSessionService
	appBaseURL string
	termsTTL   int
}

func NewTermsSessionHandler(service service.TermsSessionService, appBaseURL string, termsTTL int) *TermsSessionHandler {
	return &TermsSessionHandler{
		service:    service,
		appBaseURL: appBaseURL,
		termsTTL:   termsTTL,
	}
}

// CreateInfobipSession maneja la creación de una sesión desde Infobip
// POST /api/infobip/session
func (h *TermsSessionHandler) CreateInfobipSession(c *gin.Context) {
	ctx := c.Request.Context()

	var req dto.InfobipSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn().Err(err).Msg(constants.LogErrorValidatingInfobip)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   constants.MsgInvalidRequest,
			"details": err.Error(),
		})
		return
	}

	log.Info().
		Str("session_id", req.SessionID).
		Msg(constants.LogCreatingTermsSession)

	response, err := h.service.CreateSession(ctx, req.SessionID, h.appBaseURL, h.termsTTL)
	if err != nil {
		log.Error().Err(err).Str("session_id", req.SessionID).Msg(constants.LogErrorCreatingSession)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": constants.MsgErrorCreatingTermsSession,
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetTermsStatus obtiene el estado de una sesión de términos
// GET /api/terms/:token
func (h *TermsSessionHandler) GetTermsStatus(c *gin.Context) {
	ctx := c.Request.Context()
	token := c.Param("token")

	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": constants.MsgTokenRequired,
		})
		return
	}

	log.Debug().Str("token", token).Msg(constants.LogQueryingTermsStatus)

	response, err := h.service.GetSessionStatus(ctx, token)
	if err != nil {
		log.Error().Err(err).Str("token", token).Msg(constants.LogErrorGettingTermsStatus)
		c.JSON(http.StatusNotFound, gin.H{
			"error": constants.MsgSessionNotFound,
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// AcceptTerms acepta los términos y condiciones
// POST /api/terms/:token/accept
func (h *TermsSessionHandler) AcceptTerms(c *gin.Context) {
	ctx := c.Request.Context()
	token := c.Param("token")

	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": constants.MsgTokenRequired,
		})
		return
	}

	// Obtener IP y User-Agent del cliente
	ip := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")

	log.Info().
		Str("token", token).
		Str("ip", ip).
		Msg(constants.MsgAcceptingTerms)

	response, err := h.service.AcceptTerms(ctx, token, ip, userAgent)
	if err != nil {
		log.Error().Err(err).Str("token", token).Msg(constants.LogErrorAcceptingTerms)

		// Usar helper para determinar código de estado
		statusCode := GetHTTPStatusFromError(err)

		c.JSON(statusCode, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetSessionBySessionID obtiene el estado de una sesión usando el sessionID (para frontend)
// GET /api/v1/terms/by-session/:sessionId
func (h *TermsSessionHandler) GetSessionBySessionID(c *gin.Context) {
	ctx := c.Request.Context()
	sessionID := c.Param("sessionId")

	if sessionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": constants.MsgSessionIDRequired,
		})
		return
	}

	log.Debug().Str("session_id", sessionID).Msg(constants.LogQueryingBySessionID)

	response, err := h.service.GetSessionBySessionID(ctx, sessionID)
	if err != nil {
		log.Error().Err(err).Str("session_id", sessionID).Msg(constants.LogErrorGettingBySessionID)
		c.JSON(http.StatusNotFound, gin.H{
			"error": constants.MsgSessionNotFound,
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// RejectTerms rechaza los términos y condiciones
// POST /api/terms/:token/reject
func (h *TermsSessionHandler) RejectTerms(c *gin.Context) {
	ctx := c.Request.Context()
	token := c.Param("token")

	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": constants.MsgTokenRequired,
		})
		return
	}

	// Obtener IP y User-Agent del cliente
	ip := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")

	log.Info().
		Str("token", token).
		Str("ip", ip).
		Msg(constants.MsgRejectingTerms)

	response, err := h.service.RejectTerms(ctx, token, ip, userAgent)
	if err != nil {
		log.Error().Err(err).Str("token", token).Msg(constants.LogErrorRejectingTerms)

		// Usar helper para determinar código de estado
		statusCode := GetHTTPStatusFromError(err)

		c.JSON(statusCode, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}
