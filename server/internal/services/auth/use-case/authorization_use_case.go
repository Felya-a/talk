package usecase

import (
	"context"
	"talk/internal/config"
	. "talk/internal/lib/logger"
	models "talk/internal/services/auth/model"
	sso "talk/internal/services/auth/sso"
)

type AuthorizationUseCase struct {
	SsoProvider sso.SsoProvider
}

func (uc *AuthorizationUseCase) Execute(
	ctx context.Context,
	accessToken string,
	refreshToken string,
	authorizationCode string,
) (redirectUrl string, user *models.UserModel, jwtTokens *models.JwtTokens, err error) {
	var tokens *models.JwtTokens

	if accessToken == "" &&
		refreshToken == "" &&
		authorizationCode == "" {
		return config.Get().Sso.HttpClientUrl, nil, nil, nil
	}

	if authorizationCode != "" {
		tokens, err = uc.SsoProvider.Tokens(ctx, authorizationCode)
		if err != nil {
			logTokensError(authorizationCode, err)
			return "", nil, nil, err
		}

		user, err = uc.SsoProvider.UserInfo(ctx, tokens.AccessJwtToken)
		if err != nil {
			logUserInfoError(tokens.AccessJwtToken, err)
			return "", nil, nil, err
		}
	}

	if accessToken != "" ||
		refreshToken != "" {
		user, err = uc.SsoProvider.UserInfo(ctx, accessToken)
		if err != nil {
			logUserInfoError(accessToken, err)
			tokens, err = uc.SsoProvider.Refresh(ctx, refreshToken)
			if err != nil {
				logRefreshError(refreshToken, err)
				return "", nil, nil, err
			}

			user, err = uc.SsoProvider.UserInfo(ctx, tokens.AccessJwtToken)
			if err != nil {
				logUserInfoError(tokens.AccessJwtToken, err)
				return "", nil, nil, err
			}
		}
	}

	return "", user, tokens, nil
}

func logTokensError(authorizationCode string, err error) {
	Log.Error(
		"error on call rpc method Tokens",
		LogFields{"authorizationCode": authorizationCode, "error": err},
	)
}

func logRefreshError(refreshToken string, err error) {
	Log.Error(
		"error on call rpc method Refresh",
		LogFields{"refresh_token": refreshToken, "error": err},
	)
}

func logUserInfoError(accessToken string, err error) {
	Log.Error(
		"error on call rpc method UserInfo",
		LogFields{"access_token": accessToken, "error": err},
	)
}
