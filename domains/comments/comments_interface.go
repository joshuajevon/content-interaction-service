package comments

import (
	"bootcamp-content-interaction-service/domains/comments/entities"
	"context"

	"github.com/google/uuid"
)

type CommentsUseCase interface {
	CreateComment(ctx context.Context, userId, postId, msg string, replyId *string) error
	UpdateComment(ctx context.Context, id, userId, msg string) error
	ReplyComment(ctx context.Context, id, userId, postId, msg string) error
	FindAllComment(ctx context.Context, postId string) (*[]entities.Comments,error)
	DeleteComment(ctx context.Context, id uuid.UUID) (error) 
}

type CommentsRepository interface {
	CreateComment(ctx context.Context, userId, postId, msg string, replyId *string) error
	UpdateComment(ctx context.Context, id, userId, msg string) error
	ReplyComment(ctx context.Context, id, userId, postId, msg string) error
	FindAllComment(ctx context.Context, postId string) (*[]entities.Comments, error)
	DeleteComment(ctx context.Context, id uuid.UUID) (error) 
}
