package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/linspacestrom/go-project/internal/domain"
)

func (r *Repository) CreateSlot(ctx context.Context, slot *domain.Slot) (*domain.Slot, error) {
	query, args, err := psql.Insert("slots").
		Columns("mentor_id", "room_id", "city_id", "type", "start_at", "end_at", "meeting_url", "status", "capacity", "created_by", "approved_by").
		Values(slot.MentorID, slot.RoomID, slot.CityID, slot.Type, slot.StartAt, slot.EndAt, slot.MeetingURL, slot.Status, slot.Capacity, slot.CreatedBy, slot.ApprovedBy).
		Suffix("RETURNING id, created_at, updated_at").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build create slot query: %w", err)
	}

	created := *slot
	if err = r.GetConn(ctx).QueryRow(ctx, query, args...).Scan(&created.ID, &created.CreatedAt, &created.UpdatedAt); err != nil {
		return nil, fmt.Errorf("create slot: %w", err)
	}

	return &created, nil
}

func (r *Repository) UpdateSlot(ctx context.Context, slotID uuid.UUID, patch map[string]any) (*domain.Slot, error) {
	builder := psql.Update("slots").Set("updated_at", squirrel.Expr("now()"))
	for col, val := range patch {
		builder = builder.Set(col, val)
	}

	query, args, err := builder.Where(squirrel.Eq{"id": slotID}).Suffix("RETURNING id, mentor_id, room_id, city_id, type, start_at, end_at, meeting_url, status, capacity, created_by, approved_by, created_at, updated_at").ToSql()
	if err != nil {
		return nil, fmt.Errorf("build update slot query: %w", err)
	}

	var slot domain.Slot
	if err = r.GetConn(ctx).QueryRow(ctx, query, args...).Scan(
		&slot.ID,
		&slot.MentorID,
		&slot.RoomID,
		&slot.CityID,
		&slot.Type,
		&slot.StartAt,
		&slot.EndAt,
		&slot.MeetingURL,
		&slot.Status,
		&slot.Capacity,
		&slot.CreatedBy,
		&slot.ApprovedBy,
		&slot.CreatedAt,
		&slot.UpdatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrSlotNotFound
		}

		return nil, fmt.Errorf("update slot: %w", err)
	}

	return &slot, nil
}

func (r *Repository) SoftDeleteSlot(ctx context.Context, slotID uuid.UUID) error {
	query, args, err := psql.Update("slots").
		Set("status", domain.SlotStatusDeleted).
		Set("updated_at", squirrel.Expr("now()")).
		Where(squirrel.Eq{"id": slotID}).
		ToSql()
	if err != nil {
		return fmt.Errorf("build delete slot query: %w", err)
	}
	tag, err := r.GetConn(ctx).Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("delete slot: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrSlotNotFound
	}

	return nil
}

func (r *Repository) LockSlotByID(ctx context.Context, slotID uuid.UUID) (*domain.Slot, error) {
	query, args, err := psql.Select("id", "mentor_id", "room_id", "city_id", "type", "start_at", "end_at", "meeting_url", "status", "capacity", "created_by", "approved_by", "created_at", "updated_at").
		From("slots").
		Where(squirrel.Eq{"id": slotID}).
		Suffix("FOR UPDATE").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build lock slot query: %w", err)
	}

	var slot domain.Slot
	if err = r.GetConn(ctx).QueryRow(ctx, query, args...).Scan(
		&slot.ID,
		&slot.MentorID,
		&slot.RoomID,
		&slot.CityID,
		&slot.Type,
		&slot.StartAt,
		&slot.EndAt,
		&slot.MeetingURL,
		&slot.Status,
		&slot.Capacity,
		&slot.CreatedBy,
		&slot.ApprovedBy,
		&slot.CreatedAt,
		&slot.UpdatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrSlotNotFound
		}

		return nil, fmt.Errorf("lock slot: %w", err)
	}

	return &slot, nil
}

func (r *Repository) GetSlotByID(ctx context.Context, slotID uuid.UUID) (*domain.Slot, error) {
	query, args, err := psql.Select("id", "mentor_id", "room_id", "city_id", "type", "start_at", "end_at", "meeting_url", "status", "capacity", "created_by", "approved_by", "created_at", "updated_at").
		From("slots").
		Where(squirrel.Eq{"id": slotID}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build get slot query: %w", err)
	}

	var slot domain.Slot
	if err = r.GetConn(ctx).QueryRow(ctx, query, args...).Scan(
		&slot.ID,
		&slot.MentorID,
		&slot.RoomID,
		&slot.CityID,
		&slot.Type,
		&slot.StartAt,
		&slot.EndAt,
		&slot.MeetingURL,
		&slot.Status,
		&slot.Capacity,
		&slot.CreatedBy,
		&slot.ApprovedBy,
		&slot.CreatedAt,
		&slot.UpdatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrSlotNotFound
		}
		return nil, fmt.Errorf("get slot: %w", err)
	}

	return &slot, nil
}

func (r *Repository) ListSlots(ctx context.Context, filter SlotFilter) ([]domain.Slot, uint64, error) {
	q := psql.Select("s.id", "s.mentor_id", "s.room_id", "s.city_id", "s.type", "s.start_at", "s.end_at", "s.meeting_url", "s.status", "s.capacity", "s.created_by", "s.approved_by", "s.created_at", "s.updated_at").
		From("slots s")
	countQ := psql.Select("COUNT(1)").From("slots s")

	if filter.HubID != nil {
		q = q.Join("rooms r ON r.id = s.room_id").Where(squirrel.Eq{"r.hub_id": *filter.HubID})
		countQ = countQ.Join("rooms r ON r.id = s.room_id").Where(squirrel.Eq{"r.hub_id": *filter.HubID})
	}
	if filter.CityID != nil {
		q = q.Where(squirrel.Eq{"s.city_id": *filter.CityID})
		countQ = countQ.Where(squirrel.Eq{"s.city_id": *filter.CityID})
	}
	if filter.RoomID != nil {
		q = q.Where(squirrel.Eq{"s.room_id": *filter.RoomID})
		countQ = countQ.Where(squirrel.Eq{"s.room_id": *filter.RoomID})
	}
	if filter.MentorID != nil {
		q = q.Where(squirrel.Eq{"s.mentor_id": *filter.MentorID})
		countQ = countQ.Where(squirrel.Eq{"s.mentor_id": *filter.MentorID})
	}
	if filter.Status != nil {
		q = q.Where(squirrel.Eq{"s.status": *filter.Status})
		countQ = countQ.Where(squirrel.Eq{"s.status": *filter.Status})
	}
	if filter.Type != nil {
		q = q.Where(squirrel.Eq{"s.type": *filter.Type})
		countQ = countQ.Where(squirrel.Eq{"s.type": *filter.Type})
	}
	if filter.StartFrom != nil {
		q = q.Where(squirrel.GtOrEq{"s.start_at": *filter.StartFrom})
		countQ = countQ.Where(squirrel.GtOrEq{"s.start_at": *filter.StartFrom})
	}
	if filter.EndTo != nil {
		q = q.Where(squirrel.LtOrEq{"s.end_at": *filter.EndTo})
		countQ = countQ.Where(squirrel.LtOrEq{"s.end_at": *filter.EndTo})
	}

	query, args, err := q.OrderBy("s.start_at ASC").Limit(filter.Limit).Offset(filter.Offset).ToSql()
	if err != nil {
		return nil, 0, fmt.Errorf("build list slots query: %w", err)
	}

	rows, err := r.GetConn(ctx).Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("list slots: %w", err)
	}
	defer rows.Close()

	slots := make([]domain.Slot, 0)
	for rows.Next() {
		var slot domain.Slot
		if scanErr := rows.Scan(
			&slot.ID,
			&slot.MentorID,
			&slot.RoomID,
			&slot.CityID,
			&slot.Type,
			&slot.StartAt,
			&slot.EndAt,
			&slot.MeetingURL,
			&slot.Status,
			&slot.Capacity,
			&slot.CreatedBy,
			&slot.ApprovedBy,
			&slot.CreatedAt,
			&slot.UpdatedAt,
		); scanErr != nil {
			return nil, 0, fmt.Errorf("scan slot: %w", scanErr)
		}
		slots = append(slots, slot)
	}

	countQuery, countArgs, err := countQ.ToSql()
	if err != nil {
		return nil, 0, fmt.Errorf("build count slots query: %w", err)
	}
	var total uint64
	if err = r.GetConn(ctx).QueryRow(ctx, countQuery, countArgs...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count slots: %w", err)
	}

	return slots, total, nil
}

func (r *Repository) CountSlotConflicts(ctx context.Context, mentorID uuid.UUID, roomID *uuid.UUID, startAt, endAt time.Time, excludeSlotID *uuid.UUID) (int64, error) {
	cond := squirrel.And{
		squirrel.Expr("status IN (?, ?, ?)", domain.SlotStatusActive, domain.SlotStatusPending, domain.SlotStatusDraft),
		squirrel.Expr("start_at < ? AND end_at > ?", endAt, startAt),
		squirrel.Eq{"mentor_id": mentorID},
	}
	if excludeSlotID != nil {
		cond = append(cond, squirrel.NotEq{"id": *excludeSlotID})
	}

	query, args, err := psql.Select("COUNT(1)").From("slots").Where(cond).ToSql()
	if err != nil {
		return 0, fmt.Errorf("build mentor slot conflict query: %w", err)
	}
	var mentorConflicts int64
	if err = r.GetConn(ctx).QueryRow(ctx, query, args...).Scan(&mentorConflicts); err != nil {
		return 0, fmt.Errorf("count mentor slot conflict: %w", err)
	}

	if roomID == nil {
		return mentorConflicts, nil
	}

	roomCond := squirrel.And{
		squirrel.Expr("status IN (?, ?, ?)", domain.SlotStatusActive, domain.SlotStatusPending, domain.SlotStatusDraft),
		squirrel.Expr("start_at < ? AND end_at > ?", endAt, startAt),
		squirrel.Eq{"room_id": *roomID},
	}
	if excludeSlotID != nil {
		roomCond = append(roomCond, squirrel.NotEq{"id": *excludeSlotID})
	}
	roomQ, roomArgs, err := psql.Select("COUNT(1)").From("slots").Where(roomCond).ToSql()
	if err != nil {
		return 0, fmt.Errorf("build room slot conflict query: %w", err)
	}
	var roomConflicts int64
	if err = r.GetConn(ctx).QueryRow(ctx, roomQ, roomArgs...).Scan(&roomConflicts); err != nil {
		return 0, fmt.Errorf("count room slot conflict: %w", err)
	}

	return mentorConflicts + roomConflicts, nil
}
