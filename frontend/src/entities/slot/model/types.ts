export type SlotType = "online" | "offline";

export type Slot = {
  id: string;
  mentor_id: string;
  room_id?: string;
  city_id: string;
  type: SlotType;
  start_at: string;
  end_at: string;
  meeting_url?: string;
  status: string;
  capacity?: number;
};
