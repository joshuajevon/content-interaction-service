package util

import (
	"bootcamp-content-interaction-service/domains/users/models/dto"
	"context"
	"errors"
)

func GetAuthUser(ctx context.Context) (*dto.AuthUserDto, error) {
	userRaw := ctx.Value("user")
	user, ok := userRaw.(*dto.AuthUserDto)
	if !ok || user == nil {
		return nil, errors.New("unauthorized: user not found in context")
	}
	return user, nil
}
