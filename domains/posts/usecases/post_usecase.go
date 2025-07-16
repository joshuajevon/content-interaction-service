package usecases

import (
	"bootcamp-content-interaction-service/domains/posts"
	"bootcamp-content-interaction-service/domains/posts/entities"
	"bootcamp-content-interaction-service/domains/posts/models/requests"
	"bootcamp-content-interaction-service/domains/posts/models/responses"
	"bootcamp-content-interaction-service/shared/util"
	sharedResponse "bootcamp-content-interaction-service/shared/models/responses"
	"fmt"

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

func (p PostUseCase) DeletePost(ctx context.Context, id string) (*sharedResponse.BasicResponse, error) {
	user, err := util.GetAuthUser(ctx)

	if err != nil {
		return nil, err
	}

 	post, err := p.postRepository.FindById(ctx, id);
	if err != nil {
		return nil, err
	}

	if post.UserID != uuid.MustParse(user.UserId) {
        return nil, fmt.Errorf("unauthorized: cannot delete someone else's post")
    }

   	p.postRepository.DeletePost(ctx, id);

	return &sharedResponse.BasicResponse{
		Data: struct {
			Message string
		}{
			Message: "Post with " + id + " deleted",
		},
	}, nil
}

func (p PostUseCase) ViewPostById(ctx context.Context, id string) (*responses.PostResponse, error) {
	post, err := p.postRepository.FindById(ctx, id)

	if err != nil {
		return nil, err
	}

	return &responses.PostResponse{
		ID:        post.ID,
		UserID:    post.UserID,
		ImageURLs: post.ImageURLs,
		Caption:   post.Caption,
		Tags:      post.Tags,
		CreatedAt: post.CreatedAt,
		UpdatedAt: post.UpdatedAt,
	}, nil
}

func (p PostUseCase) ViewAllPostByUserId(ctx context.Context) ([]*responses.PostResponse, error) {
	user, err := util.GetAuthUser(ctx)

	if err != nil {
		return nil, err
	}

	posts, err := p.postRepository.FindAllByUserId(ctx, user.UserId)

	if err != nil {
		return nil, err
	}

	var responseList []*responses.PostResponse
	for _, post := range posts {
		response := &responses.PostResponse{
			ID:        post.ID,
			UserID:    post.UserID,
			ImageURLs: post.ImageURLs,
			Caption:   post.Caption,
			Tags:      post.Tags,
			CreatedAt: post.CreatedAt,
			UpdatedAt: post.UpdatedAt,
		}
		responseList = append(responseList, response)
	}
	return responseList, nil
}

func (p PostUseCase) ViewAllPost(ctx context.Context) ([]*responses.PostResponse, error) {
	posts, err := p.postRepository.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	var responseList []*responses.PostResponse
	for _, post := range posts {
		response := &responses.PostResponse{
			ID:        post.ID,
			UserID:    post.UserID,
			ImageURLs: post.ImageURLs,
			Caption:   post.Caption,
			Tags:      post.Tags,
			CreatedAt: post.CreatedAt,
			UpdatedAt: post.UpdatedAt,
		}
		responseList = append(responseList, response)
	}
	return responseList, nil
}

func (p PostUseCase) CreatePost(ctx context.Context, request *requests.CreatePostRequest) (*responses.PostResponse, error) {
	user, err := util.GetAuthUser(ctx)

	if err != nil {
		return nil, err
	}
	
	postObject := &entities.Post{
		UserID:    uuid.MustParse(user.UserId),
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
		UpdatedAt: savedPost.UpdatedAt,
	}, nil
}
