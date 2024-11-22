package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"talk/internal/app"
	"talk/internal/config"
	"talk/internal/lib/logger"
	"talk/internal/utils"
)

func main() {
	config := config.MustLoad()
	logger.SetEnv(config.Env)
	log := logger.Logger()

	db := utils.MustConnectPostgres(config)
	utils.Migrate(db)

	application := app.New()
	fmt.Println(application) // TODO: DELETE LINE

	// Graceful shutdown
	sgnl := gracefulShutdown()
	log.Info("Stopping application", slog.String("signal", sgnl.String()))

	// application.GrpcServer.Stop()
	// application.HttpServer.Stop()
	db.Close()

	log.Info("Application stopped")
}

func gracefulShutdown() os.Signal {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	sgnl := <-stop
	return sgnl
}