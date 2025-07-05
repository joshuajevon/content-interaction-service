package wizards

import (
	"bootcamp-content-interaction-service/config"
	"bootcamp-content-interaction-service/infrastructures"
)

var (
	Config             = config.GetConfig()
	PostgresDatabase   = infrastructures.NewPostgresDatabase(Config)
)
