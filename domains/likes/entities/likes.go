package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Likes struct {
	UserID    uuid.UUID      `json:"user_id" gorm:"type:uuid;primaryKey"`
	PostId    uuid.UUID      `json:"post_id" gorm:"type:uuid;primaryKey"`
	CreatedAt time.Time      `json:"created_at" gorm:"type:timestamp"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"type:timestamp"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}
