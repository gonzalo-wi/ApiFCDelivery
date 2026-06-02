package dto

import (
	"GoFrioCalor/internal/models"
	"time"
)

type InfobipSessionRequest struct {
	SessionID      string `json:"sessionId" binding:"required"`
	ConversationID string `json:"conversationId" binding:"required"`
	Delivery       int    `json:"delivery"`
}

type CreateTermsSessionResponse struct {
	Token     string    `json:"token"`
	URL       string    `json:"url"`
	ExpiresAt time.Time `json:"expiresAt"`
	Company   string    `json:"company,omitempty"`
}

type TermsSessionStatusResponse struct {
	Token      string                    `json:"token,omitempty"`
	Status     models.TermsSessionStatus `json:"status"`
	ExpiresAt  time.Time                 `json:"expiresAt"`
	AcceptedAt *time.Time                `json:"acceptedAt,omitempty"`
	RejectedAt *time.Time                `json:"rejectedAt,omitempty"`
	Company    string                    `json:"company,omitempty"`
}

type TermsActionRequest struct {
	IP        string `json:"ip,omitempty"`
	UserAgent string `json:"userAgent,omitempty"`
}

type TermsActionResponse struct {
	Status     models.TermsSessionStatus `json:"status"`
	Message    string                    `json:"message"`
	AcceptedAt *time.Time                `json:"acceptedAt,omitempty"`
	RejectedAt *time.Time                `json:"rejectedAt,omitempty"`
}

type InfobipWebhookPayload struct {
	Acepta bool `json:"acepta"`
}
