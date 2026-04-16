package dto

type NotificationResponse struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id"`
	Type      string `json:"type"`
	Title     string `json:"title"`
	Body      string `json:"body"`
	IsAlert   bool   `json:"is_alert"`
	IsRead    bool   `json:"is_read"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
}
