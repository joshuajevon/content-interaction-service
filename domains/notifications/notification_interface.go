package notifications

import (
	"bootcamp-content-interaction-service/domains/notifications/entities"
	"bootcamp-content-interaction-service/domains/notifications/models/requests"
	"bootcamp-content-interaction-service/domains/notifications/models/responses"
	"context"
)

type NotificationUseCase interface {
	NotifyNewPost(ctx context.Context, request *requests.PostNotificationRequest) (*responses.PostNotificationResponse, error)
}

type NotificationRepository interface {
	SaveNotification(ctx context.Context, notif *entities.Notification) (*entities.Notification, error)
}