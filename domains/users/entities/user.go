package entities

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey"`
	Name      string         `gorm:"type:varchar(255)"`
	Username  string         `gorm:"type:varchar(255);uniqueIndex"`
	Email     string         `gorm:"type:varchar(255);uniqueIndex"`
	Password  string         `gorm:"type:varchar(255)"`
	Bio       string         `gorm:"type:varchar(255)"`
	Gender    string         `gorm:"type:varchar(1)"`
	Phone     string         `gorm:"type:varchar(255)"`
	Country   string         `gorm:"type:varchar(255)"`
	Profile   string         `gorm:"type:varchar(255)"`
	CreatedAt time.Time      `gorm:"type:timestamp"`
	UpdatedAt time.Time      `gorm:"type:timestamp"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}