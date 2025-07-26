package usecases

import (
	"bootcamp-content-interaction-service/domains/notifications"
	notification "bootcamp-content-interaction-service/domains/notifications/entities"
	"bootcamp-content-interaction-service/domains/posts"
	"bootcamp-content-interaction-service/domains/posts/entities"
	"bootcamp-content-interaction-service/domains/posts/handlers/http"
	"bootcamp-content-interaction-service/domains/posts/models/requests"
	"bootcamp-content-interaction-service/domains/posts/models/responses"
	sharedResponse "bootcamp-content-interaction-service/shared/models/responses"
	"bootcamp-content-interaction-service/shared/util"
	"fmt"
	"os"
	"time"

	"context"

	"github.com/google/uuid"
)

type PostUseCase struct {
	postRepository posts.PostRepository
	userGraphService http.UserGraphService
	notifRepository notifications.NotificationRepository
}

func NewPostUseCase(postRepo posts.PostRepository, userGraph http.UserGraphService, notifRepository notifications.NotificationRepository) posts.PostUseCase {
    return PostUseCase{
        postRepository: postRepo,
		userGraphService: userGraph,
		notifRepository: notifRepository,
	}
}

func (p PostUseCase) UpdatePost(ctx context.Context, postID string, request *requests.UpdatePostRequest) (*responses.PostResponse, error) {
    existing, err := p.postRepository.FindById(ctx, postID)
    if err != nil {
        return nil, err
    }

    user, err := util.GetAuthUser(ctx)
    if err != nil {
        return nil, err
    }

    if existing.UserID != uuid.MustParse(user.UserId) {
        return nil, fmt.Errorf("unauthorized: cannot update someone else's post")
    }

    if request.Caption != "" {
        existing.Caption = request.Caption
    }
    if request.Tags != nil {
        existing.Tags = request.Tags
    }
    if request.ImageURLs != nil {
        for _, old := range existing.ImageURLs {
            _ = os.Remove(old)
        }
        existing.ImageURLs = request.ImageURLs
    }
	
	existing.UpdatedAt = time.Now()

    updated, err := p.postRepository.UpdatePost(ctx, existing)
    if err != nil {
        return nil, err
    }

    return &responses.PostResponse{
        ID:        updated.ID,
        UserID:    updated.UserID,
        Caption:   updated.Caption,
        Tags:      updated.Tags,
        ImageURLs: updated.ImageURLs,
        UpdatedAt: updated.UpdatedAt,
    }, nil
}

func (p PostUseCase) DeletePost(ctx context.Context, id string) (*sharedResponse.BasicResponse, error) {
    user, err := util.GetAuthUser(ctx)
    if err != nil {
        return nil, err
    }

    post, err := p.postRepository.FindById(ctx, id)
    if err != nil {
        return nil, err
    }

    if post.UserID != uuid.MustParse(user.UserId) {
        return nil, fmt.Errorf("unauthorized: cannot delete someone else's post")
    }

    for _, imagePath := range post.ImageURLs {
        if err := os.Remove(imagePath); err != nil {
            fmt.Printf("Failed to remove image %s: %v\n", imagePath, err)
        }
    }

    if err := p.postRepository.DeletePost(ctx, id); err != nil {
        return nil, err
    }

    return &sharedResponse.BasicResponse{
        Data: struct {
            Message string
        }{
            Message: "Post with " + id + " deleted along with images",
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

	followers, err := p.userGraphService.GetFollowers(postObject.UserID.String())

	if err != nil {
		return &responses.PostResponse{
			ID:        savedPost.ID,
			UserID:    savedPost.UserID,
			ImageURLs: savedPost.ImageURLs,
			Caption:   savedPost.Caption,
			Tags:      savedPost.Tags,
			CreatedAt: savedPost.CreatedAt,
		}, nil
	}

	for _, follower := range followers {
		notif := &notification.Notification{
			SourceUserID: postObject.UserID,
			RecipientID:  uuid.MustParse(follower),
			PostID:       savedPost.ID,
			Type:         "NEW_POST",
			Content:      savedPost.Caption,
		}

		_, err := p.notifRepository.SaveNotification(ctx, notif)
		if err != nil {
			continue
		}
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

func (p PostUseCase) ViewPostByUserId(ctx context.Context, userId string, limit int, offset int) ([]*responses.PostResponse, error) {
	followingIDs, err := p.userGraphService.GetFollowings(userId)
	if err != nil {
		return nil, err
	}

	posts, err := p.postRepository.FindByUserIDs(ctx, followingIDs, limit, offset)
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

