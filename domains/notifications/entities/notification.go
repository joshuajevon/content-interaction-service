package entities

import (
	posts "bootcamp-content-interaction-service/domains/posts/entities"
	users "bootcamp-content-interaction-service/domains/users/entities"
	"time"

	"github.com/google/uuid"
)

type Notification struct {
	ID           	uuid.UUID 		`gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	SourceUserID 	uuid.UUID 		`gorm:"type:uuid;not null"`
	SourceUser   	users.User		`gorm:"foreignKey:SourceUserID;constraint:OnDelete:CASCADE"`
	RecipientID  	uuid.UUID 		`gorm:"type:uuid;not null"`
	RecipientUser 	users.User		`gorm:"foreignKey:RecipientID;constraint:OnDelete:CASCADE"`
	PostID 			uuid.UUID     	`gorm:"type:uuid"`
	Post   			posts.Post    	`gorm:"foreignKey:PostID;references:ID;constraint:OnDelete:CASCADE"`
	Type         	string    		`gorm:"type:varchar(255)"`
	Content      	string    		`gorm:"type:varchar(255)"`
	CreatedAt    	time.Time 		`gorm:"type:timestamp"`
	UpdatedAt    	time.Time 		`gorm:"type:timestamp"`
}