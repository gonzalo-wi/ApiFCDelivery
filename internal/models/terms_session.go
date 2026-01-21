package models

import "time"

type TermsSessionStatus string

const (
	StatusPending  TermsSessionStatus = "PENDING"
	StatusAccepted TermsSessionStatus = "ACCEPTED"
	StatusRejected TermsSessionStatus = "REJECTED"
	StatusExpired  TermsSessionStatus = "EXPIRED"
)

type NotifyStatus string

const (
	NotifyPending NotifyStatus = "PENDING"
	NotifySent    NotifyStatus = "SENT"
	NotifyFailed  NotifyStatus = "FAILED"
)

type TermsSession struct {
	ID             int64              `gorm:"primaryKey" json:"id"`
	Token          string             `gorm:"uniqueIndex;not null" json:"token"`
	SessionID      string             `gorm:"not null;index" json:"session_id"`
	Status         TermsSessionStatus `gorm:"not null;index" json:"status"`
	DeliveryData   string             `gorm:"type:text" json:"delivery_data,omitempty"`
	CreatedAt      time.Time          `gorm:"autoCreateTime" json:"created_at"`
	ExpiresAt      time.Time          `gorm:"not null;index" json:"expires_at"`
	AcceptedAt     *time.Time         `json:"accepted_at,omitempty"`
	RejectedAt     *time.Time         `json:"rejected_at,omitempty"`
	IP             string             `json:"ip,omitempty"`
	UserAgent      string             `json:"user_agent,omitempty"`
	NotifyStatus   NotifyStatus       `gorm:"not null;default:'PENDING'" json:"notify_status"`
	NotifyAttempts int                `gorm:"default:0" json:"notify_attempts"`
	LastError      string             `json:"last_error,omitempty"`
}
