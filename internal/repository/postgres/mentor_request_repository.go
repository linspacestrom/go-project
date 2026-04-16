package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/linspacestrom/go-project/internal/domain"
)

func (r *Repository) CreateMentorRequest(ctx context.Context, req *domain.MentorRequest) (*domain.MentorRequest, error) {
	query, args, err := psql.Insert("mentor_requests").
		Columns("slot_id", "mentor_id", "mentee_id", "city_id", "skill_id", "request_type", "status", "comment").
		Values(req.SlotID, req.MentorID, req.MenteeID, req.CityID, req.SkillID, req.RequestType, req.Status, req.Comment).
		Suffix("RETURNING id, created_at, updated_at").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build create mentor request query: %w", err)
	}
	created := *req
	if err = r.GetConn(ctx).QueryRow(ctx, query, args...).Scan(&created.ID, &created.CreatedAt, &created.UpdatedAt); err != nil {
		return nil, fmt.Errorf("create mentor request: %w", err)
	}

	return &created, nil
}

func (r *Repository) ListMentorRequests(ctx context.Context, filter MentorRequestFilter) ([]domain.MentorRequest, uint64, error) {
	q := psql.Select("id", "slot_id", "mentor_id", "mentee_id", "city_id", "skill_id", "request_type", "status", "comment", "created_at", "updated_at").From("mentor_requests")
	countQ := psql.Select("COUNT(1)").From("mentor_requests")

	if filter.MenteeID != nil {
		q = q.Where(squirrel.Eq{"mentee_id": *filter.MenteeID})
		countQ = countQ.Where(squirrel.Eq{"mentee_id": *filter.MenteeID})
	}
	if filter.MentorID != nil {
		q = q.Where(squirrel.Eq{"mentor_id": *filter.MentorID})
		countQ = countQ.Where(squirrel.Eq{"mentor_id": *filter.MentorID})
	}
	if filter.SkillID != nil {
		q = q.Where(squirrel.Eq{"skill_id": *filter.SkillID})
		countQ = countQ.Where(squirrel.Eq{"skill_id": *filter.SkillID})
	}
	if filter.CityID != nil {
		q = q.Where(squirrel.Eq{"city_id": *filter.CityID})
		countQ = countQ.Where(squirrel.Eq{"city_id": *filter.CityID})
	}
	if filter.Status != nil {
		q = q.Where(squirrel.Eq{"status": *filter.Status})
		countQ = countQ.Where(squirrel.Eq{"status": *filter.Status})
	}

	query, args, err := q.OrderBy("created_at DESC").Limit(filter.Limit).Offset(filter.Offset).ToSql()
	if err != nil {
		return nil, 0, fmt.Errorf("build list mentor requests query: %w", err)
	}

	rows, err := r.GetConn(ctx).Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("list mentor requests: %w", err)
	}
	defer rows.Close()

	requests := make([]domain.MentorRequest, 0)
	for rows.Next() {
		var req domain.MentorRequest
		if scanErr := rows.Scan(&req.ID, &req.SlotID, &req.MentorID, &req.MenteeID, &req.CityID, &req.SkillID, &req.RequestType, &req.Status, &req.Comment, &req.CreatedAt, &req.UpdatedAt); scanErr != nil {
			return nil, 0, fmt.Errorf("scan mentor request: %w", scanErr)
		}
		requests = append(requests, req)
	}

	countQuery, countArgs, err := countQ.ToSql()
	if err != nil {
		return nil, 0, fmt.Errorf("build count mentor requests query: %w", err)
	}
	var total uint64
	if err = r.GetConn(ctx).QueryRow(ctx, countQuery, countArgs...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count mentor requests: %w", err)
	}

	return requests, total, nil
}

func (r *Repository) UpdateMentorRequestStatus(ctx context.Context, requestID uuid.UUID, status string, mentorID *uuid.UUID) (*domain.MentorRequest, error) {
	builder := psql.Update("mentor_requests").
		Set("status", status).
		Set("updated_at", squirrel.Expr("now()"))
	if mentorID != nil {
		builder = builder.Set("mentor_id", *mentorID)
	}

	query, args, err := builder.Where(squirrel.Eq{"id": requestID}).
		Suffix("RETURNING id, slot_id, mentor_id, mentee_id, city_id, skill_id, request_type, status, comment, created_at, updated_at").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build update mentor request query: %w", err)
	}

	var req domain.MentorRequest
	if err = r.GetConn(ctx).QueryRow(ctx, query, args...).Scan(
		&req.ID,
		&req.SlotID,
		&req.MentorID,
		&req.MenteeID,
		&req.CityID,
		&req.SkillID,
		&req.RequestType,
		&req.Status,
		&req.Comment,
		&req.CreatedAt,
		&req.UpdatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrMentorRequestNotFound
		}
		return nil, fmt.Errorf("update mentor request status: %w", err)
	}

	return &req, nil
}

func (r *Repository) UpsertMentorSkillSubscription(ctx context.Context, mentorID, skillID uuid.UUID) error {
	query, args, err := psql.Insert("mentor_skill_subscriptions").
		Columns("mentor_id", "skill_id", "is_active").
		Values(mentorID, skillID, true).
		Suffix("ON CONFLICT (mentor_id, skill_id) DO UPDATE SET is_active = TRUE, updated_at = now()").
		ToSql()
	if err != nil {
		return fmt.Errorf("build upsert mentor subscription query: %w", err)
	}
	if _, err = r.GetConn(ctx).Exec(ctx, query, args...); err != nil {
		return fmt.Errorf("upsert mentor subscription: %w", err)
	}

	return nil
}

func (r *Repository) DeleteMentorSkillSubscription(ctx context.Context, mentorID, skillID uuid.UUID) error {
	query, args, err := psql.Update("mentor_skill_subscriptions").
		Set("is_active", false).
		Set("updated_at", squirrel.Expr("now()")).
		Where(squirrel.Eq{"mentor_id": mentorID, "skill_id": skillID}).
		ToSql()
	if err != nil {
		return fmt.Errorf("build disable mentor subscription query: %w", err)
	}
	if _, err = r.GetConn(ctx).Exec(ctx, query, args...); err != nil {
		return fmt.Errorf("disable mentor subscription: %w", err)
	}

	return nil
}

func (r *Repository) ListMentorSubscribersBySkill(ctx context.Context, skillID uuid.UUID, cityID uuid.UUID) ([]uuid.UUID, error) {
	query, args, err := psql.Select("mss.mentor_id").
		From("mentor_skill_subscriptions mss").
		Join("users u ON u.id = mss.mentor_id").
		Where(squirrel.Eq{"mss.skill_id": skillID, "mss.is_active": true, "u.city_id": cityID}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build list mentor subscribers query: %w", err)
	}

	rows, err := r.GetConn(ctx).Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list mentor subscribers: %w", err)
	}
	defer rows.Close()

	ids := make([]uuid.UUID, 0)
	for rows.Next() {
		var id uuid.UUID
		if scanErr := rows.Scan(&id); scanErr != nil {
			return nil, fmt.Errorf("scan mentor subscriber id: %w", scanErr)
		}
		ids = append(ids, id)
	}

	return ids, nil
}
