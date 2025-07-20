package repositories

import (
	comments "bootcamp-content-interaction-service/domains/comments"
	"bootcamp-content-interaction-service/domains/comments/entities"
	"bootcamp-content-interaction-service/infrastructures"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CommentsRepository struct {
	db infrastructures.Database
}

func NewCommentsRepository(db infrastructures.Database) comments.CommentsRepository {
	return &CommentsRepository{db: db}
}

func (repo *CommentsRepository) CreateComment(ctx context.Context, userId, postId, msg string, replyId *string) error {
	var comment entities.Comments

	uId, err := uuid.Parse(userId)
	if err != nil {
		return errors.New("failed parsing userId")
	}

	pId, err := uuid.Parse(postId)
	if err != nil {
		return errors.New("failed parsing postId")
	}

	var rId *uuid.UUID
	if replyId != nil {
		parsed, err := uuid.Parse(*replyId)
		if err != nil {
			return errors.New("failed parsing replyId")
		}
		rId = &parsed
	}

	comment = entities.Comments{
		ID:        uuid.New(),
		UserID:    uId,
		PostId:    pId,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		ReplyId:   rId,
		Msg:       msg,
	}

	err = repo.db.GetInstance().WithContext(ctx).Create(&comment).Error
	if err != nil {
		return errors.New("failed to create comment")
	}

	return nil
}

func (repo *CommentsRepository) UpdateComment(ctx context.Context, id, userId, msg string) error {
	var comment entities.Comments

	err := repo.db.GetInstance().WithContext(ctx).
		Unscoped().
		Where("id=?", id).
		First(&comment).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("record not found")
	} else if err != nil {
		return err
	}

	if fmt.Sprintf("%v", comment.UserID) != userId {
		return errors.New("authorization : you cannot update other comment")
	}

	err = repo.db.GetInstance().WithContext(ctx).Model(&comment).Unscoped().
		Updates(map[string]interface{}{
			"updated_at": time.Now(),
			"msg":        msg,
		}).Error

	if err != nil {
		return errors.New("failed to update comment")
	}

	return nil
}

func (repo *CommentsRepository) ReplyComment(ctx context.Context, id, userId, postId, msg string) error {
	var comment entities.Comments

	err := repo.db.GetInstance().WithContext(ctx).
		Unscoped().
		Where("id=?", id).
		First(&comment).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("the comment_id that you reply doesn't exist")
	}

	err = repo.CreateComment(ctx, userId, postId, msg, &id)
	if err != nil {
		return err
	}

	return nil
}

func (repo *CommentsRepository) FindAllComment(ctx context.Context, postId string) (*[]entities.Comments, error) {
	var comment *[]entities.Comments

	err := repo.db.GetInstance().WithContext(ctx).
		Unscoped().
		Where("post_id=?", postId).
		Find(&comment).Error

	if len(*comment) == 0 {
		return nil, errors.New("there's no comment")
	} else if err != nil {
		return nil, errors.New("post not found")
	}

	return comment, nil
}

func (repo *CommentsRepository) getAllReplyIDs(ctx context.Context, parentID uuid.UUID) ([]uuid.UUID, error) {
	var allReplies []uuid.UUID

	var queue []uuid.UUID
	queue = append(queue, parentID)

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		var children []uuid.UUID
		err := repo.db.GetInstance().WithContext(ctx).
			Model(&entities.Comments{}).
			Where("reply_id = ?", current).
			Pluck("id", &children).Error

		if err != nil {
			return nil, errors.New("failed to get reply data")
		}

		allReplies = append(allReplies, children...)
		queue = append(queue, children...)
	}

	return allReplies, nil
}

func (repo *CommentsRepository) DeleteComment(ctx context.Context, id uuid.UUID) error {
	allReplies, err := repo.getAllReplyIDs(ctx, id)
	if err != nil {
		return err
	}

	if len(allReplies) > 0{
		err = repo.db.GetInstance().WithContext(ctx).
			Where("id IN ?", allReplies).
			Delete(&entities.Comments{}).Error
		if err != nil {
			return errors.New("failed to delete comment data : reply_id")
		}
	}

	err = repo.db.GetInstance().WithContext(ctx).
		Where("id=?", id).Delete(&entities.Comments{}).Error
	if err != nil {
		return errors.New("failed to delete comment data : mother_id")
	}

	return nil
}
