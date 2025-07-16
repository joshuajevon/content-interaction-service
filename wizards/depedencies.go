package wizards

import (
	"bootcamp-content-interaction-service/config"
	postHttp "bootcamp-content-interaction-service/domains/posts/handlers/http"
	postRepo "bootcamp-content-interaction-service/domains/posts/repositories"
	postUc "bootcamp-content-interaction-service/domains/posts/usecases"
	"bootcamp-content-interaction-service/infrastructures"
)

var (
	Config             = config.GetConfig()
	PostgresDatabase   = infrastructures.NewPostgresDatabase(Config)
	PostRepository = postRepo.NewPostRepository(PostgresDatabase)
	PostUseCase = postUc.NewPostUseCase(PostRepository)
	PostHttp = postHttp.NewPostHttp(PostUseCase)
)
