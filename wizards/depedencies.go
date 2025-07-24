package wizards

import (
	"bootcamp-content-interaction-service/config"
	commentsHttp "bootcamp-content-interaction-service/domains/comments/handlers/http"
	commentsRepository "bootcamp-content-interaction-service/domains/comments/repositories"
	commentsUc "bootcamp-content-interaction-service/domains/comments/usecases"
	likesHttp "bootcamp-content-interaction-service/domains/likes/handlers/http"
	likesRepository "bootcamp-content-interaction-service/domains/likes/repositories"
	likesUc "bootcamp-content-interaction-service/domains/likes/usecases"
	postHttp "bootcamp-content-interaction-service/domains/posts/handlers/http"
	postRepo "bootcamp-content-interaction-service/domains/posts/repositories"
	postUc "bootcamp-content-interaction-service/domains/posts/usecases"
	notificationHttp "bootcamp-content-interaction-service/domains/notifications/handlers/http"
	notificationRepo "bootcamp-content-interaction-service/domains/notifications/repositories"
	notificationUc "bootcamp-content-interaction-service/domains/notifications/usecases"
	"bootcamp-content-interaction-service/infrastructures"
	"bootcamp-content-interaction-service/shared/util"
)

var (
	Config              = config.GetConfig()
	PostgresDatabase    = infrastructures.NewPostgresDatabase(Config)
	RedisClient 	    = infrastructures.InitRedis()
	LoggerInstance, _ 	= util.NewLogger()
	
	UserGraphService    = postHttp.NewUserGraphHTTP(Config.Server.UserGraphBaseURL)

	LikesRepository     = likesRepository.NewLikesRepository(PostgresDatabase)
	LikesUseCase        = likesUc.NewLikesUseCase(LikesRepository)
	LikesHttp           = likesHttp.NewLikesHandler(LikesUseCase)

	CommentsRepository  = commentsRepository.NewCommentsRepository(PostgresDatabase, RedisClient, LoggerInstance)
	CommentsUseCase     = commentsUc.NewCommentsUseCase(CommentsRepository)
	CommentsHttp        = commentsHttp.NewLikesHandler(CommentsUseCase)

	PostRepository      = postRepo.NewPostRepository(PostgresDatabase, RedisClient, LoggerInstance)
    PostUseCase         = postUc.NewPostUseCase(PostRepository, UserGraphService, NotificationRepository)
	PostHttp            = postHttp.NewPostHttp(PostUseCase)

	NotificationRepository 	= notificationRepo.NewNotificationRepository(PostgresDatabase, RedisClient, LoggerInstance)
	NotificationUseCase  	= notificationUc.NewNotificationUseCase(NotificationRepository)
	NotificationHttp		= notificationHttp.NewNotificationHttp(NotificationUseCase)
)
