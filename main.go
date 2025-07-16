package main

import (
	users "bootcamp-content-interaction-service/domains/users/entities"
	posts "bootcamp-content-interaction-service/domains/posts/entities"
	"bootcamp-content-interaction-service/wizards"
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	wizards.PostgresDatabase.GetInstance().AutoMigrate(
		&users.User{},
		&posts.Post{},
	)

	router := gin.Default()

	wizards.RegisterServer(router)

	router.Run(fmt.Sprintf(":%d", wizards.Config.Server.Port))
}