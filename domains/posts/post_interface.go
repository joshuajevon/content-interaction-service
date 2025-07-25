package posts

import (
	"bootcamp-content-interaction-service/domains/posts/entities"
	"bootcamp-content-interaction-service/domains/posts/models/requests"
	"bootcamp-content-interaction-service/domains/posts/models/responses"
	sharedResponse "bootcamp-content-interaction-service/shared/models/responses"
	"context"
)

type PostUseCase interface {
	CreatePost(ctx context.Context, request *requests.CreatePostRequest) (*responses.PostResponse, error)
	ViewAllPost(ctx context.Context) ([]*responses.PostResponse, error)
	ViewAllPostByUserId(ctx context.Context) ([]*responses.PostResponse, error)
	ViewPostById(ctx context.Context, id string) (*responses.PostResponse, error)
	DeletePost(ctx context.Context, id string) (*sharedResponse.BasicResponse, error)
	UpdatePost(ctx context.Context, postId string, request *requests.UpdatePostRequest) (*responses.PostResponse, error)
	ViewPostByUserId(ctx context.Context, userId string, limit int, offset int) ([]*responses.PostResponse, error)
}

type PostRepository interface {
	SavePost(ctx context.Context, post *entities.Post) (*entities.Post, error)
	FindAll(ctx context.Context) ([]*entities.Post, error)
	FindAllByUserId(ctx context.Context, userId string) ([]*entities.Post, error)
	FindById(ctx context.Context, id string) (*entities.Post, error)
	DeletePost(ctx context.Context, id string) (error)
	UpdatePost(ctx context.Context, post *entities.Post) (*entities.Post, error)
	FindByUserIDs(ctx context.Context, userIds []string, limit int, offset int) ([]*entities.Post, error)
}