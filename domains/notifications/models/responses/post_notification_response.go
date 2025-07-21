package responses

type PostNotificationResponse struct {
	ID			 string `json:"id"`
	SourceUserID string `json:"source_user_id"`
    RecipientID  string `json:"recipient_id"`
    PostID       string `json:"post_id"`
    Content      string `json:"content"`
}