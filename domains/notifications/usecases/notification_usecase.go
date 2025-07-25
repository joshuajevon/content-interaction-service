package usecases

import (
	"bootcamp-content-interaction-service/domains/notifications"
	"bootcamp-content-interaction-service/domains/notifications/entities"
	"bootcamp-content-interaction-service/domains/notifications/models/requests"
	"bootcamp-content-interaction-service/domains/notifications/models/responses"
	"bootcamp-content-interaction-service/shared/util"
	"context"
	"time"

	"github.com/google/uuid"
)

type NotificationUseCase struct {
	notifRepo notifications.NotificationRepository
}

func NewNotificationUseCase(notifRepo notifications.NotificationRepository) notifications.NotificationUseCase {
	return NotificationUseCase{
		notifRepo: notifRepo,
	}
}

func (n NotificationUseCase) FindAllNotification(ctx context.Context) ([]*responses.PostNotificationResponse, error) {
	user, err := util.GetAuthUser(ctx)

	if err != nil {
		return nil, err
	}
	
	notifications, err := n.notifRepo.FindAll(ctx, string(user.UserId))
	if err != nil {
		return nil, err
	}

	var responseList []*responses.PostNotificationResponse
	for _, notification := range notifications {
		response := &responses.PostNotificationResponse{
			ID:           notification.ID.String(),
			SourceUserID: notification.SourceUserID.String(),
			RecipientID:  notification.RecipientID.String(),
			PostID:       notification.PostID.String(),
			Content:      notification.Content,
		}
		responseList = append(responseList, response)
	}
	return responseList, nil
}

func (n NotificationUseCase) NotifyNewPost(ctx context.Context, request *requests.PostNotificationRequest) (*responses.PostNotificationResponse, error) {
	notifObject := &entities.Notification{
		SourceUserID: uuid.MustParse(request.SourceUserID),
		RecipientID:  uuid.MustParse(request.RecipientID),
		PostID:       uuid.MustParse(request.PostID),
		Type:         util.NOTIF_POST,
		Content:      request.Content,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	savedPost, err := n.notifRepo.SaveNotification(ctx, notifObject)
	if err != nil {
		return nil, err
	}

	return &responses.PostNotificationResponse{
		ID:           savedPost.ID.String(),
		SourceUserID: savedPost.SourceUserID.String(),
		RecipientID:  savedPost.RecipientID.String(),
		PostID:       savedPost.PostID.String(),
		Content:      savedPost.Content,
	}, nil
}
