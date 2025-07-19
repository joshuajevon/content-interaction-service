package http

import (
	"bootcamp-content-interaction-service/domains/comments"
	"bootcamp-content-interaction-service/domains/comments/models/request"
	"bootcamp-content-interaction-service/domains/comments/models/response"
	"bootcamp-content-interaction-service/domains/users/models/dto"
	"net/http"

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

	err := c.ShouldBindBodyWithJSON(req)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			gin.H{
				"error": "Msg Null",
			},
		)
	}

	authUser, ok := c.Request.Context().Value("user").(*dto.AuthUserDto)
	if !ok {
		c.JSON(http.StatusUnauthorized,
			gin.H{
				"error": "Unauthorized",
			},
		)
	}

	userId := authUser.UserId
	postId := c.Param("id")

	err = h.uc.CreateComment(ctx, userId, postId, req.Msg, nil)

	if err != nil {
		c.JSON(http.StatusBadRequest,
			gin.H{
				"error": err,
			},
		)
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

	err := c.ShouldBindBodyWithJSON(req)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			gin.H{
				"error": "Msg Null",
			},
		)
	}

	authUser, ok := c.Request.Context().Value("user").(*dto.AuthUserDto)
	if !ok {
		c.JSON(http.StatusUnauthorized,
			gin.H{
				"error": "Unauthorized",
			},
		)
	}

	userId := authUser.UserId
	postId := c.Param("id")

	if req.ReplyId == nil {
		err = h.uc.UpdateComment(ctx, userId, postId, req.Msg, nil)
	} else {
		replyIdStr := req.ReplyId.String()
		err = h.uc.UpdateComment(ctx, userId, postId, req.Msg, &replyIdStr)
	}

	if err != nil {
		c.JSON(http.StatusBadRequest,
			gin.H{
				"error": err,
			},
		)
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

	err := c.ShouldBindBodyWithJSON(req)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			gin.H{
				"error": "Msg Null",
			},
		)
	}

	authUser, ok := c.Request.Context().Value("user").(*dto.AuthUserDto)
	if !ok {
		c.JSON(http.StatusUnauthorized,
			gin.H{
				"error": "Unauthorized",
			},
		)
	}

	userId := authUser.UserId
	postId := c.Param("id")

	if req.ReplyId == nil {
		c.JSON(http.StatusBadRequest,
			gin.H{
				"error": "reply_id NULL",
			},
		)
	} else {
		replyIdStr := req.ReplyId.String()
		err = h.uc.ReplyComment(ctx, userId, postId, replyIdStr, req.Msg)
		if err != nil {
			c.JSON(http.StatusBadRequest,
				gin.H{
					"error": err,
				},
			)
		}
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
	}

	postId := c.Param("id")

	comment, err := h.uc.FindAllComment(ctx, postId)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			gin.H{
				"error": err,
			},
		)
	}

	if comment == nil || len(*comment) == 0{
		c.JSON(http.StatusOK,
			gin.H{
				"message": "There's No Comment",
				"data": comment,
			},
		)
	}

	for _, i := range *comment{
		response := response.CommentResponse{
			ID: i.ID,
			UserID: i.UserID,
			ReplyId: i.ReplyId,
			CreatedAt: i.CreatedAt,
			UpdatedAt: i.UpdatedAt,
			Msg: i.Msg,
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
	}

	commentId := c.Param("comments_id")
	parsedCommentId, err := uuid.Parse(commentId)
	if err != nil{
		c.JSON(http.StatusInternalServerError,
			gin.H{
				"error": "Error while parsing uuid",
			},
		)
	}

	err = h.uc.DeleteComment(ctx, parsedCommentId)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			gin.H{
				"error": "Failed to delete comment",
			},
		)
	}

	c.JSON(http.StatusOK, gin.H{
		"message" : "SUCCESS",
	})
}
