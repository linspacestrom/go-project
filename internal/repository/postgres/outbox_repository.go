package postgres

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/linspacestrom/go-project/internal/domain"
)

func (r *Repository) CreateOutboxEvent(ctx context.Context, aggregateType string, aggregateID *uuid.UUID, eventType string, payload []byte) error {
	query, args, err := psql.Insert("outbox_events").
		Columns("aggregate_type", "aggregate_id", "event_type", "payload", "status").
		Values(aggregateType, aggregateID, eventType, payload, domain.OutboxStatusNew).
		ToSql()
	if err != nil {
		return fmt.Errorf("build create outbox event query: %w", err)
	}
	if _, err = r.GetConn(ctx).Exec(ctx, query, args...); err != nil {
		return fmt.Errorf("create outbox event: %w", err)
	}

	return nil
}

func (r *Repository) ListOutboxForDispatch(ctx context.Context, limit uint64) ([]domain.OutboxEvent, error) {
	query, args, err := psql.Select("id", "aggregate_type", "aggregate_id", "event_type", "payload", "status", "error_message", "created_at", "updated_at").
		From("outbox_events").
		Where(squirrel.Eq{"status": domain.OutboxStatusNew}).
		OrderBy("created_at ASC").
		Limit(limit).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build list outbox query: %w", err)
	}

	rows, err := r.GetConn(ctx).Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list outbox events: %w", err)
	}
	defer rows.Close()

	items := make([]domain.OutboxEvent, 0)
	for rows.Next() {
		var item domain.OutboxEvent
		if scanErr := rows.Scan(&item.ID, &item.AggregateType, &item.AggregateID, &item.EventType, &item.Payload, &item.Status, &item.ErrorMessage, &item.CreatedAt, &item.UpdatedAt); scanErr != nil {
			return nil, fmt.Errorf("scan outbox event: %w", scanErr)
		}
		items = append(items, item)
	}

	return items, nil
}

func (r *Repository) MarkOutboxSent(ctx context.Context, id uuid.UUID) error {
	query, args, err := psql.Update("outbox_events").
		Set("status", domain.OutboxStatusSent).
		Set("updated_at", squirrel.Expr("now()")).
		Set("error_message", nil).
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return fmt.Errorf("build mark outbox sent query: %w", err)
	}
	if _, err = r.GetConn(ctx).Exec(ctx, query, args...); err != nil {
		return fmt.Errorf("mark outbox sent: %w", err)
	}

	return nil
}

func (r *Repository) MarkOutboxFailed(ctx context.Context, id uuid.UUID, errorMessage string) error {
	query, args, err := psql.Update("outbox_events").
		Set("status", domain.OutboxStatusFailed).
		Set("updated_at", squirrel.Expr("now()")).
		Set("error_message", errorMessage).
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return fmt.Errorf("build mark outbox failed query: %w", err)
	}
	if _, err = r.GetConn(ctx).Exec(ctx, query, args...); err != nil {
		return fmt.Errorf("mark outbox failed: %w", err)
	}

	return nil
}
