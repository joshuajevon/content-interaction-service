package response

import (
	"time"

	"github.com/google/uuid"
)

type CommentResponse struct{
	ID        uuid.UUID  `json:"id"`
	UserID    uuid.UUID  `json:"user_id"`
	ReplyId   *uuid.UUID `json:"reply_id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	Msg       string     `json:"msg"`
}