package domain

import "errors"

var (
	ErrRoomNotFound       = errors.New("room not found")
	ErrSlotNotFound       = errors.New("slot not found")
	ErrSlotAlreadyBooked  = errors.New("slot is already booked")
	ErrBookingNotFound    = errors.New("booking not found")
	ErrScheduleExists     = errors.New("schedule for this room already exists")
	ErrInvalidTime        = errors.New("invalid time provided")
	ErrForbidden          = errors.New("forbidden")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrEmailExists        = errors.New("email already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
)
