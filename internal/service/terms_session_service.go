package service

import (
	"GoFrioCalor/internal/constants"
	"GoFrioCalor/internal/dto"
	"GoFrioCalor/internal/models"
	"GoFrioCalor/internal/store"
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
)

type TermsSessionService interface {
	CreateSession(ctx context.Context, sessionID, appBaseURL string, ttlHours int) (*dto.CreateTermsSessionResponse, error)
	GetSessionStatus(ctx context.Context, token string) (*dto.TermsSessionStatusResponse, error)
	GetSessionBySessionID(ctx context.Context, sessionID string) (*dto.TermsSessionStatusResponse, error)
	AcceptTerms(ctx context.Context, token, ip, userAgent string) (*dto.TermsActionResponse, error)
	RejectTerms(ctx context.Context, token, ip, userAgent string) (*dto.TermsActionResponse, error)
}

type termsSessionService struct {
	store         store.TermsSessionStore
	infobipClient InfobipClient
	maxRetries    int
	retryDelays   []time.Duration
}

func NewTermsSessionService(
	store store.TermsSessionStore,
	infobipClient InfobipClient,
) TermsSessionService {
	return &termsSessionService{
		store:         store,
		infobipClient: infobipClient,
		maxRetries:    3,
		retryDelays:   []time.Duration{1 * time.Second, 3 * time.Second, 7 * time.Second},
	}
}

// CreateSession crea una nueva sesión de términos y genera un token único
func (s *termsSessionService) CreateSession(ctx context.Context, sessionID, appBaseURL string, ttlHours int) (*dto.CreateTermsSessionResponse, error) {
	existing, err := s.store.FindBySessionID(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf(constants.ErrVerifyingExistingSession, err)
	}
	if existing != nil && existing.Status == models.StatusPending && time.Now().Before(existing.ExpiresAt) {
		log.Info().
			Str("session_id", sessionID).
			Str("token", existing.Token).
			Msg(constants.LogSessionFoundReusing)
		return &dto.CreateTermsSessionResponse{
			Token:     existing.Token,
			URL:       fmt.Sprintf("%s/terms/%s", appBaseURL, existing.Token),
			ExpiresAt: existing.ExpiresAt,
		}, nil
	}
	token, err := generateSecureToken()
	if err != nil {
		return nil, fmt.Errorf(constants.ErrGeneratingToken, err)
	}
	now := time.Now()
	expiresAt := now.Add(time.Duration(ttlHours) * time.Hour)
	session := &models.TermsSession{
		Token:        token,
		SessionID:    sessionID,
		Status:       models.StatusPending,
		CreatedAt:    now,
		ExpiresAt:    expiresAt,
		NotifyStatus: models.NotifyPending,
	}
	if err := s.store.Create(ctx, session); err != nil {
		return nil, fmt.Errorf(constants.ErrCreatingSession, err)
	}
	log.Info().
		Str("session_id", sessionID).
		Str("token", token).
		Time("expires_at", expiresAt).
		Msg(constants.LogSessionCreated)
	return &dto.CreateTermsSessionResponse{
		Token:     token,
		URL:       fmt.Sprintf("%s/terms/%s", appBaseURL, token),
		ExpiresAt: expiresAt,
	}, nil
}

// GetSessionStatus obtiene el estado actual de una sesión
func (s *termsSessionService) GetSessionStatus(ctx context.Context, token string) (*dto.TermsSessionStatusResponse, error) {
	session, err := s.store.FindByToken(ctx, token)
	if err != nil {
		return nil, err
	}
	// Verificar si está expirada
	if session.Status == models.StatusPending && time.Now().After(session.ExpiresAt) {
		if err := s.store.MarkExpired(ctx, token); err != nil {
			log.Error().Err(err).Str("token", token).Msg(constants.LogSessionMarkedExpired)
		}
		session.Status = models.StatusExpired
	}
	return &dto.TermsSessionStatusResponse{
		Status:     session.Status,
		ExpiresAt:  session.ExpiresAt,
		AcceptedAt: session.AcceptedAt,
		RejectedAt: session.RejectedAt,
	}, nil
}

// AcceptTerms marca los términos como aceptados y notifica a Infobip
func (s *termsSessionService) AcceptTerms(ctx context.Context, token, ip, userAgent string) (*dto.TermsActionResponse, error) {
	session, err := s.store.FindByToken(ctx, token)
	if err != nil {
		return nil, err
	}
	// Validar estado y expiración
	if err := s.validateSessionForAction(session); err != nil {
		return nil, err
	}
	// Si ya está aceptado (idempotencia)
	if session.Status == models.StatusAccepted {
		log.Info().Str("token", token).Msg(constants.LogSessionAlreadyAccepted)
		return &dto.TermsActionResponse{
			Status:     models.StatusAccepted,
			Message:    constants.MsgTermsAlreadyAccepted,
			AcceptedAt: session.AcceptedAt,
		}, nil
	}
	// Actualizar sesión con datos de aceptación
	now := time.Now()
	session.Status = models.StatusAccepted
	session.AcceptedAt = &now
	session.IP = ip
	session.UserAgent = userAgent
	if err := s.store.Update(ctx, session); err != nil {
		return nil, fmt.Errorf(constants.ErrUpdatingSession, err)
	}
	log.Info().
		Str("token", token).
		Str("session_id", session.SessionID).
		Str("ip", ip).
		Msg(constants.LogTermsAccepted)
	// Notificar a Infobip (con reintentos)
	go s.notifyInfobipWithRetries(context.Background(), session, constants.EventTermsAccepted)
	return &dto.TermsActionResponse{
		Status:     models.StatusAccepted,
		Message:    constants.MsgTermsAcceptedSuccess,
		AcceptedAt: session.AcceptedAt,
	}, nil
}

// RejectTerms marca los términos como rechazados y notifica a Infobip
func (s *termsSessionService) RejectTerms(ctx context.Context, token, ip, userAgent string) (*dto.TermsActionResponse, error) {
	session, err := s.store.FindByToken(ctx, token)
	if err != nil {
		return nil, err
	}
	// Validar estado y expiración
	if err := s.validateSessionForAction(session); err != nil {
		return nil, err
	}
	// Si ya está rechazado (idempotencia)
	if session.Status == models.StatusRejected {
		log.Info().Str("token", token).Msg(constants.LogSessionAlreadyRejected)
		return &dto.TermsActionResponse{
			Status:     models.StatusRejected,
			Message:    constants.MsgTermsAlreadyRejected,
			RejectedAt: session.RejectedAt,
		}, nil
	}
	// Actualizar sesión con datos de rechazo
	now := time.Now()
	session.Status = models.StatusRejected
	session.RejectedAt = &now
	session.IP = ip
	session.UserAgent = userAgent
	if err := s.store.Update(ctx, session); err != nil {
		return nil, fmt.Errorf(constants.ErrUpdatingSession, err)
	}
	log.Info().
		Str("token", token).
		Str("session_id", session.SessionID).
		Str("ip", ip).
		Msg(constants.LogTermsRejected)
	// Notificar a Infobip (con reintentos)
	go s.notifyInfobipWithRetries(context.Background(), session, constants.EventTermsRejected)
	return &dto.TermsActionResponse{
		Status:     models.StatusRejected,
		Message:    constants.MsgTermsRejected,
		RejectedAt: session.RejectedAt,
	}, nil
}

// GetSessionBySessionID obtiene el estado de una sesión usando el sessionID (para frontend/Infobip)
func (s *termsSessionService) GetSessionBySessionID(ctx context.Context, sessionID string) (*dto.TermsSessionStatusResponse, error) {
	session, err := s.store.FindBySessionID(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf(constants.MsgSessionNotFound, sessionID)
	}
	// Verificar si está expirada
	if session.Status == models.StatusPending && time.Now().After(session.ExpiresAt) {
		if err := s.store.MarkExpired(ctx, session.Token); err != nil {
			log.Error().Err(err).Str("token", session.Token).Msg(constants.LogSessionMarkedExpired)
		}
		session.Status = models.StatusExpired
	}
	return &dto.TermsSessionStatusResponse{
		Token:      session.Token,
		Status:     session.Status,
		ExpiresAt:  session.ExpiresAt,
		AcceptedAt: session.AcceptedAt,
		RejectedAt: session.RejectedAt,
	}, nil
}

// validateSessionForAction valida que una sesión pueda ser aceptada/rechazada
func (s *termsSessionService) validateSessionForAction(session *models.TermsSession) error {
	// Verificar expiración
	if time.Now().After(session.ExpiresAt) {
		return fmt.Errorf(constants.MsgSessionExpired)
	}
	// Validar estado (ya aceptado o rechazado se maneja con idempotencia en los métodos Accept/Reject)
	if session.Status != models.StatusPending && session.Status != models.StatusAccepted && session.Status != models.StatusRejected {
		return fmt.Errorf(constants.MsgSessionNotAvailable, session.Status)
	}
	return nil
}

// notifyInfobipWithRetries envía notificación a Infobip con reintentos
func (s *termsSessionService) notifyInfobipWithRetries(ctx context.Context, session *models.TermsSession, event string) {
	payload := dto.InfobipWebhookPayload{
		Event:      event,
		SessionID:  session.SessionID,
		Token:      session.Token,
		AcceptedAt: session.AcceptedAt,
		RejectedAt: session.RejectedAt,
	}
	var lastError error
	for attempt := 0; attempt < s.maxRetries; attempt++ {
		if attempt > 0 {
			delay := s.retryDelays[attempt-1]
			log.Info().
				Int("attempt", attempt+1).
				Dur("delay", delay).
				Str("session_id", session.SessionID).
				Msg(constants.LogRetryingInfobip)
			time.Sleep(delay)
		}
		err := s.infobipClient.SendWebhook(ctx, session.SessionID, payload)
		if err == nil {
			// Éxito
			if err := s.store.UpdateNotifyStatus(ctx, session.ID, models.NotifySent, attempt+1, ""); err != nil {
				log.Error().Err(err).Msg(constants.LogErrorUpdatingNotifyStatus)
			}
			log.Info().
				Str("session_id", session.SessionID).
				Int("attempts", attempt+1).
				Msg(constants.LogInfobipSuccess)
			return
		}
		lastError = err
		log.Warn().
			Err(err).
			Int("attempt", attempt+1).
			Int("max_retries", s.maxRetries).
			Str("session_id", session.SessionID).
			Msg(constants.LogInfobipFailed)
	}
	// Fallo después de todos los intentos
	errorMsg := ""
	if lastError != nil {
		errorMsg = lastError.Error()
	}
	if err := s.store.UpdateNotifyStatus(ctx, session.ID, models.NotifyFailed, s.maxRetries, errorMsg); err != nil {
		log.Error().Err(err).Msg(constants.LogErrorUpdatingNotifyFailed)
	}
	log.Error().
		Err(lastError).
		Str("session_id", session.SessionID).
		Int("attempts", s.maxRetries).
		Msg(constants.LogInfobipFailedAll)
}

// generateSecureToken genera un token seguro usando crypto/rand
func generateSecureToken() (string, error) {
	bytes := make([]byte, 32) // 32 bytes = 64 caracteres hex
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
