export type BookingType = "room_only" | "room_with_mentor" | "mentor_call" | "event_seat";

export type Booking = {
  id: string;
  slot_id?: string;
  room_id?: string;
  user_id: string;
  booking_type: BookingType;
  status: string;
  start_at: string;
  end_at: string;
  meeting_url?: string;
  seat_number?: number;
  created_at: string;
  updated_at: string;
};
