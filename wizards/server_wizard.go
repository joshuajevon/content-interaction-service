package wizards

import (

	"github.com/gin-gonic/gin"
)

func RegisterServer(router *gin.Engine) {
	post := router.Group("/posts")
	{
		post.POST("/view/likes/:id", LikesHttp.LikePost)
		post.POST("/view/dislikes/:id", LikesHttp.DislikePost)
	}
}
