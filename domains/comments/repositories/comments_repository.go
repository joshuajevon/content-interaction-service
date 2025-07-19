package repositories

import (
	comments "bootcamp-content-interaction-service/domains/comments"
	"bootcamp-content-interaction-service/domains/comments/entities"
	"bootcamp-content-interaction-service/infrastructures"
	"context"
	"errors"
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
		return err
	}

	pId, err := uuid.Parse(postId)
	if err != nil {
		return err
	}

	var rId *uuid.UUID
	if replyId != nil{
		parsed, err := uuid.Parse(*replyId)
		if err != nil {
			return err
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
		return err
	}

	return nil
}

func (repo *CommentsRepository) UpdateComment(ctx context.Context, userId, postId, msg string, replyId *string) error {
	var comment entities.Comments

	err := repo.db.GetInstance().WithContext(ctx).
		Unscoped().
		Where("user_id=? AND post_id=? AND reply_id=?", userId, postId, replyId).
		First(&comment).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	err = repo.db.GetInstance().WithContext(ctx).Model(&comment).
		Updates(map[string]interface{}{
			"UpdatedAt": time.Now(),
			"Msg":       msg,
		}).Error

	if err != nil {
		return err
	}

	return nil
}

func (repo *CommentsRepository) ReplyComment(ctx context.Context, userId, postId, replyId, msg string) error {
	var comment entities.Comments

	err := repo.db.GetInstance().WithContext(ctx).
		Unscoped().
		Where("user_id=? AND post_id=? AND reply_id=?", userId, postId, replyId).
		First(&comment).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	err = repo.CreateComment(ctx, userId, postId, msg, &replyId)
	if err != nil{
		return err
	}

	return nil
}

func (repo *CommentsRepository) FindAllComment(ctx context.Context, postId string) (*[]entities.Comments, error) {
	var comment *[]entities.Comments

	err := repo.db.GetInstance().WithContext(ctx).
		Unscoped().
		Where("post_id=?", postId).
		First(&comment).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
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
			return nil, err
		}

		allReplies = append(allReplies, children...)
		queue = append(queue, children...) 
	}

	return allReplies, nil
}

func (repo *CommentsRepository) DeleteComment(ctx context.Context, id uuid.UUID) (error) {
	allReplies, err := repo.getAllReplyIDs(ctx, id)
	if err != nil{
		return nil
	}

	for len(allReplies) > 0{
		err = repo.db.GetInstance().WithContext(ctx).
			Where("id IN ?", allReplies).
			Delete(&entities.Comments{}).Error
		if err != nil{
			return err
		}
	}

	err = repo.db.GetInstance().WithContext(ctx).
		Where("id=?", id).Delete(&entities.Comments{}).Error
	if err != nil{
		return err
	}
	
	return nil
}