package wizards

import (
	"bootcamp-content-interaction-service/shared/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterServer(router *gin.Engine) {

	api := router.Group("/v1")
	{
		post := api.Group("/posts")
		{
			post.GET("/view", PostHttp.ViewAllPost)
			post.GET("/view/:id", PostHttp.ViewPostById)

			post.Use(middlewares.AuthMiddleware())
			post.POST("/create", PostHttp.CreatePost)
			post.GET("/view/user", PostHttp.ViewAllPostByUserId)
			post.DELETE("/delete/:id", PostHttp.DeletePost)
		}

	}
}
