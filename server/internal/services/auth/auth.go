package auth_service

import (
	"context"

	"talk/internal/config"
	models "talk/internal/services/auth/model"
	sso "talk/internal/services/auth/sso"
	usecase "talk/internal/services/auth/use-case"
)

type AuthService struct {
	ssoProvider   sso.SsoProvider
	authorization usecase.AuthorizationUseCase
}

func NewAuthService() *AuthService {
	grpcSsoProvider := sso.New(config.Get().Sso.GrpcServerUrl)
	authorization := usecase.AuthorizationUseCase{SsoProvider: grpcSsoProvider}

	return &AuthService{
		ssoProvider:   grpcSsoProvider,
		authorization: authorization,
	}
}

func (s *AuthService) Auth(
	ctx context.Context,
	accessToken string,
	refreshToken string,
	authorizationCode string,
) (redirectUrl string, user *models.UserModel, jwtTokens *models.JwtTokens, err error) {
	return s.authorization.Execute(ctx, accessToken, refreshToken, authorizationCode)
}
