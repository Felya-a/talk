package auth_test

import (
	"context"
	"errors"
	. "talk/internal/lib/logger"
	model "talk/internal/services/auth/model"
)

// TODO: реализовать мок ssoProvider и реализовать authService через TDD
// 12.01.25

type UserWithTokensModel struct {
	Tokens            *model.JwtTokens
	AuthorizationCode string
	model.UserModel
}

type MockSsoProvider struct {
	users []*UserWithTokensModel
}

func New() *MockSsoProvider {
	return &MockSsoProvider{}
}

func (provider *MockSsoProvider) Tokens(
	ctx context.Context,
	log *Logger,
	authorizationCode string,
) (tokens *model.JwtTokens, err error) {
	if authorizationCode == "invalid_authorization_code" {
		return nil, errors.New("fake error")
	}

	for _, user := range provider.users {
		if user.AuthorizationCode == authorizationCode {
			return user.Tokens, nil
		}
	}

	return nil, errors.New("user not found")
}

func (provider *MockSsoProvider) Refresh(
	ctx context.Context,
	log *Logger,
	refreshToken string,
) (tokens *model.JwtTokens, err error) {
	for _, user := range provider.users {
		if user.Tokens.RefreshJwtToken == refreshToken {
			return user.Tokens, nil
		}
	}

	return nil, errors.New("user not found")
}

func (provider *MockSsoProvider) UserInfo(
	ctx context.Context,
	log *Logger,
	accessToken string,
) (user *model.UserModel, err error) {
	if accessToken == "invalid_access_token" {
		return nil, errors.New("fake error")
	}

	for _, user := range provider.users {
		if user.Tokens.AccessJwtToken == accessToken {
			return &user.UserModel, nil
		}
	}

	return nil, errors.New("user not found")
}

/* FOR MOCK ONLY */

func (provider *MockSsoProvider) SaveUser(
	user *UserWithTokensModel,
) {
	provider.users = append(provider.users, user)
}

func (provider *MockSsoProvider) UpdateUser(
	index int,
	user *UserWithTokensModel,
) {
	provider.users[index] = user
}
