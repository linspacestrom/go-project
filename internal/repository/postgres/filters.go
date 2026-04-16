package postgres

import (
	"time"

	"github.com/google/uuid"
)

type SlotFilter struct {
	CityID    *uuid.UUID
	HubID     *uuid.UUID
	RoomID    *uuid.UUID
	MentorID  *uuid.UUID
	Status    *string
	Type      *string
	StartFrom *time.Time
	EndTo     *time.Time
	Limit     uint64
	Offset    uint64
}

type BookingFilter struct {
	UserID *uuid.UUID
	Status *string
	CityID *uuid.UUID
	RoomID *uuid.UUID
	SlotID *uuid.UUID
	Limit  uint64
	Offset uint64
}

type MentorRequestFilter struct {
	MenteeID *uuid.UUID
	MentorID *uuid.UUID
	SkillID  *uuid.UUID
	CityID   *uuid.UUID
	Status   *string
	Limit    uint64
	Offset   uint64
}
