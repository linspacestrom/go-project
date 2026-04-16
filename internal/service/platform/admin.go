package platform

import (
	"context"

	"github.com/google/uuid"
	"github.com/linspacestrom/go-project/internal/domain"
)

func (s *Service) CreateCity(ctx context.Context, name string, isActive bool) (*domain.City, error) {
	return s.repo.CreateCity(ctx, name, isActive)
}

func (s *Service) UpdateCity(ctx context.Context, cityID uuid.UUID, name *string, isActive *bool) (*domain.City, error) {
	return s.repo.UpdateCity(ctx, cityID, name, isActive)
}

func (s *Service) ListCities(ctx context.Context, limit, offset uint64) ([]domain.City, uint64, error) {
	return s.repo.ListCities(ctx, limit, offset)
}

func (s *Service) CreateHub(ctx context.Context, cityID uuid.UUID, name, address string, isActive bool) (*domain.Hub, error) {
	hub, err := s.repo.CreateHub(ctx, cityID, name, address, isActive)
	if err != nil {
		return nil, err
	}
	if outErr := s.publishEvent(ctx, "hub", &hub.ID, domain.EventAdminHubCreated, map[string]any{"hub_id": hub.ID, "city_id": hub.CityID}); outErr != nil {
		return nil, outErr
	}

	return hub, nil
}

func (s *Service) UpdateHub(ctx context.Context, hubID uuid.UUID, name, address *string, isActive *bool) (*domain.Hub, error) {
	return s.repo.UpdateHub(ctx, hubID, name, address, isActive)
}

func (s *Service) ListHubsByCity(ctx context.Context, cityID uuid.UUID, limit, offset uint64) ([]domain.Hub, uint64, error) {
	return s.repo.ListHubsByCity(ctx, cityID, limit, offset)
}

func (s *Service) CreateRoom(ctx context.Context, hubID uuid.UUID, name string, description, roomType *string, capacity int, isActive bool) (*domain.Room, error) {
	room, err := s.repo.CreateRoom(ctx, hubID, name, description, roomType, capacity, isActive)
	if err != nil {
		return nil, err
	}
	if outErr := s.publishEvent(ctx, "room", &room.ID, domain.EventAdminRoomCreated, map[string]any{"room_id": room.ID, "hub_id": room.HubID}); outErr != nil {
		return nil, outErr
	}

	return room, nil
}

func (s *Service) UpdateRoom(ctx context.Context, roomID uuid.UUID, name, description, roomType *string, capacity *int, isActive *bool) (*domain.Room, error) {
	return s.repo.UpdateRoom(ctx, roomID, name, description, roomType, capacity, isActive)
}

func (s *Service) DeleteRoom(ctx context.Context, roomID uuid.UUID) error {
	return s.repo.DeleteRoom(ctx, roomID)
}

func (s *Service) ListRoomsByCity(ctx context.Context, cityID uuid.UUID, limit, offset uint64) ([]domain.Room, uint64, error) {
	return s.repo.ListRoomsByCity(ctx, cityID, limit, offset)
}
