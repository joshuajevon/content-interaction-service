package wizards

import (
	"bootcamp-content-interaction-service/shared/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterServer(router *gin.Engine) {
	post := router.Group("/posts")
	{
		post.Use(middlewares.AuthMiddleware())
		post.POST("/:id/likes", LikesHttp.LikePost)
		post.POST("/:id/dislikes", LikesHttp.DislikePost)

		post.GET("/:id/comments", CommentsHttp.FindAllComment)
		post.POST("/:id/comments", CommentsHttp.CreateComment)
		post.POST("/:id/comments/:comments_id", CommentsHttp.UpdateComment)
		post.POST("/:id/comments/:comments_id/reply", CommentsHttp.ReplyComment)
		post.DELETE("/:id/comments/:comments_id", CommentsHttp.DeleteComment)

	}
}
