package platform

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/linspacestrom/go-project/internal/domain"
	"github.com/linspacestrom/go-project/internal/dto"
	"github.com/linspacestrom/go-project/internal/mapper"
	"github.com/linspacestrom/go-project/internal/repository/postgres"
	"github.com/linspacestrom/go-project/internal/server/middleware"
)

type PlatformService interface {
	GetMe(ctx context.Context, userID uuid.UUID) (*domain.User, error)
	UpdateMe(ctx context.Context, userID uuid.UUID, fullName *string, birthDate *time.Time) (*domain.User, error)
	RegisterMentor(ctx context.Context, email, password, fullName string, cityID uuid.UUID, description, title *string) (*domain.User, error)

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
	UpdateUserCity(ctx context.Context, userID, cityID uuid.UUID) error

	ListSlots(ctx context.Context, filter postgres.SlotFilter) ([]domain.Slot, uint64, error)
	GetSlotByID(ctx context.Context, slotID uuid.UUID) (*domain.Slot, error)
	CreateSlot(ctx context.Context, actorID, mentorID uuid.UUID, roomID *uuid.UUID, cityID uuid.UUID, slotType string, startAt, endAt time.Time, meetingURL *string, status string, capacity *int) (*domain.Slot, error)
	UpdateSlot(ctx context.Context, slotID, actorID uuid.UUID, patch map[string]any) (*domain.Slot, error)
	DeleteSlot(ctx context.Context, slotID uuid.UUID) error

	CreateBooking(ctx context.Context, userID uuid.UUID, booking *domain.Booking) (*domain.Booking, error)
	CancelBooking(ctx context.Context, bookingID, userID uuid.UUID, force bool) (*domain.Booking, error)
	ListBookings(ctx context.Context, filter postgres.BookingFilter) ([]domain.Booking, uint64, error)

	CreateMentorRequest(ctx context.Context, req *domain.MentorRequest) (*domain.MentorRequest, error)
	ListMentorRequests(ctx context.Context, filter postgres.MentorRequestFilter) ([]domain.MentorRequest, uint64, error)
	UpdateMentorRequestStatus(ctx context.Context, requestID uuid.UUID, status string, mentorID *uuid.UUID) (*domain.MentorRequest, error)

	SubscribeMentorSkill(ctx context.Context, mentorID, skillID uuid.UUID) error
	UnsubscribeMentorSkill(ctx context.Context, mentorID, skillID uuid.UUID) error

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
}

type Handler struct {
	s PlatformService
}

func NewHandler(s PlatformService) *Handler {
	return &Handler{s: s}
}

func (h *Handler) RegisterRoutes(r gin.IRouter) {
	authGroup := r.Group("/api/v1")
	authGroup.GET("/me", h.GetMe)
	authGroup.PATCH("/me", h.UpdateMe)
	authGroup.GET("/cities", h.ListCities)
	authGroup.GET("/cities/:cityId/hubs", h.ListHubsByCity)
	authGroup.GET("/cities/:cityId/rooms", h.ListRoomsByCity)
	authGroup.GET("/slots", h.ListSlots)
	authGroup.GET("/slots/:id", h.GetSlotByID)
	authGroup.POST("/bookings", h.CreateBooking)
	authGroup.DELETE("/bookings/:id", h.CancelBooking)
	authGroup.GET("/bookings", h.ListBookings)
	authGroup.POST("/mentor-requests", h.CreateMentorRequest)
	authGroup.GET("/mentor-requests", h.ListMentorRequests)

	admin := authGroup.Group("/admin")
	admin.Use(middleware.RequireRoles(domain.RoleAdmin))
	admin.POST("/mentors/register_mentor", h.RegisterMentor)
	admin.POST("/cities", h.CreateCity)
	admin.PATCH("/cities/:id", h.UpdateCity)
	admin.POST("/hubs", h.CreateHub)
	admin.PATCH("/hubs/:id", h.UpdateHub)
	admin.POST("/rooms", h.CreateRoom)
	admin.PATCH("/rooms/:id", h.UpdateRoom)
	admin.DELETE("/rooms/:id", h.DeleteRoom)
	admin.PATCH("/users/:id/city", h.UpdateUserCity)
	admin.GET("/analytics/business", h.BusinessAnalytics)
	admin.GET("/analytics/technical", h.TechnicalAnalytics)

	mentor := authGroup.Group("/mentor")
	mentor.Use(middleware.RequireRoles(domain.RoleMentor))
	mentor.POST("/skills/subscribe", h.SubscribeMentorSkill)
	mentor.DELETE("/skills/subscribe/:skillId", h.UnsubscribeMentorSkill)
	mentor.GET("/slots", h.ListMentorSlots)
	mentor.POST("/slots", h.CreateMentorSlot)
	mentor.PATCH("/slots/:id", h.UpdateMentorSlot)
	mentor.DELETE("/slots/:id", h.DeleteMentorSlot)
	mentor.GET("/requests", h.ListMentorRequests)
	mentor.POST("/requests/:id/approve", h.ApproveMentorRequest)
	mentor.POST("/requests/:id/reject", h.RejectMentorRequest)
}

func (h *Handler) mapError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, domain.ErrInvalidInput), errors.Is(err, domain.ErrInvalidTime):
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: dto.ErrorDetail{Code: "INVALID_REQUEST", Message: err.Error()}})
	case errors.Is(err, domain.ErrUnauthorized):
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: dto.ErrorDetail{Code: "UNAUTHORIZED", Message: err.Error()}})
	case errors.Is(err, domain.ErrForbidden):
		c.JSON(http.StatusForbidden, dto.ErrorResponse{Error: dto.ErrorDetail{Code: "FORBIDDEN", Message: err.Error()}})
	case errors.Is(err, domain.ErrRoomNotFound), errors.Is(err, domain.ErrSlotNotFound), errors.Is(err, domain.ErrBookingNotFound), errors.Is(err, domain.ErrMentorRequestNotFound), errors.Is(err, domain.ErrCityNotFound), errors.Is(err, domain.ErrHubNotFound):
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: dto.ErrorDetail{Code: "NOT_FOUND", Message: err.Error()}})
	case errors.Is(err, domain.ErrConflict), errors.Is(err, domain.ErrBookingCapacityExceeded), errors.Is(err, domain.ErrBookingCityMismatch), errors.Is(err, domain.ErrSlotTimeConflict):
		c.JSON(http.StatusConflict, dto.ErrorResponse{Error: dto.ErrorDetail{Code: "CONFLICT", Message: err.Error()}})
	default:
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: dto.ErrorDetail{Code: "INTERNAL_ERROR", Message: "internal server error"}})
	}
}

func ptrOrNilUUID(raw string) (*uuid.UUID, error) {
	if raw == "" {
		return nil, nil
	}
	parsed, err := uuid.Parse(raw)
	if err != nil {
		return nil, err
	}

	return &parsed, nil
}

func firstUserID(c *gin.Context) (uuid.UUID, error) {
	return middleware.GetUserID(c)
}

func toCityListResponse(cities []domain.City, total, limit, offset uint64) dto.ListResponse[dto.CityResponse] {
	items := make([]dto.CityResponse, 0, len(cities))
	for _, city := range cities {
		copy := city
		items = append(items, mapper.ToCityResponse(&copy))
	}

	return dto.ListResponse[dto.CityResponse]{Items: items, TotalCount: total, Limit: limit, Offset: offset}
}
