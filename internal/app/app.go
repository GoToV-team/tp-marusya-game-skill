// Package app configures and runs application.
package app

import (
	"fmt"
	"github.com/evrone/go-clean-template/pkg/game"
	"github.com/evrone/go-clean-template/pkg/game/scene"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"os"
	"os/signal"
	"syscall"

	"github.com/evrone/go-clean-template/config"
	v1 "github.com/evrone/go-clean-template/internal/controller/http/v1"
	"github.com/evrone/go-clean-template/pkg/httpserver"
	"github.com/evrone/go-clean-template/pkg/logger"
)

type GameScene struct {
	userText string
}

func (gs *GameScene) React(ctx *scene.Context) (scene.Scene, scene.Command) {
	return &GameScene{
		ctx.Request.SearchedMessage,
	}, scene.NoCommand
}

func (gs *GameScene) GetSceneInfo() scene.Info {
	return scene.Info{
		Text: scene.Text{
			BaseText:     "Ты сказал: " + gs.userText,
			TextToSpeech: gs.userText,
		},
		ExpectedMessages: []string{"*"},
	}
}

type StartScene struct{}

func (ss *StartScene) React(ctx *scene.Context) (scene.Scene, scene.Command) {
	return &GameScene{ctx.Request.SearchedMessage}, scene.NoCommand
}

func (ss *StartScene) GetSceneInfo() scene.Info {
	return scene.Info{
		Text: scene.Text{
			BaseText:     "Привет путник. Я буду повторять слова за тобой",
			TextToSpeech: "Привет п+утник. Я б+уду повтор+ять сл+ова за тобой",
		},
		ExpectedMessages: []string{"*"},
	}
}

type InitGoodByeScene struct{}

func (igs *InitGoodByeScene) React(_ *scene.Context) (scene.Scene, scene.Command) {
	return &GoodByeScene{}, scene.NoCommand
}

func (igs *InitGoodByeScene) GetSceneInfo() scene.Info {
	return scene.Info{Text: scene.Text{
		BaseText:     "Привет путник. Я буду повторять слова за тобой",
		TextToSpeech: "Привет п+утник. Я б+уду повтор+ять сл+ова за тобой",
	}}
}

type ErrorScene struct {
	userText string
}

func (es *ErrorScene) React(ctx *scene.Context) (scene.Scene, scene.Command) {
	return &ErrorScene{ctx.Request.SearchedMessage}, scene.ApplyStashedScene
}

func (es *ErrorScene) GetSceneInfo() scene.Info {
	return scene.Info{Text: scene.Text{
		BaseText:     "Я не знаю такую команду " + es.userText,
		TextToSpeech: "Я не знаю такую команду" + es.userText,
	}}
}

type GoodByeScene struct{}

func (gs *GoodByeScene) React(ctx *scene.Context) (scene.Scene, scene.Command) {
	if ctx.Request.SearchedMessage == "Точно" {
		return nil, scene.FinishScene
	}
	return nil, scene.ApplyStashedScene
}

func (gs *GoodByeScene) GetSceneInfo() scene.Info {
	return scene.Info{Text: scene.Text{
		BaseText: "Ты точно решил нас покинуть?",
	},
		Buttons: []scene.Button{
			{
				Title: "Точно",
			},
			{
				Title: "Я передумал",
			},
		},
		ExpectedMessages: []string{
			"Точно", "Я передумал",
		},
	}
}

// Run creates objects via constructors.
func Run(cfg *config.Config) {
	l := logger.New(cfg.Log.Level)

	//// Repository
	//pg, err := postgres.New(cfg.PG.URL, postgres.MaxPoolSize(cfg.PG.PoolMax))
	//if err != nil {
	//	l.Fatal(fmt.Errorf("app - Run - postgres.New: %w", err))
	//}
	//defer pg.Close()
	//
	//// Use case
	//translationUseCase := usecase.New(
	//	repo.New(pg),
	//	webapi.New(),
	//)
	//
	//// RabbitMQ RPC Server
	//rmqRouter := amqprpc.NewRouter(translationUseCase)
	//
	//rmqServer, err := server.New(cfg.RMQ.URL, cfg.RMQ.ServerExchange, rmqRouter, l)
	//if err != nil {
	//	l.Fatal(fmt.Errorf("app - Run - rmqServer - server.New: %w", err))
	//}

	// HTTP Server
	handler := gin.New()
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost", "https://skill-debugger.marusia.mail.ru"}
	corsConfig.AllowMethods = []string{"POST"}

	gameDirectorConfig := game.SceneDirectorConfig{
		EndCommand:   "Пока",
		StartScene:   &StartScene{},
		ErrorScene:   &ErrorScene{},
		GoodbyeScene: &InitGoodByeScene{},
		GoodbyeMessage: scene.Text{
			"Пока!",
			"Пока!",
		},
	}

	hub := game.NewHub()

	handler.Use(cors.New(corsConfig))
	v1.NewRouter(handler, l, gameDirectorConfig, hub)
	httpServer := httpserver.New(handler, httpserver.Port(cfg.HTTP.Port))

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	go hub.Run()

	select {
	case s := <-interrupt:
		l.Info("app - Run - signal: " + s.String())
	case err := <-httpServer.Notify():
		l.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
		//case err = <-rmqServer.Notify():
		//	l.Error(fmt.Errorf("app - Run - rmqServer.Notify: %w", err))
	}

	hub.StopHub()

	// Shutdown
	err := httpServer.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}

	//err = rmqServer.Shutdown()
	//if err != nil {
	//	l.Error(fmt.Errorf("app - Run - rmqServer.Shutdown: %w", err))
	//}
}
