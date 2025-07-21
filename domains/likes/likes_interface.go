package likes

import (
	"context"
)

type LikesUseCase interface {
	LikePost(ctx context.Context, userId, postId string)error
	DislikePost(ctx context.Context, userId, postId string) error 
}

type LikesRepository interface {
	LikePost(ctx context.Context, userId, postId string) error
	DislikePost(ctx context.Context, userId, postId string) error
}