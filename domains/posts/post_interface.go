package posts

import (
	"bootcamp-content-interaction-service/domains/posts/entities"
	"bootcamp-content-interaction-service/domains/posts/models/requests"
	"bootcamp-content-interaction-service/domains/posts/models/responses"
	"context"
)

type PostUseCase interface {
	CreatePost(ctx context.Context, request *requests.CreatePostRequest) (*responses.PostResponse, error)
	ViewAllPost(ctx context.Context) ([]*responses.PostResponse, error)
}

type PostRepository interface {
	SavePost(ctx context.Context, post *entities.Post) (*entities.Post, error)
	FindAll(ctx context.Context) ([]*entities.Post, error)
}