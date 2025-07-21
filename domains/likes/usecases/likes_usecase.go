package usecases

import (
	likes "bootcamp-content-interaction-service/domains/likes"
	"context"
)

type LikesUseCase struct {
	repo likes.LikesRepository
}

func NewLikesUseCase(repo likes.LikesRepository) likes.LikesUseCase {
	return &LikesUseCase{repo: repo}
}

func (uc *LikesUseCase) LikePost(ctx context.Context, userId, postId string) error {
	err := uc.repo.LikePost(ctx, userId, postId)
	if err != nil {
		return err
	}

	return nil
}

func (uc *LikesUseCase) DislikePost(ctx context.Context, userId, postId string) error {
	err := uc.repo.DislikePost(ctx, userId, postId)
	if err != nil {
		return err
	}

	return nil
}
