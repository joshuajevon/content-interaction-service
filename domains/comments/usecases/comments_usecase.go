package usecases

import (
	comments "bootcamp-content-interaction-service/domains/comments"
	"bootcamp-content-interaction-service/domains/comments/entities"
	"context"

	"github.com/google/uuid"
)

type CommentsUseCase struct {
	repo comments.CommentsRepository
}

func NewCommentsUseCase(repo comments.CommentsRepository) comments.CommentsUseCase {
	return &CommentsUseCase{repo: repo}
}

func (uc *CommentsUseCase) CreateComment(ctx context.Context, userId, postId, msg string, replyId *string) error {
	err := uc.repo.CreateComment(ctx, userId, postId, msg, nil)
	if err != nil {
		return err
	}

	return nil
}

func (uc *CommentsUseCase) UpdateComment(ctx context.Context, userId, postId, msg string, replyId *string) error {
	err := uc.repo.UpdateComment(ctx, userId, postId, msg, nil)
	if err != nil {
		return err
	}

	return nil
}

func (uc *CommentsUseCase) ReplyComment(ctx context.Context, userId, postId, replyId, msg string) error {
	err := uc.repo.ReplyComment(ctx, userId, postId, replyId, msg)
	if err != nil {
		return err
	}

	return nil
}

func (uc *CommentsUseCase) FindAllComment(ctx context.Context, postId string) (*[]entities.Comments,error) {
	comment, err := uc.repo.FindAllComment(ctx, postId)
	if err != nil {
		return nil, err
	}

	return comment, nil
}

func (uc *CommentsUseCase) DeleteComment(ctx context.Context, id uuid.UUID) (error)  {
	err := uc.repo.DeleteComment(ctx, id)
	if err != nil {
		return err
	}

	return err
}