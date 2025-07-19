package entities

import (
	"time"

	"github.com/google/uuid"
)

type Comments struct {
	ID        uuid.UUID  `json:"id" gorm:"type:uuid;primaryKey"`
	UserID    uuid.UUID  `json:"user_id" gorm:"type:uuid"`
	PostId    uuid.UUID  `json:"post_id" gorm:"type:uuid"`
	ReplyId   *uuid.UUID `json:"reply_id" gorm:"type:uuid"`
	CreatedAt time.Time  `json:"created_at" gorm:"type:timestamp"`
	UpdatedAt time.Time  `json:"updated_at" gorm:"type:timestamp"`
	Msg       string     `json:"msg" gorm:"type:string"`
}
