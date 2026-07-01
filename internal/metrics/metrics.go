// Package metrics define las métricas de negocio del servicio, expuestas en /metrics
// junto con las métricas HTTP (middleware) y las del runtime de Go.
//
// Todas se registran en el registry por defecto de Prometheus, el mismo que sirve
// promhttp.Handler() en /metrics.
package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// DeliveriesCreatedTotal cuenta entregas creadas, segmentadas por tipo de entrega.
	DeliveriesCreatedTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "deliveries_created_total",
			Help: "Total de entregas creadas, por tipo de entrega.",
		},
		[]string{"tipo_entrega"},
	)

	// TermsActionsTotal cuenta aceptaciones/rechazos de términos, por acción y empresa.
	TermsActionsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "terms_actions_total",
			Help: "Total de acciones sobre términos y condiciones (accepted/rejected), por empresa.",
		},
		[]string{"action", "company"},
	)

	// EmailsSentTotal cuenta intentos de envío de email, por tipo y resultado (sent/error).
	EmailsSentTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "emails_sent_total",
			Help: "Total de emails procesados, por tipo y resultado (sent/error).",
		},
		[]string{"type", "result"},
	)
)

// DeliveryCreated registra la creación de una entrega.
func DeliveryCreated(tipoEntrega string) {
	if tipoEntrega == "" {
		tipoEntrega = "unknown"
	}
	DeliveriesCreatedTotal.WithLabelValues(tipoEntrega).Inc()
}

// TermsAction registra una acción de términos (accepted/rejected).
func TermsAction(action, company string) {
	if company == "" {
		company = "unknown"
	}
	TermsActionsTotal.WithLabelValues(action, company).Inc()
}

// EmailSent registra el resultado de un envío de email.
func EmailSent(emailType string, ok bool) {
	result := "sent"
	if !ok {
		result = "error"
	}
	EmailsSentTotal.WithLabelValues(emailType, result).Inc()
}
