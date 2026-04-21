package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

const fallbackClientEmail = "gwinazki@el-jumillano.com.ar"

// ClientLookupService resuelve el email de un cliente a partir de su nro_cta
// consultando la API externa de El Jumillano.
type ClientLookupService interface {
	GetClientEmail(ctx context.Context, nroCta string) string
}

type clientLookupService struct {
	baseURL      string
	apiKey       string // token estático para el endpoint /jmap2token/token
	defaultEmail string
	httpClient   *http.Client
}

// NewClientLookupService crea un ClientLookupService.
// baseURL: e.g. "https://servicios.el-jumillano.com.ar:8443"
// apiKey:  Bearer token para el endpoint de autenticación (valor de Authorization header)
// defaultEmail: email de respaldo cuando el cliente no tiene email registrado
func NewClientLookupService(baseURL, apiKey, defaultEmail string) ClientLookupService {
	if defaultEmail == "" {
		defaultEmail = fallbackClientEmail
	}
	return &clientLookupService{
		baseURL:      strings.TrimRight(baseURL, "/"),
		apiKey:       apiKey,
		defaultEmail: defaultEmail,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// GetClientEmail obtiene el email del cliente dado su nro_cta.
// Normaliza nro_cta a minúsculas y sin espacios antes de consultar la API.
// Si el cliente no tiene email o falla la consulta, retorna el email por defecto.
func (s *clientLookupService) GetClientEmail(ctx context.Context, nroCta string) string {
	// Normalizar: minúsculas, sin espacios
	nroCta = strings.ToLower(strings.ReplaceAll(nroCta, " ", ""))

	token, err := s.getToken(ctx)
	if err != nil {
		log.Warn().Err(err).Str("nro_cta", nroCta).Msg("Client lookup: failed to get token, using default email")
		return s.defaultEmail
	}

	email, err := s.lookupClientEmail(ctx, token, nroCta)
	if err != nil {
		log.Warn().Err(err).Str("nro_cta", nroCta).Msg("Client lookup: failed to look up client, using default email")
		return s.defaultEmail
	}

	if email == "" {
		log.Info().Str("nro_cta", nroCta).Str("default", s.defaultEmail).Msg("Client lookup: no email registered, using default")
		return s.defaultEmail
	}

	log.Info().Str("nro_cta", nroCta).Str("email", email).Msg("Client lookup: email resolved")
	return email
}

// tokenResponse representa la respuesta del endpoint /jmap2token/token
type tokenResponse struct {
	Token      string `json:"token"`
	Current    string `json:"current"`
	Expiration string `json:"expiration"`
}

func (s *clientLookupService) getToken(ctx context.Context) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, s.baseURL+"/jmap2token/token", nil)
	if err != nil {
		return "", fmt.Errorf("building token request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+s.apiKey)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("calling token endpoint: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("token endpoint returned HTTP %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("reading token response body: %w", err)
	}

	var tr tokenResponse
	if err = json.Unmarshal(body, &tr); err != nil {
		return "", fmt.Errorf("parsing token response: %w", err)
	}

	if tr.Token == "" {
		return "", fmt.Errorf("empty token in response")
	}

	return tr.Token, nil
}

// clientLookupResponse representa la respuesta del endpoint /jmap2/client
type clientLookupResponse struct {
	Success bool                 `json:"success"`
	Data    []clientLookupRecord `json:"data"`
	Message *string              `json:"message"`
}

type clientLookupRecord struct {
	CodCliente int    `json:"codCliente"`
	Nombre     string `json:"nombre"`
	Emails     string `json:"emails"`
}

func (s *clientLookupService) lookupClientEmail(ctx context.Context, token, nroCta string) (string, error) {
	url := fmt.Sprintf("%s/jmap2/client?nrocta=%s", s.baseURL, nroCta)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", fmt.Errorf("building client lookup request: %w", err)
	}
	req.Header.Set("Authorization", token)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("calling client lookup endpoint: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("client lookup returned HTTP %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("reading client lookup response: %w", err)
	}

	var cr clientLookupResponse
	if err = json.Unmarshal(body, &cr); err != nil {
		return "", fmt.Errorf("parsing client lookup response: %w", err)
	}

	if !cr.Success || len(cr.Data) == 0 {
		return "", nil
	}

	return strings.TrimSpace(cr.Data[0].Emails), nil
}
