package repositories

import (
	"bootcamp-content-interaction-service/domains/posts"
	"bootcamp-content-interaction-service/domains/posts/entities"
	"bootcamp-content-interaction-service/infrastructures"
	"bootcamp-content-interaction-service/shared/util"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type PostRepository struct {
	db infrastructures.Database
	redisCache *redis.Client
    logger util.Logger
}

func NewPostRepository(db infrastructures.Database, redisClient *redis.Client, logger util.Logger) posts.PostRepository {
	return PostRepository{
		db: db,
		redisCache: redisClient,
        logger: logger,
	}
}

func (p PostRepository) UpdatePost(ctx context.Context, post *entities.Post) (*entities.Post, error) {
    result := p.db.GetInstance().WithContext(ctx).Save(post)
    if result.Error != nil {
        return nil, result.Error
    }

    // update redis
    postJSON, err := json.Marshal(post)
    if err == nil {
        key := "post:" + post.ID.String()
        _ = p.redisCache.Set(ctx, key, postJSON, time.Hour).Err()
		p.logger.Info("Update post in redis",
            zap.String("post_id", post.ID.String()),
        )
    }

    return post, nil
}

func (p PostRepository) DeletePost(ctx context.Context, id string) error {
    parsedID, err := uuid.Parse(id)
    if err != nil {
        return fmt.Errorf("invalid UUID format: %w", err)
    }

    result := p.db.GetInstance().WithContext(ctx).Where("id = ?", parsedID).Delete(&entities.Post{})
    if result.Error != nil {
        return result.Error
    }
    if result.RowsAffected == 0 {
        return fmt.Errorf("no post found with ID: %s", id)
    }

    // delete from redis
    key := "post:" + parsedID.String()
    _ = p.redisCache.Del(ctx, key).Err()
	p.logger.Info("Delete post from redis",
        zap.String("post_id", parsedID.String()),
    )

    return nil
}
func (p PostRepository) FindById(ctx context.Context, id string) (*entities.Post, error) {
	var post entities.Post
    key := "post:" + id

    // get from redis
    cached, err := p.redisCache.Get(ctx, key).Result()
    if err == nil {
		p.logger.Info("Cache hit - returning posts from redis", 
            zap.String("cache_key", key),
        )
        if err := json.Unmarshal([]byte(cached), &post); err == nil {
            return &post, nil
        }
    }

    // else, get from db
    result := p.db.GetInstance().WithContext(ctx).Where("id = ?", id).First(&post)
	p.logger.Info("Get from DB",
        zap.String("id", id),
    )
    if result.Error != nil {
        return nil, result.Error
    }

    // set in redis
    postJSON, _ := json.Marshal(post)
    _ = p.redisCache.Set(ctx, key, postJSON, time.Hour).Err()
	p.logger.Info("Set post in cache",
        zap.String("post_json", string(postJSON)),
    )

    return &post, nil
}

func (p PostRepository) FindAllByUserId(ctx context.Context, userId string) ([]*entities.Post, error) {
    key := "user_posts:" + userId
    var posts []*entities.Post

    // get from redis
    cached, err := p.redisCache.LRange(ctx, key, 0, -1).Result()
    if err == nil && len(cached) > 0 {
		p.logger.Info("Cache hit - returning posts from redis", 
            zap.String("cache_key", key),
        )
        for _, jsonItem := range cached {
            var post entities.Post
            if err := json.Unmarshal([]byte(jsonItem), &post); err == nil {
                posts = append(posts, &post)
            }
        }
        return posts, nil
    }

    // else, get from db
    result := p.db.GetInstance().WithContext(ctx).Where("user_id = ?", userId).Order("created_at DESC").Find(&posts)
    p.logger.Info("Get from DB",
        zap.String("user_id", userId),
    )
	if result.Error != nil {
        return nil, result.Error
    }

    // set in redis
    for _, post := range posts {
        postJSON, _ := json.Marshal(post)
        _ = p.redisCache.LPush(ctx, key, postJSON).Err()
		p.logger.Info("Set post in cache",
            zap.String("post_json", string(postJSON)),
        )
    }
    _ = p.redisCache.Expire(ctx, key, time.Hour)

    return posts, nil
}

func (p PostRepository) FindAll(ctx context.Context) ([]*entities.Post, error) {
    key := "feed_posts"
    var posts []*entities.Post

	// get from redis
    cached, err := p.redisCache.LRange(ctx, key, 0, -1).Result()
    if err == nil && len(cached) > 0 {
		p.logger.Info("Cache hit - returning posts from redis", 
            zap.String("cache_key", key),
        )
        for _, jsonItem := range cached {
            var post entities.Post
            if err := json.Unmarshal([]byte(jsonItem), &post); err == nil {
                posts = append(posts, &post)
            }
        }
        return posts, nil
    }

	// else, get from db
    result := p.db.GetInstance().WithContext(ctx).Order("created_at DESC").Find(&posts)
    p.logger.Info("Get all data from DB")
	if result.Error != nil {
        return nil, result.Error
    }

    // set in redis
    for _, post := range posts {
        postJSON, _ := json.Marshal(post)
        _ = p.redisCache.LPush(ctx, key, postJSON).Err()
		p.logger.Info("Set post in cache",
            zap.String("post_json", string(postJSON)),
        )
    }
    _ = p.redisCache.LTrim(ctx, key, 0, 99)          
    _ = p.redisCache.Expire(ctx, key, time.Hour)

    return posts, nil
}

func (p PostRepository) SavePost(ctx context.Context, post *entities.Post) (*entities.Post, error) {
    

	postModel := &entities.Post{
        ID:        uuid.New(),
        UserID:    post.UserID,
        ImageURLs: pq.StringArray(post.ImageURLs),
        Caption:   post.Caption,
        Tags:      pq.StringArray(post.Tags),
        CreatedAt: time.Now(),
    }

    result := p.db.GetInstance().WithContext(ctx).Create(postModel)
    if result.Error != nil {
        return nil, result.Error
    }

    postJSON, err := json.Marshal(postModel)
    if err != nil {
        fmt.Printf("Redis marshal failed: %v\n", err)
        return postModel, nil
    }

    // Find by id post
    _ = p.redisCache.Set(ctx, "post:"+postModel.ID.String(), postJSON, time.Hour).Err()
    p.logger.Info("Set in cache",
        zap.String("post_key", postModel.ID.String()),
    )

    // FindAllByUserId post
    userFeedKey := "user_posts:" + postModel.UserID.String()
    _ = p.redisCache.LPush(ctx, userFeedKey, postJSON)
    _ = p.redisCache.Expire(ctx, userFeedKey, time.Hour)
	p.logger.Info("Set in cache",
        zap.String("user_posts_key", postModel.ID.String()),
    )

    // FindAll post
    feedKey := "feed_posts"
    _ = p.redisCache.LPush(ctx, feedKey, postJSON)
    _ = p.redisCache.LTrim(ctx, feedKey, 0, 99)
    _ = p.redisCache.Expire(ctx, feedKey, time.Hour)
	p.logger.Info("Set in cache cache with key feed_posts")

    return postModel, nil
}

func (p PostRepository) FindByUserIDs(ctx context.Context, userIds []string, limit, offset int) ([]*entities.Post, error) {
	var posts []*entities.Post
	cacheKey := fmt.Sprintf("posts:user_ids:%s:limit:%d:offset:%d", strings.Join(userIds, ","), limit, offset)

	p.logger.Info("Looking up posts by user IDs with pagination",
		zap.Strings("user_ids", userIds),
		zap.String("cache_key", cacheKey),
		zap.Int("limit", limit),
		zap.Int("offset", offset),
	)

	cached, err := p.redisCache.Get(ctx, cacheKey).Result()
	if err == nil {
		if err := json.Unmarshal([]byte(cached), &posts); err == nil {
			p.logger.Info("Cache hit - returning posts from Redis",
				zap.String("cache_key", cacheKey),
				zap.Int("post_count", len(posts)),
			)
			return posts, nil
		} else {
			p.logger.Error("Failed to unmarshal cached posts",
				zap.Error(err),
				zap.String("cache_key", cacheKey),
			)
		}
	} else if err != redis.Nil {
		p.logger.Error("Redis GET operation failed",
			zap.Error(err),
			zap.String("cache_key", cacheKey),
		)
	} else {
		p.logger.Info("Cache miss - querying database",
			zap.String("cache_key", cacheKey),
		)
	}

	result := p.db.GetInstance().
		WithContext(ctx).
		Where("user_id IN ?", userIds).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&posts)

	if result.Error != nil {
		p.logger.Error("Database query failed",
			zap.Error(result.Error),
			zap.Strings("user_ids", userIds),
		)
		return nil, result.Error
	}

	bytes, err := json.Marshal(posts)
	if err != nil {
		p.logger.Error("Failed to marshal posts for caching",
			zap.Error(err),
			zap.String("cache_key", cacheKey),
		)
	} else {
		err := p.redisCache.Set(ctx, cacheKey, bytes, 30*time.Second).Err()
		if err != nil {
			p.logger.Error("Failed to set cache",
				zap.Error(err),
				zap.String("cache_key", cacheKey),
			)
		} else {
			p.logger.Info("Cached posts",
				zap.String("cache_key", cacheKey),
				zap.Int("post_count", len(posts)),
			)
		}
	}

	return posts, nil
}

