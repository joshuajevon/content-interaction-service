package entities

import (
	"bootcamp-content-interaction-service/domains/users/entities"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type Post struct {
    ID         uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
    UserID     uuid.UUID      `gorm:"type:uuid;not null"`
	User       entities.User  `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
    ImageURLs  pq.StringArray `gorm:"type:text[]"`
    Caption    string         `gorm:"type:text"`
    Tags       pq.StringArray `gorm:"type:text[]"`
    CreatedAt  time.Time      `gorm:"type:timestamp"`
    UpdatedAt  time.Time      `gorm:"type:timestamp"`
}