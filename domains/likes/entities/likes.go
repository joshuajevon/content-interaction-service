package entities

import (
	post "bootcamp-content-interaction-service/domains/posts/entities"
	user "bootcamp-content-interaction-service/domains/users/entities"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Likes struct {
	ID		  uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	UserID    uuid.UUID      `gorm:"type:uuid"`
	User      user.User  	`gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	PostId    uuid.UUID      `gorm:"type:uuid;not null"`
	Post      post.Post  	`gorm:"foreignKey:PostId;constraint:OnDelete:CASCADE"`
	CreatedAt time.Time      `gorm:"type:timestamp"`
	UpdatedAt time.Time      `gorm:"type:timestamp"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
