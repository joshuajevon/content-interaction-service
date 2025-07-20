package http

import (
	"bootcamp-content-interaction-service/domains/likes"
	"bootcamp-content-interaction-service/domains/users/models/dto"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type LikesHttp struct {
	uc likes.LikesUseCase
}

func NewLikesHandler(uc likes.LikesUseCase) LikesHttp {
	return LikesHttp{uc: uc}
}

func (h *LikesHttp) LikePost(c *gin.Context) {
	ctx := c.Request.Context()
	authUser, ok := c.Request.Context().Value("user").(*dto.AuthUserDto)
	if !ok {
		c.JSON(http.StatusUnauthorized,
			gin.H{
				"error": "Unauthorized",
			},
		)
		return
	}

	userId := authUser.UserId
	postId := strings.TrimPrefix(c.Param("id"), ":")

	err := h.uc.LikePost(ctx, userId, postId)

	if err != nil {
		c.JSON(http.StatusBadRequest,
			gin.H{
				"error": err,
			},
		)
		return
	}

	c.JSON(http.StatusOK,
		gin.H{
			"message": "SUCCESS",
		},
	)
}

func (h *LikesHttp) DislikePost(c *gin.Context) {
	ctx := c.Request.Context()
	authUser, ok := c.Request.Context().Value("user").(*dto.AuthUserDto)
	if !ok {
		c.JSON(http.StatusUnauthorized,
			gin.H{
				"error": "Unauthorized",
			},
		)
		return
	}

	userId := authUser.UserId
	postId := strings.TrimPrefix(c.Param("id"), ":")

	err := h.uc.DislikePost(ctx, userId, postId)

	if err != nil {
		c.JSON(http.StatusBadRequest,
			gin.H{
				"error": err,
			},
		)
		return
	}

	c.JSON(http.StatusOK,
		gin.H{
			"message": "SUCCESS",
		},
	)
}
