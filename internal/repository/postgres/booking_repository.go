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

func (r *Repository) FindBookingByIdempotency(ctx context.Context, userID uuid.UUID, key string) (*domain.Booking, error) {
	query, args, err := psql.Select("id", "slot_id", "room_id", "user_id", "booking_type", "status", "start_at", "end_at", "meeting_url", "seat_number", "idempotency_key", "created_at", "updated_at").
		From("bookings").
		Where(squirrel.Eq{"user_id": userID, "idempotency_key": key}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build idempotency booking query: %w", err)
	}

	var b domain.Booking
	err = r.GetConn(ctx).QueryRow(ctx, query, args...).Scan(
		&b.ID,
		&b.SlotID,
		&b.RoomID,
		&b.UserID,
		&b.BookingType,
		&b.Status,
		&b.StartAt,
		&b.EndAt,
		&b.MeetingURL,
		&b.SeatNumber,
		&b.IdempotencyKey,
		&b.CreatedAt,
		&b.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrBookingNotFound
		}
		return nil, fmt.Errorf("get booking by idempotency: %w", err)
	}

	return &b, nil
}

func (r *Repository) CountRoomBookingConflicts(ctx context.Context, roomID uuid.UUID, startAt, endAt time.Time) (int64, error) {
	query, args, err := psql.Select("COUNT(1)").
		From("bookings").
		Where(squirrel.And{
			squirrel.Eq{"room_id": roomID},
			squirrel.Eq{"status": domain.BookingStatusActive},
			squirrel.Expr("start_at < ? AND end_at > ?", endAt, startAt),
		}).
		ToSql()
	if err != nil {
		return 0, fmt.Errorf("build room booking conflict query: %w", err)
	}
	var conflicts int64
	if err = r.GetConn(ctx).QueryRow(ctx, query, args...).Scan(&conflicts); err != nil {
		return 0, fmt.Errorf("count room booking conflict: %w", err)
	}

	return conflicts, nil
}

func (r *Repository) CountSlotActiveBookings(ctx context.Context, slotID uuid.UUID) (int64, error) {
	query, args, err := psql.Select("COUNT(1)").
		From("bookings").
		Where(squirrel.Eq{"slot_id": slotID, "status": domain.BookingStatusActive}).
		ToSql()
	if err != nil {
		return 0, fmt.Errorf("build count slot bookings query: %w", err)
	}
	var cnt int64
	if err = r.GetConn(ctx).QueryRow(ctx, query, args...).Scan(&cnt); err != nil {
		return 0, fmt.Errorf("count slot bookings: %w", err)
	}

	return cnt, nil
}

func (r *Repository) CreateBooking(ctx context.Context, booking *domain.Booking) (*domain.Booking, error) {
	query, args, err := psql.Insert("bookings").
		Columns("slot_id", "room_id", "user_id", "booking_type", "status", "start_at", "end_at", "meeting_url", "seat_number", "idempotency_key").
		Values(booking.SlotID, booking.RoomID, booking.UserID, booking.BookingType, booking.Status, booking.StartAt, booking.EndAt, booking.MeetingURL, booking.SeatNumber, booking.IdempotencyKey).
		Suffix("RETURNING id, created_at, updated_at").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build create booking query: %w", err)
	}

	created := *booking
	if err = r.GetConn(ctx).QueryRow(ctx, query, args...).Scan(&created.ID, &created.CreatedAt, &created.UpdatedAt); err != nil {
		return nil, fmt.Errorf("create booking: %w", err)
	}

	return &created, nil
}

func (r *Repository) CancelBooking(ctx context.Context, bookingID, userID uuid.UUID, force bool) (*domain.Booking, error) {
	where := squirrel.Eq{"id": bookingID}
	if !force {
		where["user_id"] = userID
	}
	query, args, err := psql.Update("bookings").
		Set("status", domain.BookingStatusCancelled).
		Set("updated_at", squirrel.Expr("now()")).
		Where(where).
		Suffix("RETURNING id, slot_id, room_id, user_id, booking_type, status, start_at, end_at, meeting_url, seat_number, idempotency_key, created_at, updated_at").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build cancel booking query: %w", err)
	}

	var b domain.Booking
	if err = r.GetConn(ctx).QueryRow(ctx, query, args...).Scan(
		&b.ID,
		&b.SlotID,
		&b.RoomID,
		&b.UserID,
		&b.BookingType,
		&b.Status,
		&b.StartAt,
		&b.EndAt,
		&b.MeetingURL,
		&b.SeatNumber,
		&b.IdempotencyKey,
		&b.CreatedAt,
		&b.UpdatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrBookingNotFound
		}
		return nil, fmt.Errorf("cancel booking: %w", err)
	}

	return &b, nil
}

func (r *Repository) ListBookings(ctx context.Context, filter BookingFilter) ([]domain.Booking, uint64, error) {
	q := psql.Select("id", "slot_id", "room_id", "user_id", "booking_type", "status", "start_at", "end_at", "meeting_url", "seat_number", "idempotency_key", "created_at", "updated_at").
		From("bookings")
	countQ := psql.Select("COUNT(1)").From("bookings")

	if filter.UserID != nil {
		q = q.Where(squirrel.Eq{"user_id": *filter.UserID})
		countQ = countQ.Where(squirrel.Eq{"user_id": *filter.UserID})
	}
	if filter.Status != nil {
		q = q.Where(squirrel.Eq{"status": *filter.Status})
		countQ = countQ.Where(squirrel.Eq{"status": *filter.Status})
	}
	if filter.CityID != nil {
		q = q.Join("slots s ON s.id = bookings.slot_id").Where(squirrel.Eq{"s.city_id": *filter.CityID})
		countQ = countQ.Join("slots s ON s.id = bookings.slot_id").Where(squirrel.Eq{"s.city_id": *filter.CityID})
	}
	if filter.RoomID != nil {
		q = q.Where(squirrel.Eq{"room_id": *filter.RoomID})
		countQ = countQ.Where(squirrel.Eq{"room_id": *filter.RoomID})
	}
	if filter.SlotID != nil {
		q = q.Where(squirrel.Eq{"slot_id": *filter.SlotID})
		countQ = countQ.Where(squirrel.Eq{"slot_id": *filter.SlotID})
	}

	query, args, err := q.OrderBy("start_at DESC").Limit(filter.Limit).Offset(filter.Offset).ToSql()
	if err != nil {
		return nil, 0, fmt.Errorf("build list bookings query: %w", err)
	}
	rows, err := r.GetConn(ctx).Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("list bookings: %w", err)
	}
	defer rows.Close()

	bookings := make([]domain.Booking, 0)
	for rows.Next() {
		var b domain.Booking
		if scanErr := rows.Scan(
			&b.ID,
			&b.SlotID,
			&b.RoomID,
			&b.UserID,
			&b.BookingType,
			&b.Status,
			&b.StartAt,
			&b.EndAt,
			&b.MeetingURL,
			&b.SeatNumber,
			&b.IdempotencyKey,
			&b.CreatedAt,
			&b.UpdatedAt,
		); scanErr != nil {
			return nil, 0, fmt.Errorf("scan booking: %w", scanErr)
		}
		bookings = append(bookings, b)
	}

	countQuery, countArgs, err := countQ.ToSql()
	if err != nil {
		return nil, 0, fmt.Errorf("build count bookings query: %w", err)
	}
	var total uint64
	if err = r.GetConn(ctx).QueryRow(ctx, countQuery, countArgs...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count bookings: %w", err)
	}

	return bookings, total, nil
}
