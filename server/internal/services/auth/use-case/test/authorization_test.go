package auth_test

import (
	"context"
	"log/slog"
	"os"
	"testing"

	"talk/internal/config"
	auth_service "talk/internal/services/auth/model"
	sso "talk/internal/services/auth/sso"
	ssoFake "talk/internal/services/auth/sso/mock"
	usecase "talk/internal/services/auth/use-case"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("AuthorizationUseCase", Label("unit"), func() {
	var log *Logger
	var fakeUser *ssoFake.UserWithTokensModel

	var authorization usecase.AuthorizationUseCase
	var ssoProvider sso.SsoProvider
	var fakeSsoProvider *ssoFake.MockSsoProvider

	config.MustLoad()

	BeforeEach(func() {
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

		fakeSsoProvider = ssoFake.New()
		ssoProvider = fakeSsoProvider

		authorization = usecase.AuthorizationUseCase{SsoProvider: ssoProvider}

		fakeUser = FakeUser

		fakeSsoProvider.SaveUser(fakeUser)
	})

	Context("all empty", func() {
		It("should return redirectUrl", func() {
			// Arrange
			accessToken := ""
			refreshToken := ""
			authorizationCode := ""

			// Action
			redirectUrl, user, tokens, err := authorization.Execute(context.Background(), log, accessToken, refreshToken, authorizationCode)

			// Assert
			Expect(err).To(BeNil())
			Expect(tokens).To(BeNil())
			Expect(user).To(BeNil())

			Expect(redirectUrl).NotTo(BeZero())
		})
	})

	Context("authorizationCode is not empty", func() {
		It("should return user and tokens", func() {
			// Arrange
			accessToken := ""
			refreshToken := ""
			authorizationCode := fakeUser.AuthorizationCode

			// Action
			redirectUrl, user, tokens, err := authorization.Execute(context.Background(), log, accessToken, refreshToken, authorizationCode)

			// Assert
			Expect(err).To(BeNil())
			Expect(redirectUrl).To(BeZero())

			Expect(user).NotTo(BeNil())
			Expect(user.ID).NotTo(Equal(0))
			Expect(user.Email).NotTo(BeZero())
			Expect(tokens).NotTo(BeNil())
			Expect(tokens.AccessJwtToken).NotTo(BeZero())
			Expect(tokens.RefreshJwtToken).NotTo(BeZero())
		})

		It("should error if error on call Tokens", func() {
			// Arrange
			accessToken := ""
			refreshToken := ""
			authorizationCode := "invalid_authorization_code"

			// Action
			redirectUrl, user, tokens, err := authorization.Execute(context.Background(), log, accessToken, refreshToken, authorizationCode)

			// Assert
			Expect(redirectUrl).To(BeZero())
			Expect(user).To(BeNil())
			Expect(tokens).To(BeNil())

			Expect(err).NotTo(BeNil())
		})

		It("should error if error on call UserInfo", func() {
			// Arrange
			accessToken := ""
			refreshToken := ""
			authorizationCode := fakeUser.AuthorizationCode

			fakeSsoProvider.UpdateUser(0, &ssoFake.UserWithTokensModel{
				Tokens: &auth_service.JwtTokens{
					AccessJwtToken:  "invalid_access_token",
					RefreshJwtToken: fakeUser.Tokens.RefreshJwtToken,
				},
				UserModel:         fakeUser.UserModel,
				AuthorizationCode: fakeUser.AuthorizationCode,
			})

			// Action
			redirectUrl, user, tokens, err := authorization.Execute(context.Background(), log, accessToken, refreshToken, authorizationCode)

			// Assert
			Expect(redirectUrl).To(BeZero())
			Expect(user).To(BeNil())
			Expect(tokens).To(BeNil())

			Expect(err).NotTo(BeNil())
		})
	})

	Context("jwt tokens is not empty", func() {
		It("should return user", func() {
			// Arrange
			accessToken := fakeUser.Tokens.AccessJwtToken
			refreshToken := fakeUser.Tokens.RefreshJwtToken
			authorizationCode := ""

			// Action
			redirectUrl, user, tokens, err := authorization.Execute(context.Background(), log, accessToken, refreshToken, authorizationCode)

			// Assert
			Expect(err).To(BeNil())
			Expect(redirectUrl).To(BeZero())
			Expect(tokens).To(BeNil())

			Expect(user).NotTo(BeNil())
			Expect(user.ID).NotTo(Equal(0))
		})

		It("should return error if error on call UserInfo", func() {
			// Arrange
			accessToken := "invalid_access_token"
			refreshToken := fakeUser.Tokens.RefreshJwtToken
			authorizationCode := ""

			// Action
			redirectUrl, user, tokens, err := authorization.Execute(context.Background(), log, accessToken, refreshToken, authorizationCode)

			// Assert
			// Ошибки быть не должно. Принудительно вызываются методы Refresh и UserInfo
			Expect(err).To(BeNil())
			Expect(redirectUrl).To(BeZero())

			Expect(user).NotTo(BeNil())
			Expect(user.ID).NotTo(Equal(0))
			Expect(user.Email).NotTo(BeZero())
			Expect(tokens).NotTo(BeNil())
			Expect(tokens.AccessJwtToken).NotTo(BeZero())
			Expect(tokens.RefreshJwtToken).NotTo(BeZero())
		})

		It("should return error if error on call UserInfo and error on call Refresh", func() {
			// Arrange
			accessToken := "invalid_access_token"
			refreshToken := "invalid_refresh_token"
			authorizationCode := ""

			// Action
			redirectUrl, user, tokens, err := authorization.Execute(context.Background(), log, accessToken, refreshToken, authorizationCode)

			// Assert
			Expect(redirectUrl).To(BeZero())
			Expect(user).To(BeNil())
			Expect(tokens).To(BeNil())

			Expect(err).NotTo(BeNil())
		})

		It("should return error if error on call UserInfo and error on repeat call UserInfo", func() {
			// Arrange
			accessToken := "invalid_access_token"
			refreshToken := fakeUser.Tokens.RefreshJwtToken
			authorizationCode := ""

			fakeSsoProvider.UpdateUser(0, &ssoFake.UserWithTokensModel{
				Tokens: &auth_service.JwtTokens{
					AccessJwtToken:  "invalid_access_token",
					RefreshJwtToken: fakeUser.Tokens.RefreshJwtToken,
				},
				UserModel:         fakeUser.UserModel,
				AuthorizationCode: fakeUser.AuthorizationCode,
			})

			// Action
			redirectUrl, user, tokens, err := authorization.Execute(context.Background(), log, accessToken, refreshToken, authorizationCode)

			// Assert
			Expect(redirectUrl).To(BeZero())
			Expect(user).To(BeNil())
			Expect(tokens).To(BeNil())

			Expect(err).NotTo(BeNil())
		})
	})

})

func TestAuthorizationUseCase(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "AuthorizationUseCase Suite")
}
