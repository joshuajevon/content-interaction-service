package responses

import (
	"time"

	"github.com/google/uuid"
)

type PostResponse struct {
	ID		  uuid.UUID   `json:"id"`
	UserID    uuid.UUID   `json:"user_id"`
	ImageURLs []string    `json:"image_urls"`
	Caption   string      `json:"caption"`
	Tags      []string    `json:"tags"`
	CreatedAt time.Time	  `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}
