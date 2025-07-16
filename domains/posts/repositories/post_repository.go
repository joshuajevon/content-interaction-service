package repositories

import (
	"bootcamp-content-interaction-service/domains/posts"
	"bootcamp-content-interaction-service/domains/posts/entities"
	"bootcamp-content-interaction-service/infrastructures"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type PostRepository struct {
	db infrastructures.Database
}

func NewPostRepository(db infrastructures.Database) posts.PostRepository {
	return PostRepository{
		db: db,
	}
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

    return nil
}

func (p PostRepository) FindById(ctx context.Context, id string) (*entities.Post, error) {
	var post *entities.Post

	result := p.db.GetInstance().WithContext(ctx).Where("id = ?", id).First(&post)

	if result.Error != nil {
		return nil, result.Error
	}

	return post, nil
}

func (p PostRepository) FindAllByUserId(ctx context.Context, userId string) ([]*entities.Post, error) {
	var posts []*entities.Post

	result := p.db.GetInstance().WithContext(ctx).Where("user_id = ?", userId).Order("created_at DESC").Find(&posts)
	if result.Error != nil {
		return nil, result.Error
	}
	return posts, nil
}

func (p PostRepository) FindAll(ctx context.Context) ([]*entities.Post, error) {
	var posts []*entities.Post

	result := p.db.GetInstance().WithContext(ctx).Order("created_at DESC").Find(&posts)

	if result.Error != nil {
		return nil, result.Error
	}
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

	return postModel, nil
}
