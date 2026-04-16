export type Notification = {
  id: string;
  user_id: string;
  type: string;
  title: string;
  body: string;
  payload: Record<string, unknown>;
  is_alert: boolean;
  is_read: boolean;
  delivery_channel: "MAIL" | "IN_APP" | "WS" | "SSE" | string;
  status: "NEW" | "SENT" | "FAILED" | "READ" | string;
  created_at: string;
  updated_at: string;
};
