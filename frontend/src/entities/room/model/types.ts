export type Room = {
  id: string;
  hub_id: string;
  name: string;
  description?: string;
  room_type?: string;
  capacity: number;
  is_active: boolean;
};
