package platform

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/linspacestrom/go-project/internal/domain"
	"github.com/linspacestrom/go-project/internal/repository/postgres"
)

func (s *Service) CreateBooking(ctx context.Context, userID uuid.UUID, booking *domain.Booking) (*domain.Booking, error) {
	booking.UserID = userID
	booking.Status = domain.BookingStatusActive

	if !booking.StartAt.Before(booking.EndAt) {
		return nil, domain.ErrInvalidTime
	}
	if booking.IdempotencyKey != nil {
		existing, err := s.repo.FindBookingByIdempotency(ctx, userID, *booking.IdempotencyKey)
		if err == nil {
			return existing, nil
		}
		if !errors.Is(err, domain.ErrBookingNotFound) {
			return nil, err
		}
	}

	var created *domain.Booking
	if err := s.trm.Do(ctx, func(txCtx context.Context) error {
		user, err := s.repo.GetUserByID(txCtx, userID)
		if err != nil {
			return err
		}
		if user.CityID == nil {
			return domain.ErrBookingCityMismatch
		}

		if booking.SlotID != nil {
			slot, slotErr := s.repo.LockSlotByID(txCtx, *booking.SlotID)
			if slotErr != nil {
				return slotErr
			}
			if slot.Status != domain.SlotStatusActive {
				return domain.ErrSlotNotFound
			}
			if slot.CityID != *user.CityID {
				return domain.ErrBookingCityMismatch
			}
			currentBookings, cntErr := s.repo.CountSlotActiveBookings(txCtx, slot.ID)
			if cntErr != nil {
				return cntErr
			}
			if slot.Capacity != nil && currentBookings >= int64(*slot.Capacity) {
				return domain.ErrBookingCapacityExceeded
			}
			if slot.RoomID != nil {
				booking.RoomID = slot.RoomID
			}
		}

		if booking.RoomID != nil {
			room, roomErr := s.repo.LockRoomByID(txCtx, *booking.RoomID)
			if roomErr != nil {
				return roomErr
			}
			if room.CityID != *user.CityID {
				return domain.ErrBookingCityMismatch
			}

			conflicts, conflictErr := s.repo.CountRoomBookingConflicts(txCtx, room.ID, booking.StartAt, booking.EndAt)
			if conflictErr != nil {
				return conflictErr
			}
			if conflicts >= int64(room.Capacity) {
				_ = s.repo.InsertAnalyticsEvent(txCtx, "booking_conflict", "booking", nil, []byte(`{"reason":"capacity"}`))
				return domain.ErrBookingCapacityExceeded
			}
		}

		createdBooking, createErr := s.repo.CreateBooking(txCtx, booking)
		if createErr != nil {
			return createErr
		}
		created = createdBooking

		analyticsPayload, _ := json.Marshal(map[string]any{"booking_id": created.ID, "type": created.BookingType})
		if analyticsErr := s.repo.InsertAnalyticsEvent(txCtx, "booking_created", "booking", &created.ID, analyticsPayload); analyticsErr != nil {
			return analyticsErr
		}

		if publishErr := s.publishEvent(txCtx, "booking", &created.ID, domain.EventBookingCreated, map[string]any{"booking_id": created.ID, "user_id": userID}); publishErr != nil {
			return publishErr
		}

		return nil
	}); err != nil {
		if errors.Is(err, domain.ErrBookingCapacityExceeded) || errors.Is(err, domain.ErrBookingCityMismatch) {
			return nil, err
		}
		return nil, fmt.Errorf("transaction create booking: %w", err)
	}

	return created, nil
}

func (s *Service) CancelBooking(ctx context.Context, bookingID, userID uuid.UUID, force bool) (*domain.Booking, error) {
	var canceled *domain.Booking
	if err := s.trm.Do(ctx, func(txCtx context.Context) error {
		b, err := s.repo.CancelBooking(txCtx, bookingID, userID, force)
		if err != nil {
			return err
		}
		canceled = b

		payload, _ := json.Marshal(map[string]any{"booking_id": b.ID, "status": b.Status})
		if analyticsErr := s.repo.InsertAnalyticsEvent(txCtx, "booking_canceled", "booking", &b.ID, payload); analyticsErr != nil {
			return analyticsErr
		}
		if b.BookingType == domain.BookingTypeRoomWithMentor {
			_ = s.repo.InsertAnalyticsEvent(txCtx, "booking_cancel_after_mentor_approval", "booking", &b.ID, payload)
		}
		if publishErr := s.publishEvent(txCtx, "booking", &b.ID, domain.EventBookingCanceled, map[string]any{"booking_id": b.ID, "user_id": b.UserID}); publishErr != nil {
			return publishErr
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return canceled, nil
}

func (s *Service) ListBookings(ctx context.Context, filter postgres.BookingFilter) ([]domain.Booking, uint64, error) {
	return s.repo.ListBookings(ctx, filter)
}
