package main

import (
	users "bootcamp-content-interaction-service/domains/users/entities"
	"bootcamp-content-interaction-service/wizards"
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	wizards.PostgresDatabase.GetInstance().AutoMigrate(
		&users.User{},
	)

	router := gin.Default()

	wizards.RegisterServer(router)

	router.Run(fmt.Sprintf(":%d", wizards.Config.Server.Port))
}