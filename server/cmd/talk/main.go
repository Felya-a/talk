package main

import (
	"os"
	"os/signal"
	"syscall"
	"talk/internal/app"
	"talk/internal/config"
	. "talk/internal/lib/logger"
	"talk/internal/lib/logger/logrus"
	. "talk/internal/services/auth"
	"talk/internal/utils"
)

func main() {
	config := config.MustLoad()

	db := utils.MustConnectPostgres(config)
	utils.Migrate(db)

	InitGlobalLogger(logger.NewLogrusLogger)

	authService := NewAuthService()

	application := app.New(config, authService)

	go application.WsServer.MustRun()
	go application.HttpServer.MustRun()

	Log.Info("Starting application", LogFields{"env": config.Env})

	// Graceful shutdown
	sgnl := gracefulShutdown()
	Log.Info("Stopping application", LogFields{"signal": sgnl.String()})

	// application.GrpcServer.Stop()
	application.WsServer.Stop()
	db.Close()

	Log.Info("Application stopped", nil)
}

func gracefulShutdown() os.Signal {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	sgnl := <-stop
	return sgnl
}
