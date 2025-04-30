package http_handlers_v1

import (
	"net/http"
	"net/url"
	"talk/internal/config"
	. "talk/internal/http/handlers"
	. "talk/internal/lib/logger"
	authService "talk/internal/services/auth"
	models "talk/internal/services/auth/model"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetAuthHandler(
	authService *authService.AuthService,
) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var dto AuthRequestDto
		var err error

		log := NewLogger()
		log.AttachFields(LogFields{"requestid": uuid.New().String()})

		dto.AccessToken, err = ctx.Cookie("access_token")
		if err != nil && err != http.ErrNoCookie {
			response := ErrorResponse{
				Status:  "error",
				Message: "error on auth",
				Error:   "error on read access_token from cookie",
			}
			log.Error("error on read access_token from cookie", Log.Err(err))
			ctx.JSON(400, response)
		}

		dto.RefreshToken, err = ctx.Cookie("refresh_token")
		if err != nil && err != http.ErrNoCookie {
			response := ErrorResponse{
				Status:  "error",
				Message: "error on auth",
				Error:   "error on read refresh_token from cookie",
			}
			log.Error("error on read refresh_token from cookie", Log.Err(err))
			ctx.JSON(400, response)
		}

		dto.AuthorizationCode = ctx.Query("authorization_code")

		redirectUrl, user, tokens, err := authService.Auth(ctx, dto.AccessToken, dto.RefreshToken, dto.AuthorizationCode)
		if err != nil {
			log.Info("error on auth", Log.Err(err))
			response := ErrorResponse{
				Status:  "error",
				Message: "error on auth",
				Error:   "internal error",
			}

			if models.IsDefinedError(err) {
				response.Error = err.Error()
				// TODO: По хорошему сделать в контроллере дополнительную мапу ошибка-код
				ctx.JSON(400, response)
				return
			}

			ctx.JSON(500, response)
			return
		}

		if tokens != nil {
			ctx.SetCookie("access_token", tokens.AccessJwtToken, 30*24*60*60, "/", "", true, true)
			ctx.SetCookie("refresh_token", tokens.RefreshJwtToken, 30*24*60*60, "/", "", true, true)
		}

		if redirectUrl != "" {
			queryParams := url.Values{}
			queryParams.Add("redirect_url", config.Get().Talk.HttpClientUrl)

			// ctx.Redirect(307, redirectUrl+"?"+queryParams.Encode())
			response := SuccessResponse{
				Status:  "ok",
				Message: "success login",
				Data:    AuthResponseDto{RedirectUrl: redirectUrl + "?" + queryParams.Encode()},
			}
			ctx.JSON(200, response)
			return
		}

		response := SuccessResponse{
			Status:  "ok",
			Message: "success login",
			Data:    AuthResponseDto{UserData: UserData{UserId: user.ID, Email: user.Email}},
		}
		ctx.JSON(200, response)
	}
}
