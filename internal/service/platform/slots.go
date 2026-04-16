package platform

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/linspacestrom/go-project/internal/domain"
	"github.com/linspacestrom/go-project/internal/repository/postgres"
)

func (s *Service) ListSlots(ctx context.Context, filter postgres.SlotFilter) ([]domain.Slot, uint64, error) {
	return s.repo.ListSlots(ctx, filter)
}

func (s *Service) GetSlotByID(ctx context.Context, slotID uuid.UUID) (*domain.Slot, error) {
	return s.repo.GetSlotByID(ctx, slotID)
}

func (s *Service) CreateSlot(ctx context.Context, actorID, mentorID uuid.UUID, roomID *uuid.UUID, cityID uuid.UUID, slotType string, startAt, endAt time.Time, meetingURL *string, status string, capacity *int) (*domain.Slot, error) {
	if !startAt.Before(endAt) {
		return nil, domain.ErrInvalidTime
	}
	if slotType == domain.SlotTypeOffline && roomID == nil {
		return nil, domain.ErrInvalidInput
	}
	if slotType == domain.SlotTypeOnline && roomID != nil {
		return nil, domain.ErrInvalidInput
	}
	if meetingURL == nil {
		return nil, domain.ErrInvalidInput
	}

	conflicts, err := s.repo.CountSlotConflicts(ctx, mentorID, roomID, startAt, endAt, nil)
	if err != nil {
		return nil, err
	}
	if conflicts > 0 {
		return nil, domain.ErrSlotTimeConflict
	}

	slot := &domain.Slot{
		MentorID:   mentorID,
		RoomID:     roomID,
		CityID:     cityID,
		Type:       slotType,
		StartAt:    startAt,
		EndAt:      endAt,
		MeetingURL: meetingURL,
		Status:     status,
		Capacity:   capacity,
		CreatedBy:  actorID,
	}
	created, err := s.repo.CreateSlot(ctx, slot)
	if err != nil {
		return nil, err
	}

	if err = s.publishEvent(ctx, "slot", &created.ID, domain.EventSlotCreated, map[string]any{"slot_id": created.ID, "mentor_id": created.MentorID}); err != nil {
		return nil, err
	}

	return created, nil
}

func (s *Service) UpdateSlot(ctx context.Context, slotID, actorID uuid.UUID, patch map[string]any) (*domain.Slot, error) {
	_ = actorID
	return s.repo.UpdateSlot(ctx, slotID, patch)
}

func (s *Service) DeleteSlot(ctx context.Context, slotID uuid.UUID) error {
	return s.repo.SoftDeleteSlot(ctx, slotID)
}
