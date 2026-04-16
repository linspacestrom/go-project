package postgres

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

func (r *Repository) InsertAnalyticsEvent(ctx context.Context, eventType, entityType string, entityID *uuid.UUID, payload []byte) error {
	query, args, err := psql.Insert("analytics_events").
		Columns("event_type", "entity_type", "entity_id", "payload").
		Values(eventType, entityType, entityID, payload).
		ToSql()
	if err != nil {
		return fmt.Errorf("build insert analytics event query: %w", err)
	}
	if _, err = r.GetConn(ctx).Exec(ctx, query, args...); err != nil {
		return fmt.Errorf("insert analytics event: %w", err)
	}

	return nil
}

func (r *Repository) GetBusinessAnalytics(ctx context.Context) (*struct {
	MentorApprovals                int64
	MentorRejections               int64
	ActiveBookings                 int64
	MentorRequestsBySkill          int64
	MentorRequestsByCategory       int64
	UserCancelsAfterMentorApproval int64
}, error) {
	result := &struct {
		MentorApprovals                int64
		MentorRejections               int64
		ActiveBookings                 int64
		MentorRequestsBySkill          int64
		MentorRequestsByCategory       int64
		UserCancelsAfterMentorApproval int64
	}{}

	queries := []struct {
		builder squirrel.SelectBuilder
		dest    *int64
	}{
		{psql.Select("COUNT(1)").From("mentor_requests").Where(squirrel.Eq{"status": "approved"}), &result.MentorApprovals},
		{psql.Select("COUNT(1)").From("mentor_requests").Where(squirrel.Eq{"status": "rejected"}), &result.MentorRejections},
		{psql.Select("COUNT(1)").From("bookings").Where(squirrel.Eq{"status": "active"}), &result.ActiveBookings},
		{psql.Select("COUNT(1)").From("mentor_requests").Where(squirrel.Eq{"request_type": "skill"}), &result.MentorRequestsBySkill},
		{psql.Select("COUNT(1)").From("mentor_requests").Where(squirrel.Eq{"request_type": "category"}), &result.MentorRequestsByCategory},
		{psql.Select("COUNT(1)").From("analytics_events").Where(squirrel.Eq{"event_type": "booking_cancel_after_mentor_approval"}), &result.UserCancelsAfterMentorApproval},
	}

	for _, q := range queries {
		query, args, err := q.builder.ToSql()
		if err != nil {
			return nil, fmt.Errorf("build business analytics query: %w", err)
		}
		if err = r.GetConn(ctx).QueryRow(ctx, query, args...).Scan(q.dest); err != nil {
			return nil, fmt.Errorf("business analytics query failed: %w", err)
		}
	}

	return result, nil
}

func (r *Repository) GetTechnicalAnalytics(ctx context.Context) (*struct {
	BookingConflictCount int64
	FailedOutboxEvents   int64
	OutboxLagCount       int64
}, error) {
	result := &struct {
		BookingConflictCount int64
		FailedOutboxEvents   int64
		OutboxLagCount       int64
	}{}

	queries := []struct {
		builder squirrel.SelectBuilder
		dest    *int64
	}{
		{psql.Select("COUNT(1)").From("analytics_events").Where(squirrel.Eq{"event_type": "booking_conflict"}), &result.BookingConflictCount},
		{psql.Select("COUNT(1)").From("outbox_events").Where(squirrel.Eq{"status": "failed"}), &result.FailedOutboxEvents},
		{psql.Select("COUNT(1)").From("outbox_events").Where(squirrel.Eq{"status": "new"}), &result.OutboxLagCount},
	}

	for _, q := range queries {
		query, args, err := q.builder.ToSql()
		if err != nil {
			return nil, fmt.Errorf("build technical analytics query: %w", err)
		}
		if err = r.GetConn(ctx).QueryRow(ctx, query, args...).Scan(q.dest); err != nil {
			return nil, fmt.Errorf("technical analytics query failed: %w", err)
		}
	}

	return result, nil
}
