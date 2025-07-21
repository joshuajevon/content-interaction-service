package repositories

import (
	likes "bootcamp-content-interaction-service/domains/likes"
	"bootcamp-content-interaction-service/domains/likes/entities"
	"bootcamp-content-interaction-service/infrastructures"
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LikesRepository struct {
	db infrastructures.Database
}

func NewLikesRepository(db infrastructures.Database) likes.LikesRepository {
	return &LikesRepository{db: db}
}

func (repo *LikesRepository) LikePost(ctx context.Context, userId, postId string) error {
	var likes entities.Likes

	uId, err := uuid.Parse(userId)
	if err != nil {
		return errors.New("failed to parse userId")
	}

	pId, err := uuid.Parse(postId)
	if err != nil {
		return errors.New("failed to parse postId")
	}

	err = repo.db.GetInstance().WithContext(ctx).
		Unscoped().
		Where("user_id=? AND post_id=?", userId, postId).
		First(&likes).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		likes = entities.Likes{
			UserID:    uId,
			PostId:    pId,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		err = repo.db.GetInstance().WithContext(ctx).Create(&likes).Error
		if err != nil {
			return errors.New("failed adding to like database")
		}
	}

	if likes.DeletedAt.Valid {
		err = repo.db.GetInstance().WithContext(ctx).Unscoped().Model(&likes).
			Updates(map[string]interface{}{
				"deleted_at": nil,
				"updated_at": time.Now(),
			}).Error

		if err != nil {
			return errors.New("failed to update like data")
		}
	}

	return nil
}

func (repo *LikesRepository) DislikePost(ctx context.Context, userId, postId string) error {
	var likes entities.Likes

	err := repo.db.GetInstance().WithContext(ctx).
		Unscoped().
		Where("user_id=? AND post_id=?", userId, postId).
		First(&likes).Error

	if err != nil {
		return errors.New("failed to search like data")
	}

	if !likes.DeletedAt.Valid {
		err = repo.db.GetInstance().WithContext(ctx).Model(&likes).
			Updates(map[string]interface{}{
				"deleted_at": time.Now(),
				"updated_at": time.Now(),
			}).Error

		if err != nil {
			return errors.New("failed to update like database")
		}
	}else{
		return errors.New("you have dislike this post")
	}

	return nil
}
