package main

import (
	comments "bootcamp-content-interaction-service/domains/comments/entities"
	likes "bootcamp-content-interaction-service/domains/likes/entities"
	users "bootcamp-content-interaction-service/domains/users/entities"
	"bootcamp-content-interaction-service/wizards"
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	wizards.PostgresDatabase.GetInstance().AutoMigrate(
		&users.User{},
		&likes.Likes{},
		&comments.Comments{},
	)

	router := gin.Default()

	wizards.RegisterServer(router)

	router.Run(fmt.Sprintf(":%d", wizards.Config.Server.Port))
}