package requests

type PostNotificationRequest struct {
	SourceUserID string `json:"source_user_id" validate:"required"`
    RecipientID  string `json:"recipient_id" validate:"required"`
    PostID       string `json:"post_id" validate:"required"`
    Content      string `json:"content"`
}