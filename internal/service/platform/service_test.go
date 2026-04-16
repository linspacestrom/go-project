package platform

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/linspacestrom/go-project/internal/domain"
	"github.com/linspacestrom/go-project/internal/repository/postgres"
)

type fakeTRM struct{}

func (fakeTRM) Do(ctx context.Context, fn func(ctx context.Context) error) error { return fn(ctx) }

type fakeRepo struct {
	user                  *domain.User
	slot                  *domain.Slot
	room                  *domain.RoomView
	existingByIdempotency *domain.Booking
	bookingCount          int64
	createdBooking        *domain.Booking
}

func (f *fakeRepo) GetUserByID(context.Context, uuid.UUID) (*domain.User, error) { return f.user, nil }
func (f *fakeRepo) CheckUserExistsByEmail(context.Context, string) (bool, error) { return false, nil }
func (f *fakeRepo) CreateUser(context.Context, postgres.CreateUserParams) (*domain.User, error) {
	return nil, errors.New("not implemented")
}
func (f *fakeRepo) CreateMentorProfile(context.Context, uuid.UUID, *string, *string) error {
	return nil
}
func (f *fakeRepo) UpdateUserProfile(context.Context, uuid.UUID, *string, *time.Time) (*domain.User, error) {
	return nil, errors.New("not implemented")
}
func (f *fakeRepo) UpdateUserCity(context.Context, uuid.UUID, uuid.UUID) error { return nil }
func (f *fakeRepo) CreateCity(context.Context, string, bool) (*domain.City, error) {
	return nil, errors.New("not implemented")
}
func (f *fakeRepo) UpdateCity(context.Context, uuid.UUID, *string, *bool) (*domain.City, error) {
	return nil, errors.New("not implemented")
}
func (f *fakeRepo) ListCities(context.Context, uint64, uint64) ([]domain.City, uint64, error) {
	return nil, 0, errors.New("not implemented")
}
func (f *fakeRepo) CreateHub(context.Context, uuid.UUID, string, string, bool) (*domain.Hub, error) {
	return nil, errors.New("not implemented")
}
func (f *fakeRepo) UpdateHub(context.Context, uuid.UUID, *string, *string, *bool) (*domain.Hub, error) {
	return nil, errors.New("not implemented")
}
func (f *fakeRepo) ListHubsByCity(context.Context, uuid.UUID, uint64, uint64) ([]domain.Hub, uint64, error) {
	return nil, 0, errors.New("not implemented")
}
func (f *fakeRepo) CreateRoom(context.Context, uuid.UUID, string, *string, *string, int, bool) (*domain.Room, error) {
	return nil, errors.New("not implemented")
}
func (f *fakeRepo) UpdateRoom(context.Context, uuid.UUID, *string, *string, *string, *int, *bool) (*domain.Room, error) {
	return nil, errors.New("not implemented")
}
func (f *fakeRepo) DeleteRoom(context.Context, uuid.UUID) error { return errors.New("not implemented") }
func (f *fakeRepo) ListRoomsByCity(context.Context, uuid.UUID, uint64, uint64) ([]domain.Room, uint64, error) {
	return nil, 0, errors.New("not implemented")
}
func (f *fakeRepo) LockRoomByID(context.Context, uuid.UUID) (*domain.RoomView, error) {
	return f.room, nil
}
func (f *fakeRepo) CreateSlot(context.Context, *domain.Slot) (*domain.Slot, error) {
	return nil, errors.New("not implemented")
}
func (f *fakeRepo) UpdateSlot(context.Context, uuid.UUID, map[string]any) (*domain.Slot, error) {
	return nil, errors.New("not implemented")
}
func (f *fakeRepo) SoftDeleteSlot(context.Context, uuid.UUID) error {
	return errors.New("not implemented")
}
func (f *fakeRepo) LockSlotByID(context.Context, uuid.UUID) (*domain.Slot, error) { return f.slot, nil }
func (f *fakeRepo) GetSlotByID(context.Context, uuid.UUID) (*domain.Slot, error)  { return f.slot, nil }
func (f *fakeRepo) ListSlots(context.Context, postgres.SlotFilter) ([]domain.Slot, uint64, error) {
	return nil, 0, errors.New("not implemented")
}
func (f *fakeRepo) CountSlotConflicts(context.Context, uuid.UUID, *uuid.UUID, time.Time, time.Time, *uuid.UUID) (int64, error) {
	return 0, nil
}
func (f *fakeRepo) FindBookingByIdempotency(context.Context, uuid.UUID, string) (*domain.Booking, error) {
	if f.existingByIdempotency != nil {
		return f.existingByIdempotency, nil
	}
	return nil, domain.ErrBookingNotFound
}
func (f *fakeRepo) CountRoomBookingConflicts(context.Context, uuid.UUID, time.Time, time.Time) (int64, error) {
	return f.bookingCount, nil
}
func (f *fakeRepo) CountSlotActiveBookings(context.Context, uuid.UUID) (int64, error) { return 0, nil }
func (f *fakeRepo) CreateBooking(_ context.Context, booking *domain.Booking) (*domain.Booking, error) {
	copy := *booking
	copy.ID = uuid.New()
	copy.CreatedAt = time.Now().UTC()
	copy.UpdatedAt = copy.CreatedAt
	f.createdBooking = &copy
	return &copy, nil
}
func (f *fakeRepo) CancelBooking(context.Context, uuid.UUID, uuid.UUID, bool) (*domain.Booking, error) {
	return nil, errors.New("not implemented")
}
func (f *fakeRepo) ListBookings(context.Context, postgres.BookingFilter) ([]domain.Booking, uint64, error) {
	return nil, 0, errors.New("not implemented")
}
func (f *fakeRepo) CreateMentorRequest(context.Context, *domain.MentorRequest) (*domain.MentorRequest, error) {
	return nil, errors.New("not implemented")
}
func (f *fakeRepo) ListMentorRequests(context.Context, postgres.MentorRequestFilter) ([]domain.MentorRequest, uint64, error) {
	return nil, 0, errors.New("not implemented")
}
func (f *fakeRepo) UpdateMentorRequestStatus(context.Context, uuid.UUID, string, *uuid.UUID) (*domain.MentorRequest, error) {
	return nil, errors.New("not implemented")
}
func (f *fakeRepo) UpsertMentorSkillSubscription(context.Context, uuid.UUID, uuid.UUID) error {
	return nil
}
func (f *fakeRepo) DeleteMentorSkillSubscription(context.Context, uuid.UUID, uuid.UUID) error {
	return nil
}
func (f *fakeRepo) ListMentorSubscribersBySkill(context.Context, uuid.UUID, uuid.UUID) ([]uuid.UUID, error) {
	return nil, nil
}
func (f *fakeRepo) InsertAnalyticsEvent(context.Context, string, string, *uuid.UUID, []byte) error {
	return nil
}
func (f *fakeRepo) GetBusinessAnalytics(context.Context) (*struct {
	MentorApprovals                int64
	MentorRejections               int64
	ActiveBookings                 int64
	MentorRequestsBySkill          int64
	MentorRequestsByCategory       int64
	UserCancelsAfterMentorApproval int64
}, error) {
	return nil, errors.New("not implemented")
}
func (f *fakeRepo) GetTechnicalAnalytics(context.Context) (*struct {
	BookingConflictCount int64
	FailedOutboxEvents   int64
	OutboxLagCount       int64
}, error) {
	return nil, errors.New("not implemented")
}
func (f *fakeRepo) CreateOutboxEvent(context.Context, string, *uuid.UUID, string, []byte) error {
	return nil
}

func TestCreateBooking_Idempotency(t *testing.T) {
	userCity := uuid.New()
	repo := &fakeRepo{
		user:                  &domain.User{ID: uuid.New(), CityID: &userCity},
		existingByIdempotency: &domain.Booking{ID: uuid.New()},
	}
	svc := NewService(repo, fakeTRM{})

	booking, err := svc.CreateBooking(context.Background(), repo.user.ID, &domain.Booking{
		BookingType:    domain.BookingTypeRoomOnly,
		StartAt:        time.Now().UTC(),
		EndAt:          time.Now().UTC().Add(time.Hour),
		IdempotencyKey: ptr("idem-1"),
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if booking.ID != repo.existingByIdempotency.ID {
		t.Fatalf("expected existing booking id")
	}
}

func TestCreateBooking_CapacityConflict(t *testing.T) {
	cityID := uuid.New()
	roomID := uuid.New()
	repo := &fakeRepo{
		user:         &domain.User{ID: uuid.New(), CityID: &cityID},
		room:         &domain.RoomView{Room: domain.Room{ID: roomID, Capacity: 1}, CityID: cityID},
		bookingCount: 1,
	}
	svc := NewService(repo, fakeTRM{})

	_, err := svc.CreateBooking(context.Background(), repo.user.ID, &domain.Booking{
		RoomID:      &roomID,
		BookingType: domain.BookingTypeRoomOnly,
		StartAt:     time.Now().UTC(),
		EndAt:       time.Now().UTC().Add(time.Hour),
	})
	if !errors.Is(err, domain.ErrBookingCapacityExceeded) {
		t.Fatalf("expected ErrBookingCapacityExceeded, got %v", err)
	}
}

func ptr(v string) *string { return &v }
