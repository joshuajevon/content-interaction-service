package http

import (
	"bootcamp-content-interaction-service/domains/comments"
	"bootcamp-content-interaction-service/domains/comments/models/request"
	"bootcamp-content-interaction-service/domains/comments/models/response"
	"bootcamp-content-interaction-service/domains/users/models/dto"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CommentsHttp struct {
	uc comments.CommentsUseCase
}

func NewLikesHandler(uc comments.CommentsUseCase) CommentsHttp {
	return CommentsHttp{uc: uc}
}

func (h *CommentsHttp) CreateComment(c *gin.Context) {
	var req request.CommentRequest
	ctx := c.Request.Context()

	err := c.ShouldBindBodyWithJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			gin.H{
				"error": "Msg Null",
			},
		)
		return
	}

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

	err = h.uc.CreateComment(ctx, userId, postId, req.Msg, nil)

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

func (h *CommentsHttp) UpdateComment(c *gin.Context) {
	var req request.CommentRequest
	ctx := c.Request.Context()

	err := c.ShouldBindBodyWithJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			gin.H{
				"error": "Msg Null",
			},
		)
		return
	}

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
	commentId := strings.TrimPrefix(c.Param("comments_id"), ":")

	err = h.uc.UpdateComment(ctx, commentId, userId, req.Msg)

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

func (h *CommentsHttp) ReplyComment(c *gin.Context) {
	var req request.CommentRequest
	ctx := c.Request.Context()

	err := c.ShouldBindBodyWithJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			gin.H{
				"error": "Msg Null",
			},
		)
		return
	}

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
	commentId := strings.TrimPrefix(c.Param("comments_id"), ":")

	err = h.uc.ReplyComment(ctx, commentId, userId, postId, req.Msg)
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

func (h *CommentsHttp) FindAllComment(c *gin.Context) {
	var res []response.CommentResponse
	ctx := c.Request.Context()
	_, ok := c.Request.Context().Value("user").(*dto.AuthUserDto)
	if !ok {
		c.JSON(http.StatusUnauthorized,
			gin.H{
				"error": "Unauthorized",
			},
		)
		return
	}

	postId := strings.TrimPrefix(c.Param("id"), ":")

	comment, err := h.uc.FindAllComment(ctx, postId)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			gin.H{
				"error": err,
			},
		)
		return
	}

	if comment == nil || len(*comment) == 0 {
		c.JSON(http.StatusOK,
			gin.H{
				"message": "There's No Comment",
				"data":    comment,
			},
		)
		return
	}

	for _, i := range *comment {
		response := response.CommentResponse{
			ID:        i.ID,
			UserID:    i.UserID,
			ReplyId:   i.ReplyId,
			CreatedAt: i.CreatedAt,
			UpdatedAt: i.UpdatedAt,
			Msg:       i.Msg,
		}
		res = append(res, response)
	}

	c.JSON(http.StatusOK, res)
}

func (h *CommentsHttp) DeleteComment(c *gin.Context) {
	ctx := c.Request.Context()
	_, ok := c.Request.Context().Value("user").(*dto.AuthUserDto)
	if !ok {
		c.JSON(http.StatusUnauthorized,
			gin.H{
				"error": "Unauthorized",
			},
		)
		return
	}

	commentId := strings.TrimPrefix(c.Param("comments_id"), ":")
	parsedCommentId, err := uuid.Parse(commentId)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			gin.H{
				"error": "Error while parsing uuid",
			},
		)
		return
	}

	err = h.uc.DeleteComment(ctx, parsedCommentId)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			gin.H{
				"error": "Failed to delete comment",
			},
		)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "SUCCESS",
	})
}
