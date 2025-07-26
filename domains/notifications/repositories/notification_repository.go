package repositories

import (
	"bootcamp-content-interaction-service/domains/notifications"
	"bootcamp-content-interaction-service/domains/notifications/entities"
	"bootcamp-content-interaction-service/infrastructures"
	"bootcamp-content-interaction-service/shared/util"
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type NotificationRepository struct {
	db          infrastructures.Database
	redisClient *redis.Client
	logger      util.Logger
}

func NewNotificationRepository(db infrastructures.Database, redisClient *redis.Client, logger util.Logger) notifications.NotificationRepository {
	return NotificationRepository{
		db:          db,
		redisClient: redisClient,
		logger:      logger,
	}
}

func (n NotificationRepository) FindAll(ctx context.Context, recipientID string) ([]*entities.Notification, error) {
	key :="post_notifications:" + recipientID
	var notifications []*entities.Notification

	//get from redis
	cached, err := n.redisClient.LRange(ctx, key, 0, -1).Result()
	if err == nil && len(cached) > 0 {
		n.logger.Info("Cache hit - returning notifications from redis",
			zap.String("cache_key", key),
		)
		for _, jsonItem := range cached {
			var notification entities.Notification
			if err := json.Unmarshal([]byte(jsonItem), &notification); err == nil {
				notifications = append(notifications, &notification)
			}
		}
		return notifications, nil
	}

	//else , get from db
	result := n.db.GetInstance().WithContext(ctx).Where("recipient_id = ?", recipientID).Order("created_at DESC").Find(&notifications)
	n.logger.Info("Get all data from DB")
	if result.Error != nil {
		return nil, result.Error
	}

	// set in redis
	for _, notification := range notifications {
		notifJSON, _ := json.Marshal(notification)
		_ = n.redisClient.LPush(ctx, key, notifJSON).Err()
		n.logger.Info("Set notif in cache",
			zap.String("notif_json", string(notifJSON)),
		)
	}
	_ = n.redisClient.LTrim(ctx, key, 0, 99)
	_ = n.redisClient.Expire(ctx, key,  7*24*time.Hour)
	return notifications, nil
}

func (n NotificationRepository) SaveNotification(ctx context.Context, notif *entities.Notification) (*entities.Notification, error) {

	notifmodel := &entities.Notification{
		ID:           uuid.New(),
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

	n.logger.Info("Notification saved to DB",
		zap.String("notification_id", notifmodel.ID.String()),
		zap.String("recipient_id", notifmodel.RecipientID.String()),
	)

	inboxKey := "post_notifications:" + notifmodel.RecipientID.String()
	notifJSON, err := json.Marshal(notifmodel)
	if err != nil {
		n.logger.Warn("Redis marshal failed", zap.Error(err))
		return notifmodel, nil
	}

	if err := n.redisClient.LPush(ctx, inboxKey, notifJSON).Err(); err != nil {
		n.logger.Warn("Redis LPUSH failed", zap.Error(err))
	} else {
		_ = n.redisClient.Expire(ctx, inboxKey, 7*24*time.Hour).Err()
		n.logger.Info("Notification stored in Redis",
			zap.String("redis_key", inboxKey),
			zap.String("recipient_id", notifmodel.RecipientID.String()),
		)
	}

	return notifmodel, nil
}
