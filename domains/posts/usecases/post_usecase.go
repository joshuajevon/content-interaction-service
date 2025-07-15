package usecases

import (
	"bootcamp-content-interaction-service/domains/posts"
	"bootcamp-content-interaction-service/domains/posts/entities"
	"bootcamp-content-interaction-service/domains/posts/models/requests"
	"bootcamp-content-interaction-service/domains/posts/models/responses"

	"context"

	"github.com/google/uuid"
)

type PostUseCase struct {
	postRepository posts.PostRepository
}

func NewPostUseCase(postRepo posts.PostRepository) posts.PostUseCase {
	return PostUseCase{
		postRepository: postRepo,
	}
}
func (p PostUseCase) CreatePost(ctx context.Context, request *requests.CreatePostRequest) (*responses.PostResponse, error) {
    postObject := &entities.Post{
        UserID:    uuid.MustParse(request.UserID),
        ImageURLs: request.ImageURLs,
        Caption:   request.Caption,
        Tags:      request.Tags,
    }

    savedPost, err := p.postRepository.SavePost(ctx, postObject)
    if err != nil {
        return nil, err
    }

    return &responses.PostResponse{
        ID:        savedPost.ID,
        UserID:    savedPost.UserID,
        ImageURLs: savedPost.ImageURLs,
        Caption:   savedPost.Caption,
        Tags:      savedPost.Tags,
        CreatedAt: savedPost.CreatedAt,
    }, nil
}