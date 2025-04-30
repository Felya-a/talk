package auth_service

import (
	"context"
	. "talk/internal/lib/logger"
	model "talk/internal/services/auth/model"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	ssov1 "github.com/Felya-a/chat-app-protos/gen/go/sso"
)

type GrpcSsoProvider struct {
	grpcClient ssov1.AuthClient
}

func New(grpcUrl string) *GrpcSsoProvider {
	connection, _ := grpc.NewClient(grpcUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	grpcClient := ssov1.NewAuthClient(connection)

	return &GrpcSsoProvider{grpcClient: grpcClient}
}

func (provider *GrpcSsoProvider) Tokens(
	ctx context.Context,
	authorizationCode string,
) (tokens *model.JwtTokens, err error) {
	Log.Info("call rpc method Tokens", nil)
	response, err := provider.grpcClient.Tokens(ctx, &ssov1.TokensRequest{
		AuthorizationCode: authorizationCode,
	})
	if err != nil {
		Log.Error("error on send grpc request", Log.Err(err))
		return nil, err
	}
	Log.Info("success response from rpc method Tokens", nil)

	return &model.JwtTokens{
		AccessJwtToken:  response.AccessToken,
		RefreshJwtToken: response.RefreshToken,
	}, nil
}

func (provider *GrpcSsoProvider) Refresh(
	ctx context.Context,
	refreshToken string,
) (tokens *model.JwtTokens, err error) {
	Log.Info("call rpc method Refresh", nil)
	response, err := provider.grpcClient.Refresh(ctx, &ssov1.RefreshRequest{
		RefreshToken: refreshToken,
	})
	if err != nil {
		Log.Error("error on send grpc request", Log.Err(err))
		return nil, err
	}

	Log.Info("success response from rpc method Refresh", nil)
	return &model.JwtTokens{
		AccessJwtToken:  response.AccessToken,
		RefreshJwtToken: response.RefreshToken,
	}, nil
}

func (provider *GrpcSsoProvider) UserInfo(
	ctx context.Context,
	accessToken string,
) (user *model.UserModel, err error) {
	Log.Info("call rpc method UserInfo", nil)
	response, err := provider.grpcClient.UserInfo(ctx, &ssov1.UserInfoRequest{
		AccessToken: accessToken,
	})
	if err != nil {
		Log.Error("error on send grpc request", Log.Err(err))

		st, ok := status.FromError(err)
		if !ok {
			return nil, err
		}

		switch st.Code() {
		case codes.Unauthenticated:
			return nil, model.ErrUnathorizated
		}

		return nil, err
	}

	Log.Info("success response from rpc method UserInfo", nil)
	return &model.UserModel{
		ID:    response.UserId,
		Email: response.Email,
	}, nil
}
