package wizards

import (

	"github.com/gin-gonic/gin"
)

func RegisterServer(router *gin.Engine) {
	post := router.Group("/posts")
	{
		post.POST("/create", PostHttp.CreatePost)
		post.GET("/view", PostHttp.ViewAllPost)
	}
}
