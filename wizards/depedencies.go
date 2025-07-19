package wizards

import (
	"bootcamp-content-interaction-service/config"
	likesHttp "bootcamp-content-interaction-service/domains/likes/handlers/http"
	likesRepository "bootcamp-content-interaction-service/domains/likes/repositories"
	likesUc "bootcamp-content-interaction-service/domains/likes/usecases"

	commentsHttp "bootcamp-content-interaction-service/domains/comments/handlers/http"
	commentsRepository "bootcamp-content-interaction-service/domains/comments/repositories"
	commentsUc "bootcamp-content-interaction-service/domains/comments/usecases"
	"bootcamp-content-interaction-service/infrastructures"
)

var (
	Config           = config.GetConfig()
	PostgresDatabase = infrastructures.NewPostgresDatabase(Config)

	LikesRepository  = likesRepository.NewLikesRepository(PostgresDatabase)
	LikesUseCase     = likesUc.NewLikesUseCase(LikesRepository)
	LikesHttp        = likesHttp.NewLikesHandler(LikesUseCase)

	CommentsRepository = commentsRepository.NewCommentsRepository(PostgresDatabase)
	CommentsUseCase     = commentsUc.NewCommentsUseCase(CommentsRepository)
	CommentsHttp        = commentsHttp.NewLikesHandler(CommentsUseCase)
)
