package entities

import (
	"bootcamp-content-interaction-service/domains/posts/entities"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Likes struct {
	UserID    uuid.UUID      `gorm:"type:uuid;primaryKey"`
	PostId    uuid.UUID      `gorm:"type:uuid;not null"`
	Post      entities.Post  `gorm:"foreignKey:PostId;constraint:OnDelete:CASCADE"`
	CreatedAt time.Time      `gorm:"type:timestamp"`
	UpdatedAt time.Time      `gorm:"type:timestamp"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
