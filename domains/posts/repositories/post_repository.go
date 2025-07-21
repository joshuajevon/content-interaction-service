package repositories

import (
	"bootcamp-content-interaction-service/domains/posts"
	"bootcamp-content-interaction-service/domains/posts/entities"
	"bootcamp-content-interaction-service/infrastructures"
	"context"
	"encoding/json"
	"fmt"
	"time"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/lib/pq"
	"go.uber.org/zap"
)

type PostRepository struct {
	db infrastructures.Database
	redisCache *redis.Client
}

func NewPostRepository(db infrastructures.Database, redisClient *redis.Client) posts.PostRepository {
	return PostRepository{
		db: db,
		redisCache: redisClient,
	}
}

func (p PostRepository) UpdatePost(ctx context.Context, post *entities.Post) (*entities.Post, error) {
	var logger = zap.NewExample()
    result := p.db.GetInstance().WithContext(ctx).Save(post)
    if result.Error != nil {
        return nil, result.Error
    }

    // update redis
    postJSON, err := json.Marshal(post)
    if err == nil {
        key := "post:" + post.ID.String()
        _ = p.redisCache.Set(ctx, key, postJSON, time.Hour).Err()
		logger.Info("update post:"+post.ID.String()+" in redis")
    }

    return post, nil
}

func (p PostRepository) DeletePost(ctx context.Context, id string) error {
	var logger = zap.NewExample()
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
	logger.Info("delete post:"+parsedID.String()+" from redis")

    return nil
}
func (p PostRepository) FindById(ctx context.Context, id string) (*entities.Post, error) {
    var logger = zap.NewExample()
	var post entities.Post
    key := "post:" + id

    // get from redis
    cached, err := p.redisCache.Get(ctx, key).Result()
    if err == nil {
		logger.Info("Redis HIT: " + key)
        if err := json.Unmarshal([]byte(cached), &post); err == nil {
            return &post, nil
        }
    }

    // else, get from db
    result := p.db.GetInstance().WithContext(ctx).Where("id = ?", id).First(&post)
	logger.Info("Get from DB with id" + id)
    if result.Error != nil {
        return nil, result.Error
    }

    // set in redis
    postJSON, _ := json.Marshal(post)
    _ = p.redisCache.Set(ctx, key, postJSON, time.Hour).Err()
	logger.Info("Set post in redis cache " + string(postJSON))

    return &post, nil
}

func (p PostRepository) FindAllByUserId(ctx context.Context, userId string) ([]*entities.Post, error) {
    key := "user_posts:" + userId
    var posts []*entities.Post
	var logger = zap.NewExample()

    // get from redis
    cached, err := p.redisCache.LRange(ctx, key, 0, -1).Result()
    if err == nil && len(cached) > 0 {
		logger.Info("Redis HIT: " + key)
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
    logger.Info("Get from DB with userid" + userId)
	if result.Error != nil {
        return nil, result.Error
    }

    // set in redis
    for _, post := range posts {
        postJSON, _ := json.Marshal(post)
        _ = p.redisCache.LPush(ctx, key, postJSON).Err()
		logger.Info("Set post in redis cache " + string(postJSON))
    }
    _ = p.redisCache.Expire(ctx, key, time.Hour)

    return posts, nil
}

func (p PostRepository) FindAll(ctx context.Context) ([]*entities.Post, error) {
    key := "feed_posts"
    var posts []*entities.Post
	var logger = zap.NewExample()

	// get from redis
    cached, err := p.redisCache.LRange(ctx, key, 0, -1).Result()
    if err == nil && len(cached) > 0 {
		logger.Info("Redis HIT: " + key)
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
    logger.Info("Get all data from DB")
	if result.Error != nil {
        return nil, result.Error
    }

    // set in redis
    for _, post := range posts {
        postJSON, _ := json.Marshal(post)
        _ = p.redisCache.LPush(ctx, key, postJSON).Err()
		logger.Info("Set post in redis cache " + string(postJSON))
    }
    _ = p.redisCache.LTrim(ctx, key, 0, 99)          
    _ = p.redisCache.Expire(ctx, key, time.Hour)

    return posts, nil
}

func (p PostRepository) SavePost(ctx context.Context, post *entities.Post) (*entities.Post, error) {
    var logger = zap.NewExample()

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
	logger.Info("set redis cache with key post:" + postModel.ID.String())

    // FindAllByUserId post
    userFeedKey := "user_posts:" + postModel.UserID.String()
    _ = p.redisCache.LPush(ctx, userFeedKey, postJSON)
    _ = p.redisCache.Expire(ctx, userFeedKey, time.Hour)
	logger.Info("set redis cache with key user_posts:"+postModel.ID.String())

    // FindAll post
    feedKey := "feed_posts"
    _ = p.redisCache.LPush(ctx, feedKey, postJSON)
    _ = p.redisCache.LTrim(ctx, feedKey, 0, 99)
    _ = p.redisCache.Expire(ctx, feedKey, time.Hour)
	logger.Info("set redis cache with key feed_posts")

    return postModel, nil
}

func (p PostRepository) FindByUserIDs(ctx context.Context, userIds []string) ([]*entities.Post, error) {
	var posts []*entities.Post

	result := p.db.GetInstance().
		WithContext(ctx).
		Where("user_id IN ?", userIds).
		Order("created_at DESC").
		Find(&posts)

	if result.Error != nil {
		return nil, result.Error
	}
	
	return posts, nil
}