// Package app configures and runs application.
package app

import (
	"fmt"
	"github.com/evrone/go-clean-template/pkg/game/scg/botanicalgardengame/manager/usecase"
	"github.com/evrone/go-clean-template/pkg/game/scg/botanicalgardengame/script"
	redis2 "github.com/evrone/go-clean-template/pkg/game/scg/botanicalgardengame/store/redis"
	"github.com/evrone/go-clean-template/pkg/game/scg/botanicalgardengame/store/storesaver"
	"github.com/go-redis/redis/v8"
	"os"
	"os/signal"
	"syscall"

	"github.com/ThCompiler/go_game_constractor/marusia/hub"
	"github.com/evrone/go-clean-template/pkg/game/botanicalgardenscript"
	"github.com/evrone/go-clean-template/pkg/game/lemonadescript"
	grpc2 "github.com/evrone/go-clean-template/pkg/grpc"
	"github.com/evrone/go-clean-template/pkg/grpc/client/garden"
	"github.com/evrone/go-clean-template/pkg/grpc/client/lemonade"
	"github.com/gin-gonic/gin"

	"github.com/evrone/go-clean-template/config"
	v1 "github.com/evrone/go-clean-template/internal/controller/http/v1"
	"github.com/evrone/go-clean-template/pkg/httpserver"
	"github.com/evrone/go-clean-template/pkg/logger"
)

// Run creates objects via constructors.
func Run(cfg *config.Config) {
	l := logger.New(cfg.Log.Level)

	// Redis
	opt, err := redis.ParseURL(cfg.Redis.URL)
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - redis - redis.New: %w", err))
	}
	rdb := redis.NewClient(opt)

	// Repository
	botanicStore := redis2.NewScriptRepository(rdb)

	// GRPC
	grpc, err := grpc2.NewGrpcConnection(cfg.GRPC.URL)
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - grpc - grpc.New: %w", err))
	}

	err = storesaver.SaveScripts(botanicStore)
	if err != nil && err != storesaver.ScriptAlreadySaveError {
		l.Fatal(fmt.Errorf("app - Run - store - saver.SaveStore: %w", err))
	}

	// HTTP Server
	gameDirectorConfigLemonade := lemonadescript.NewLemonadeScript(lemonade.NewLemonadeGame(grpc))
	gameDirectorConfigGarden := botanicalgardenscript.NewBotanicalGardenScript(garden.NewBotanicalGardenGame(grpc))
	gameDirectorConfigBotanicGarden := script.NewBotanicalGardenGameScript(usecase.NewTextUsecase(botanicStore))

	hub := hub.NewHub()

	appHandler := gin.New()
	metricHandler := gin.New()
	v1.NewRouter(appHandler, l, gameDirectorConfigLemonade, gameDirectorConfigGarden, gameDirectorConfigBotanicGarden, hub)
	v1.NewMetricRouter(metricHandler, appHandler)
	httpServer := httpserver.New(appHandler, httpserver.Port(cfg.HTTP.Port))
	httpMetricServer := httpserver.New(appHandler, httpserver.Port(cfg.HTTP.MetricPort))

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	go hub.Run()

	select {
	case s := <-interrupt:
		l.Info("app - Run - signal: " + s.String())
	case err := <-httpServer.Notify():
		l.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	case err := <-httpMetricServer.Notify():
		l.Error(fmt.Errorf("app - Run - httpMetricServer.Notify: %w", err))
		//case err = <-rmqServer.Notify():
		//	l.Error(fmt.Errorf("app - Run - rmqServer.Notify: %w", err))
	}

	hub.StopHub()

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}
}
