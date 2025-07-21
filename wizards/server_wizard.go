package wizards

import (
	"github.com/gin-gonic/gin"
)

func RegisterServer(router *gin.Engine) {
  
  api := router.Group("/v1")
	{
    post := api.Group("/posts")
    {
      post.GET("/view", PostHttp.ViewAllPost)
			post.GET("/view/:id", PostHttp.ViewPostById)
			post.GET("/view/feed/:id", PostHttp.ViewPersonalFeed)
      
      post.Use(middlewares.AuthMiddleware())
      post.POST("/:id/likes", LikesHttp.LikePost)
      post.POST("/:id/dislikes", LikesHttp.DislikePost)

      post.GET("/:id/comments", CommentsHttp.FindAllComment)
      post.POST("/:id/comments", CommentsHttp.CreateComment)
      post.POST("/:id/comments/:comments_id", CommentsHttp.UpdateComment)
      post.POST("/:id/comments/:comments_id/reply", CommentsHttp.ReplyComment)
      post.DELETE("/:id/comments/:comments_id", CommentsHttp.DeleteComment)

			post.POST("/create", PostHttp.CreatePost)
			post.GET("/view/user", PostHttp.ViewAllPostByUserId)
			post.DELETE("/delete/:id", PostHttp.DeletePost)
			post.PATCH("/update/:id", PostHttp.UpdatePost)
		}
	}
}
