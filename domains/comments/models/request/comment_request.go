package request

import "github.com/google/uuid"

type CommentRequest struct {
	ReplyId *uuid.UUID `json:"reply_id,omitempty"`
	Msg     string     `json:"msg" validate:"required"`
}
