package platform

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/linspacestrom/go-project/internal/domain"
	"github.com/linspacestrom/go-project/internal/repository/postgres"
)

type TransactionManager interface {
	Do(ctx context.Context, fn func(ctx context.Context) error) error
}

type Repository interface {
	GetUserByID(ctx context.Context, userID uuid.UUID) (*domain.User, error)
	CheckUserExistsByEmail(ctx context.Context, email string) (bool, error)
	CreateUser(ctx context.Context, params postgres.CreateUserParams) (*domain.User, error)
	CreateMentorProfile(ctx context.Context, userID uuid.UUID, description, title *string) error
	UpdateUserProfile(ctx context.Context, userID uuid.UUID, fullName *string, birthDate *time.Time) (*domain.User, error)
	UpdateUserCity(ctx context.Context, userID, cityID uuid.UUID) error

	CreateCity(ctx context.Context, name string, isActive bool) (*domain.City, error)
	UpdateCity(ctx context.Context, cityID uuid.UUID, name *string, isActive *bool) (*domain.City, error)
	ListCities(ctx context.Context, limit, offset uint64) ([]domain.City, uint64, error)

	CreateHub(ctx context.Context, cityID uuid.UUID, name, address string, isActive bool) (*domain.Hub, error)
	UpdateHub(ctx context.Context, hubID uuid.UUID, name, address *string, isActive *bool) (*domain.Hub, error)
	ListHubsByCity(ctx context.Context, cityID uuid.UUID, limit, offset uint64) ([]domain.Hub, uint64, error)

	CreateRoom(ctx context.Context, hubID uuid.UUID, name string, description, roomType *string, capacity int, isActive bool) (*domain.Room, error)
	UpdateRoom(ctx context.Context, roomID uuid.UUID, name, description, roomType *string, capacity *int, isActive *bool) (*domain.Room, error)
	DeleteRoom(ctx context.Context, roomID uuid.UUID) error
	ListRoomsByCity(ctx context.Context, cityID uuid.UUID, limit, offset uint64) ([]domain.Room, uint64, error)
	LockRoomByID(ctx context.Context, roomID uuid.UUID) (*domain.RoomView, error)

	CreateSlot(ctx context.Context, slot *domain.Slot) (*domain.Slot, error)
	UpdateSlot(ctx context.Context, slotID uuid.UUID, patch map[string]any) (*domain.Slot, error)
	SoftDeleteSlot(ctx context.Context, slotID uuid.UUID) error
	LockSlotByID(ctx context.Context, slotID uuid.UUID) (*domain.Slot, error)
	GetSlotByID(ctx context.Context, slotID uuid.UUID) (*domain.Slot, error)
	ListSlots(ctx context.Context, filter postgres.SlotFilter) ([]domain.Slot, uint64, error)
	CountSlotConflicts(ctx context.Context, mentorID uuid.UUID, roomID *uuid.UUID, startAt, endAt time.Time, excludeSlotID *uuid.UUID) (int64, error)

	FindBookingByIdempotency(ctx context.Context, userID uuid.UUID, key string) (*domain.Booking, error)
	CountRoomBookingConflicts(ctx context.Context, roomID uuid.UUID, startAt, endAt time.Time) (int64, error)
	CountSlotActiveBookings(ctx context.Context, slotID uuid.UUID) (int64, error)
	CreateBooking(ctx context.Context, booking *domain.Booking) (*domain.Booking, error)
	CancelBooking(ctx context.Context, bookingID, userID uuid.UUID, force bool) (*domain.Booking, error)
	ListBookings(ctx context.Context, filter postgres.BookingFilter) ([]domain.Booking, uint64, error)

	CreateMentorRequest(ctx context.Context, req *domain.MentorRequest) (*domain.MentorRequest, error)
	ListMentorRequests(ctx context.Context, filter postgres.MentorRequestFilter) ([]domain.MentorRequest, uint64, error)
	UpdateMentorRequestStatus(ctx context.Context, requestID uuid.UUID, status string, mentorID *uuid.UUID) (*domain.MentorRequest, error)
	UpsertMentorSkillSubscription(ctx context.Context, mentorID, skillID uuid.UUID) error
	DeleteMentorSkillSubscription(ctx context.Context, mentorID, skillID uuid.UUID) error
	ListMentorSubscribersBySkill(ctx context.Context, skillID uuid.UUID, cityID uuid.UUID) ([]uuid.UUID, error)

	InsertAnalyticsEvent(ctx context.Context, eventType, entityType string, entityID *uuid.UUID, payload []byte) error
	GetBusinessAnalytics(ctx context.Context) (*struct {
		MentorApprovals                int64
		MentorRejections               int64
		ActiveBookings                 int64
		MentorRequestsBySkill          int64
		MentorRequestsByCategory       int64
		UserCancelsAfterMentorApproval int64
	}, error)
	GetTechnicalAnalytics(ctx context.Context) (*struct {
		BookingConflictCount int64
		FailedOutboxEvents   int64
		OutboxLagCount       int64
	}, error)

	CreateOutboxEvent(ctx context.Context, aggregateType string, aggregateID *uuid.UUID, eventType string, payload []byte) error
}

type Service struct {
	repo Repository
	trm  TransactionManager
}

func NewService(repo Repository, trm TransactionManager) *Service {
	return &Service{repo: repo, trm: trm}
}

func (s *Service) publishEvent(ctx context.Context, aggregateType string, aggregateID *uuid.UUID, eventType string, payload map[string]any) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal event payload: %w", err)
	}

	return s.repo.CreateOutboxEvent(ctx, aggregateType, aggregateID, eventType, body)
}
