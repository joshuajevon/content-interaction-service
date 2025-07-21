package repositories

import (
	"bootcamp-content-interaction-service/domains/notifications"
	"bootcamp-content-interaction-service/domains/notifications/entities"
	"bootcamp-content-interaction-service/infrastructures"
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type NotificationRepository struct {
	db infrastructures.Database
	redisClient *redis.Client
}

func NewNotificationRepository(db infrastructures.Database, redisClient *redis.Client) notifications.NotificationRepository {
	return NotificationRepository{
		db: db,
		redisClient: redisClient,
	}
}

func (n NotificationRepository) SaveNotification(ctx context.Context, notif *entities.Notification) (*entities.Notification, error) {
	var logger = zap.NewExample()

	notifmodel := &entities.Notification{
		ID: uuid.New(),
		SourceUserID: notif.SourceUserID,
		RecipientID: notif.RecipientID,
		PostID: notif.PostID,
		Type: notif.Type,
		Content: notif.Content,
		CreatedAt: time.Now(),
	}

	result := n.db.GetInstance().WithContext(ctx).Create(notifmodel)
	if result.Error != nil {
		return nil, result.Error
	}

	logger.Info("Saving notification to DB with id: " + notif.ID.String())

	notifJSON, err := json.Marshal(notifmodel)
	if err != nil {
        logger.Warn("Redis marshal failed:" + err.Error())
        return notifmodel, nil
    }

	if err := n.redisClient.Publish(ctx, "post_notifications", notifJSON).Err(); err != nil {
        logger.Warn("Redis publish failed", zap.Error(err))
    } else {
        logger.Info("Notification published", zap.String("recipient", notif.RecipientID.String()))
    }
	return notifmodel, nil
}
