package outbox

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/linspacestrom/go-project/internal/domain"
	"go.uber.org/zap"
)

type Repository interface {
	ListOutboxForDispatch(ctx context.Context, limit uint64) ([]domain.OutboxEvent, error)
	MarkOutboxSent(ctx context.Context, id uuid.UUID) error
	MarkOutboxFailed(ctx context.Context, id uuid.UUID, errorMessage string) error
}

type Producer interface {
	Publish(ctx context.Context, key uuid.UUID, eventType string, payload []byte) error
}

type Service struct {
	log      *zap.Logger
	repo     Repository
	producer Producer
	interval time.Duration
	limit    uint64
}

func NewService(log *zap.Logger, repo Repository, producer Producer, interval time.Duration, limit uint64) *Service {
	return &Service{log: log, repo: repo, producer: producer, interval: interval, limit: limit}
}

func (s *Service) Run(ctx context.Context) {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s.dispatchOnce(ctx)
		}
	}
}

func (s *Service) dispatchOnce(ctx context.Context) {
	items, err := s.repo.ListOutboxForDispatch(ctx, s.limit)
	if err != nil {
		s.log.Error("failed to fetch outbox events", zap.Error(err))
		return
	}

	for _, item := range items {
		aggregateID := uuid.Nil
		if item.AggregateID != nil {
			aggregateID = *item.AggregateID
		}
		if pubErr := s.producer.Publish(ctx, aggregateID, item.EventType, item.Payload); pubErr != nil {
			s.log.Error("failed to publish outbox event", zap.Error(pubErr), zap.String("event_id", item.ID.String()))
			if markErr := s.repo.MarkOutboxFailed(ctx, item.ID, pubErr.Error()); markErr != nil {
				s.log.Error("failed to mark outbox failed", zap.Error(markErr), zap.String("event_id", item.ID.String()))
			}
			continue
		}

		if markErr := s.repo.MarkOutboxSent(ctx, item.ID); markErr != nil {
			s.log.Error("failed to mark outbox sent", zap.Error(markErr), zap.String("event_id", item.ID.String()))
		}
	}
}

func (s *Service) DispatchNow(ctx context.Context) error {
	items, err := s.repo.ListOutboxForDispatch(ctx, s.limit)
	if err != nil {
		return err
	}

	for _, item := range items {
		aggregateID := uuid.Nil
		if item.AggregateID != nil {
			aggregateID = *item.AggregateID
		}
		if pubErr := s.producer.Publish(ctx, aggregateID, item.EventType, item.Payload); pubErr != nil {
			if markErr := s.repo.MarkOutboxFailed(ctx, item.ID, pubErr.Error()); markErr != nil {
				return fmt.Errorf("publish failed: %w, mark failed: %w", pubErr, markErr)
			}
			continue
		}
		if markErr := s.repo.MarkOutboxSent(ctx, item.ID); markErr != nil {
			return markErr
		}
	}

	return nil
}
