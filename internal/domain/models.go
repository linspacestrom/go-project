package domain

import (
	"time"

	"github.com/google/uuid"
)

const (
	BookingStatusActive    = "active"
	BookingStatusCancelled = "cancelled"
)

type User struct {
	ID        uuid.UUID
	Email     string
	Role      string
	CreatedAt time.Time
}

type Room struct {
	ID          uuid.UUID
	Name        string
	Description *string
	Capacity    *int
	CreatedAt   time.Time
}

type Schedule struct {
	ID                 uuid.UUID
	RoomID             uuid.UUID
	DaysOfWeek         []int
	StartTime          string
	EndTime            string
	LastGeneratedUntil *time.Time
}

type Slot struct {
	ID     uuid.UUID
	RoomID uuid.UUID
	Start  time.Time
	End    time.Time
}

type Booking struct {
	ID             uuid.UUID
	SlotID         uuid.UUID
	UserID         uuid.UUID
	Status         string
	ConferenceLink *string
	CreatedAt      time.Time
}

type UserCredentials struct {
	User
	PasswordHash string
}
