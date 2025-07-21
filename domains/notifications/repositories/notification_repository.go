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
		ID:          uuid.New(),
		SourceUserID: notif.SourceUserID,
		RecipientID:  notif.RecipientID,
		PostID:       notif.PostID,
		Type:         notif.Type,
		Content:      notif.Content,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	result := n.db.GetInstance().WithContext(ctx).Create(notifmodel)
	if result.Error != nil {
		return nil, result.Error
	}

	logger.Info("Notification saved to DB",
		zap.String("notification_id", notifmodel.ID.String()),
		zap.String("recipient_id", notifmodel.RecipientID.String()),
	)

	inboxKey := "post_notifications:" + notifmodel.RecipientID.String()
	notifJSON, err := json.Marshal(notifmodel)
	if err != nil {
		logger.Warn("Redis marshal failed", zap.Error(err))
		return notifmodel, nil
	}

	if err := n.redisClient.LPush(ctx, inboxKey, notifJSON).Err(); err != nil {
		logger.Warn("Redis LPUSH failed", zap.Error(err))
	} else {
		_ = n.redisClient.Expire(ctx, inboxKey, 7*24*time.Hour).Err()
		logger.Info("Notification stored in Redis",
			zap.String("redis_key", inboxKey),
			zap.String("recipient_id", notifmodel.RecipientID.String()),
		)
	}

	return notifmodel, nil
}