package store

import (
	"GoFrioCalor/internal/models"
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type TermsSessionStore interface {
	Create(ctx context.Context, session *models.TermsSession) error
	FindByToken(ctx context.Context, token string) (*models.TermsSession, error)
	FindBySessionID(ctx context.Context, sessionID string) (*models.TermsSession, error)
	Update(ctx context.Context, session *models.TermsSession) error
	UpdateStatus(ctx context.Context, token string, status models.TermsSessionStatus) error
	UpdateNotifyStatus(ctx context.Context, id int64, notifyStatus models.NotifyStatus, attempts int, lastError string) error
	MarkExpired(ctx context.Context, token string) error
}

type termsSessionStore struct {
	db *gorm.DB
}

func NewTermsSessionStore(db *gorm.DB) TermsSessionStore {
	return &termsSessionStore{db: db}
}

func (s *termsSessionStore) Create(ctx context.Context, session *models.TermsSession) error {
	if err := s.db.WithContext(ctx).Create(session).Error; err != nil {
		return fmt.Errorf("error creando sesión de términos: %w", err)
	}
	return nil
}

func (s *termsSessionStore) FindByToken(ctx context.Context, token string) (*models.TermsSession, error) {
	var session models.TermsSession
	if err := s.db.WithContext(ctx).Where("token = ?", token).First(&session).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("sesión de términos no encontrada")
		}
		return nil, fmt.Errorf("error buscando sesión por token: %w", err)
	}
	return &session, nil
}

func (s *termsSessionStore) FindBySessionID(ctx context.Context, sessionID string) (*models.TermsSession, error) {
	var session models.TermsSession
	if err := s.db.WithContext(ctx).Where("session_id = ?", sessionID).First(&session).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // No error si no existe
		}
		return nil, fmt.Errorf("error buscando sesión por sessionID: %w", err)
	}
	return &session, nil
}

func (s *termsSessionStore) Update(ctx context.Context, session *models.TermsSession) error {
	if err := s.db.WithContext(ctx).Save(session).Error; err != nil {
		return fmt.Errorf("error actualizando sesión de términos: %w", err)
	}
	return nil
}

func (s *termsSessionStore) UpdateStatus(ctx context.Context, token string, status models.TermsSessionStatus) error {
	now := time.Now()
	updates := map[string]interface{}{
		"status": status,
	}

	if status == models.StatusAccepted {
		updates["accepted_at"] = now
	} else if status == models.StatusRejected {
		updates["rejected_at"] = now
	}

	if err := s.db.WithContext(ctx).Model(&models.TermsSession{}).
		Where("token = ?", token).
		Updates(updates).Error; err != nil {
		return fmt.Errorf("error actualizando estado de sesión: %w", err)
	}
	return nil
}

func (s *termsSessionStore) UpdateNotifyStatus(ctx context.Context, id int64, notifyStatus models.NotifyStatus, attempts int, lastError string) error {
	updates := map[string]interface{}{
		"notify_status":   notifyStatus,
		"notify_attempts": attempts,
		"last_error":      lastError,
	}

	if err := s.db.WithContext(ctx).Model(&models.TermsSession{}).
		Where("id = ?", id).
		Updates(updates).Error; err != nil {
		return fmt.Errorf("error actualizando estado de notificación: %w", err)
	}
	return nil
}

func (s *termsSessionStore) MarkExpired(ctx context.Context, token string) error {
	if err := s.db.WithContext(ctx).Model(&models.TermsSession{}).
		Where("token = ? AND status = ?", token, models.StatusPending).
		Update("status", models.StatusExpired).Error; err != nil {
		return fmt.Errorf("error marcando sesión como expirada: %w", err)
	}
	return nil
}
