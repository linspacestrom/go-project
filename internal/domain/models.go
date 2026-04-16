package domain

import (
	"time"

	"github.com/google/uuid"
)

const (
	BookingStatusActive    = "active"
	BookingStatusCancelled = "canceled"
	BookingStatusCompleted = "completed"
	BookingStatusDeleted   = "deleted"
)

const (
	SlotTypeOnline  = "online"
	SlotTypeOffline = "offline"
)

const (
	SlotStatusDraft     = "draft"
	SlotStatusPending   = "pending"
	SlotStatusActive    = "active"
	SlotStatusCanceled  = "canceled"
	SlotStatusCompleted = "completed"
	SlotStatusDeleted   = "deleted"
)

const (
	BookingTypeRoomOnly       = "room_only"
	BookingTypeRoomWithMentor = "room_with_mentor"
	BookingTypeMentorCall     = "mentor_call"
	BookingTypeEventSeat      = "event_seat"
)

const (
	RequestTypeCategory     = "category"
	RequestTypeSkill        = "skill"
	RequestTypeDirectMentor = "direct_mentor"
	RequestTypeOther        = "other"
)

const (
	MentorRequestStatusPending  = "pending"
	MentorRequestStatusInReview = "in_review"
	MentorRequestStatusApproved = "approved"
	MentorRequestStatusRejected = "rejected"
	MentorRequestStatusCanceled = "canceled"
	MentorRequestStatusExpired  = "expired"
)

const (
	OutboxStatusNew     = "new"
	OutboxStatusSent    = "sent"
	OutboxStatusFailed  = "failed"
	OutboxStatusPending = "pending"
)

const (
	EventUserRegistered      = "user_registered"
	EventMentorRegistered    = "mentor_registered"
	EventSlotCreated         = "slot_created"
	EventBookingCreated      = "booking_created"
	EventBookingCanceled     = "booking_canceled"
	EventMentorRequestCreate = "mentor_request_created"
	EventMentorRequestUpdate = "mentor_request_updated"
	EventMentorSkillMatched  = "mentor_skill_matched"
	EventAdminCityChanged    = "admin_city_changed"
	EventAdminRoomCreated    = "admin_room_created"
	EventAdminHubCreated     = "admin_hub_created"
)

type User struct {
	ID        uuid.UUID
	FullName  string
	BirthDate *time.Time
	Email     string
	Role      string
	CityID    *uuid.UUID
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UserCredentials struct {
	User
	PasswordHash string
}

type StudentProfile struct {
	ID         uuid.UUID
	UserID     uuid.UUID
	University string
	Course     int
	DegreeType string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type MentorProfile struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	Description *string
	Title       *string
	IsVerified  bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type City struct {
	ID        uuid.UUID
	Name      string
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Hub struct {
	ID        uuid.UUID
	CityID    uuid.UUID
	Name      string
	Address   string
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Room struct {
	ID          uuid.UUID
	HubID       uuid.UUID
	Name        string
	Description *string
	RoomType    *string
	Capacity    int
	IsActive    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type WorkingHours struct {
	ID                uuid.UUID
	HubID             uuid.UUID
	DayOfWeek         int
	StartTime         string
	EndTime           string
	IsHolidayOverride bool
	HolidayDate       *time.Time
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

type Skill struct {
	ID          uuid.UUID
	Name        string
	Code        string
	Description *string
	IsActive    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type MentorSkillSubscription struct {
	ID        uuid.UUID
	MentorID  uuid.UUID
	SkillID   uuid.UUID
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Slot struct {
	ID         uuid.UUID
	MentorID   uuid.UUID
	RoomID     *uuid.UUID
	CityID     uuid.UUID
	Type       string
	StartAt    time.Time
	EndAt      time.Time
	MeetingURL *string
	Status     string
	Capacity   *int
	CreatedBy  uuid.UUID
	ApprovedBy *uuid.UUID
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type Booking struct {
	ID             uuid.UUID
	SlotID         *uuid.UUID
	RoomID         *uuid.UUID
	UserID         uuid.UUID
	BookingType    string
	Status         string
	StartAt        time.Time
	EndAt          time.Time
	MeetingURL     *string
	SeatNumber     *int
	IdempotencyKey *string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type MentorRequest struct {
	ID          uuid.UUID
	SlotID      *uuid.UUID
	MentorID    *uuid.UUID
	MenteeID    uuid.UUID
	CityID      uuid.UUID
	SkillID     *uuid.UUID
	RequestType string
	Status      string
	Comment     *string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type AnalyticsEvent struct {
	ID         uuid.UUID
	EventType  string
	EntityType string
	EntityID   *uuid.UUID
	Payload    []byte
	CreatedAt  time.Time
}

type OutboxEvent struct {
	ID            uuid.UUID
	AggregateType string
	AggregateID   *uuid.UUID
	EventType     string
	Payload       []byte
	Status        string
	ErrorMessage  *string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// Lightweight projection for room listing with hub/city constraints.
type RoomView struct {
	Room
	HubName   string
	HubIDView uuid.UUID
	CityID    uuid.UUID
}
