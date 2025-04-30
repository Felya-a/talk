package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"talk/internal/config"
	. "talk/internal/lib/logger"
	authService "talk/internal/services/auth"
	"talk/internal/transport/http/middleware"
	"talk/internal/transport/http/router"

	"github.com/gin-gonic/gin"
)

type HttpTransport struct {
	httpServer *http.Server
	port       string
}

func New(
	port string,
	authService *authService.AuthService,
) *HttpTransport {
	gin.SetMode(getGinMode())
	handler := gin.Default()
	handler.Use(middleware.CORSMiddleware()) // TODO: Возможно на проде нужно отключать

	router.SetupRoutes(handler, authService)

	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: handler,
	}

	return &HttpTransport{
		httpServer,
		port,
	}
}

func (t *HttpTransport) MustRun() {
	if err := t.run(); err != nil {
		panic(err)
	}
}

func (t *HttpTransport) run() error {
	const op = "http_app.Run"

	Log.Info("http server is running", LogFields{"op": op, "port": t.port, "addr": t.httpServer.Addr})
	if err := t.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		Log.Error("error on start http server", Log.Err(err))
		return fmt.Errorf("%s:%w", op, err)
	}

	return nil
}

func getGinMode() string {
	var mode string

	switch config.Get().Env {
	case "local", "test":
		mode = "debug"
	case "stage", "prod":
		mode = "release"
	default:
		mode = "release"
	}

	return mode
}

func (t *HttpTransport) Stop() {
	Log.Info("stopping http server", nil)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := t.httpServer.Shutdown(ctx); err != nil {
		Log.Error("server forced to shutdown", Log.Err(err))
	}
}
