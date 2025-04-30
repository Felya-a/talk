package auth_service

import (
	"context"
	model "talk/internal/services/auth/model"
)

type SsoProvider interface {
	Tokens(
		ctx context.Context,
		authorizationCode string,
	) (tokens *model.JwtTokens, err error)
	Refresh(
		ctx context.Context,
		refreshToken string,
	) (tokens *model.JwtTokens, err error)
	UserInfo(
		ctx context.Context,
		accessToken string,
	) (user *model.UserModel, err error)
}
