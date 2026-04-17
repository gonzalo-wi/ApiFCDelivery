package service

import (
	"GoFrioCalor/internal/store"
	"context"
	"time"

	"github.com/rs/zerolog/log"
)

// Scheduler ejecuta tareas programadas periódicas
type Scheduler struct {
	deliveryStore store.DeliveryStore
	stopCh        chan struct{}
}

func NewScheduler(deliveryStore store.DeliveryStore) *Scheduler {
	return &Scheduler{
		deliveryStore: deliveryStore,
		stopCh:        make(chan struct{}),
	}
}

// Start inicia el scheduler. Calcula el tiempo hasta las 00:00 y luego repite cada 24hs.
func (s *Scheduler) Start() {
	go func() {
		// Calcular duración hasta la próxima medianoche
		now := time.Now()
		nextMidnight := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())
		untilMidnight := time.Until(nextMidnight)

		log.Info().
			Str("next_run", nextMidnight.Format("2006-01-02 15:04:05")).
			Str("wait", untilMidnight.String()).
			Msg("Scheduler started, waiting for first run at midnight")

		select {
		case <-time.After(untilMidnight):
			s.cancelExpiredDeliveries()
		case <-s.stopCh:
			log.Info().Msg("Scheduler stopped before first run")
			return
		}

		// Repetir cada 24 horas
		ticker := time.NewTicker(24 * time.Hour)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				s.cancelExpiredDeliveries()
			case <-s.stopCh:
				log.Info().Msg("Scheduler stopped")
				return
			}
		}
	}()
}

// Stop detiene el scheduler
func (s *Scheduler) Stop() {
	close(s.stopCh)
}

func (s *Scheduler) cancelExpiredDeliveries() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	count, err := s.deliveryStore.CancelExpiredPending(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Scheduler: error cancelling expired deliveries")
		return
	}

	log.Info().
		Int64("cancelled", count).
		Msg("Scheduler: expired pending deliveries cancelled")
}
