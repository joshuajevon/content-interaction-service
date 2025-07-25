package http

import (
	"bootcamp-content-interaction-service/domains/notifications"
	"bootcamp-content-interaction-service/domains/notifications/models/requests"
	"bootcamp-content-interaction-service/shared/models/responses"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type NotificationHttp struct {
	notifUc notifications.NotificationUseCase
}

func NewNotificationHttp(notifUc notifications.NotificationUseCase) *NotificationHttp {
	return &NotificationHttp{
		notifUc: notifUc,
	}
}

func (handler *NotificationHttp) ViewAllNotification(c *gin.Context) {
	ctx := c.Request.Context()

	result, err := handler.notifUc.FindAllNotification(ctx)

    if err != nil {
        c.JSON(http.StatusInternalServerError, responses.BasicResponse{Error: err.Error()})
        return
    }

    c.JSON(http.StatusOK, result)
}

func (handler *NotificationHttp) CreatePostNotification(c *gin.Context) {
	ctx := c.Request.Context()
	var req requests.PostNotificationRequest

	if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, responses.BasicResponse{Error: err.Error()})
        return
    }

	validate := validator.New()
	err := validate.StructCtx(ctx, req)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, responses.BasicResponse{
			Error: err.Error(),
		})
		return
	}

	result, err := handler.notifUc.NotifyNewPost(ctx, &req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, responses.BasicResponse{
			Error: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, result)
}