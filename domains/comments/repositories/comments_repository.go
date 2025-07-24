package repositories

import (
	comments "bootcamp-content-interaction-service/domains/comments"
	"bootcamp-content-interaction-service/domains/comments/entities"
	"bootcamp-content-interaction-service/infrastructures"
	"bootcamp-content-interaction-service/shared/util"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type CommentsRepository struct {
	db infrastructures.Database
	redisCache *redis.Client
    logger util.Logger
}

func NewCommentsRepository(db infrastructures.Database, redisClient *redis.Client, logger util.Logger) comments.CommentsRepository {
	return &CommentsRepository{db: db, redisCache: redisClient, logger: logger}
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

	cacheKey := "comments:post:" + postId
	delErr := repo.redisCache.Del(ctx, cacheKey).Err()
	if delErr != nil {
		repo.logger.Warn("failed to delete redis cache",
			zap.String("cacheKey", cacheKey),
			zap.Error(delErr),
		)
	} else {
		repo.logger.Info("deleted redis cache for comments",
			zap.String("postId", postId),
		)
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

	cacheKey := "comments:post:" + comment.PostId.String()
	delErr := repo.redisCache.Del(ctx, cacheKey).Err()
	if delErr != nil {
		repo.logger.Warn("failed to delete redis cache",
			zap.String("cacheKey", cacheKey),
			zap.Error(delErr),
		)
	} else {
		repo.logger.Info("deleted redis cache for comments",
			zap.String("postId", comment.PostId.String()),
		)
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

	cacheKey := "comments:post:" + postId
	delErr := repo.redisCache.Del(ctx, cacheKey).Err()
	if delErr != nil {
		repo.logger.Warn("failed to delete redis cache",
			zap.String("cacheKey", cacheKey),
			zap.Error(delErr),
		)
	} else {
		repo.logger.Info("deleted redis cache for comments",
			zap.String("postId", postId),
		)
	}

	return nil
}

func (repo *CommentsRepository) FindAllComment(ctx context.Context, postId string) (*[]entities.Comments, error) {
	var comment *[]entities.Comments
	cacheKey := "comments:posts:" + postId

	val, err := repo.redisCache.Get(ctx, cacheKey).Result()
	if err == nil {
		if unmarshalErr := json.Unmarshal([]byte(val), &comment); unmarshalErr == nil {
			repo.logger.Info("Cache hit for comments",
				zap.String("postId", postId),
			)
			return comment, nil
		} else {
			repo.logger.Warn("Failed to unmarshal cached comments",
				zap.Error(unmarshalErr),
			)
		}
	} else if err != redis.Nil {
		repo.logger.Error("Redis error on Get",
			zap.Error(err),
		)
	} else {
		repo.logger.Info("Cache miss for comments",
			zap.String("postId", postId),
		)
	}

	err = repo.db.GetInstance().WithContext(ctx).
		Unscoped().
		Where("post_id=?", postId).
		Find(&comment).Error

	if err != nil {
		repo.logger.Error("Database error while getting comments",
			zap.String("postId", postId),
			zap.Error(err),
		)
		return nil, errors.New("post not found")
	}

	if len(*comment) == 0 {
		repo.logger.Info("No comments found",
			zap.String("postId", postId),
		)
		return nil, errors.New("there's no comment")
	}

	data, marshalErr := json.Marshal(comment)
	if marshalErr == nil {
		cacheErr := repo.redisCache.Set(ctx, cacheKey, data, time.Minute*5).Err()
		if cacheErr != nil {
			repo.logger.Warn("Failed to set Redis cache for comments",
				zap.String("postId", postId),
				zap.Error(cacheErr),
			)
		} else {
			repo.logger.Info("Cached comments to Redis",
				zap.String("postId", postId),
			)
		}
	} else {
		repo.logger.Warn("Failed to marshal comments to JSON",
			zap.Error(marshalErr),
		)
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
	var parentComment entities.Comments
	err := repo.db.GetInstance().WithContext(ctx).
		Where("id = ?", id).
		First(&parentComment).Error
	if err != nil {
		return errors.New("failed to find parent comment")
	}

	postId := parentComment.PostId.String()
	
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

	cacheKey := "comments:post:" + postId
	delErr := repo.redisCache.Del(ctx, cacheKey).Err()
	if delErr != nil {
		repo.logger.Warn("failed to delete redis cache",
			zap.String("cacheKey", cacheKey),
			zap.Error(delErr),
		)
	} else {
		repo.logger.Info("deleted redis cache for comments",
			zap.String("postId", postId),
		)
	}

	return nil
}
