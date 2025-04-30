package http_routes_v1

import (
	handlers "talk/internal/http/handlers/v1"
	authService "talk/internal/services/auth"

	"github.com/gin-gonic/gin"
)

func SetupV1Routes(
	r *gin.RouterGroup,
	authService *authService.AuthService,
) {
	r.GET("/auth", handlers.GetAuthHandler(authService))

	/* FOR DEBUG ONLY */
	r.GET("/redirect", func(ctx *gin.Context) { ctx.Redirect(301, "https://google.com") })
	r.POST("/redirect", func(ctx *gin.Context) { ctx.Redirect(301, "https://google.com") })
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
}
