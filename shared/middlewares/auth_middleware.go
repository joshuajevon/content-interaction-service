package middlewares

import (
	"bootcamp-content-interaction-service/domains/users/models/dto"
	"bootcamp-content-interaction-service/shared/constant"
	"bootcamp-content-interaction-service/shared/models/responses"
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		accessToken := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")

		token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
			return constant.JWT_SECRET, nil
		})
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, responses.BasicResponse{Error: "Unauthorized"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, responses.BasicResponse{Error: "Invalid claims"})
			return
		}

		authUser := &dto.AuthUserDto{
			UserId: fmt.Sprintf("%v", claims["id"]),
			Name:   fmt.Sprintf("%v", claims["name"]),
			Email:  fmt.Sprintf("%v", claims["email"]),
		}

		ctx := context.WithValue(c.Request.Context(), "user", authUser)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
