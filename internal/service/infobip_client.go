package service

import (
	"GoFrioCalor/internal/dto"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

type InfobipClient interface {
	SendWebhook(ctx context.Context, sessionID string, payload dto.InfobipWebhookPayload) error
}

type infobipClient struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

func NewInfobipClient(baseURL, apiKey string) InfobipClient {
	return &infobipClient{
		baseURL: baseURL,
		apiKey:  apiKey,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *infobipClient) SendWebhook(ctx context.Context, sessionID string, payload dto.InfobipWebhookPayload) error {
	url := fmt.Sprintf("%s/bots/webhook/%s", c.baseURL, sessionID)

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error serializando payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error creando request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("App %s", c.apiKey))

	log.Debug().
		Str("url", url).
		Str("session_id", sessionID).
		Str("payload", string(jsonData)).
		Msg("Enviando webhook a Infobip")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("error enviando webhook: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		log.Error().
			Int("status_code", resp.StatusCode).
			Str("response_body", string(body)).
			Msg("Respuesta error de Infobip")
		return fmt.Errorf("infobip respondió con status %d: %s", resp.StatusCode, string(body))
	}

	log.Info().
		Str("session_id", sessionID).
		Int("status_code", resp.StatusCode).
		Msg("Webhook enviado exitosamente a Infobip")

	return nil
}
