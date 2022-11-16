// Package v1 implements routing paths. Each services in own file.
package v1

import (
    http2 "github.com/ThCompiler/go_game_constractor/pkg/logger/http"
    "net/http"

    game "github.com/ThCompiler/go_game_constractor/director/scriptdirector"
    "github.com/ThCompiler/go_game_constractor/marusia/runner/hub"
    "github.com/gin-contrib/cors"
    "github.com/gin-gonic/gin"
    swaggerFiles "github.com/swaggo/files"
    ginSwagger "github.com/swaggo/gin-swagger"

    "github.com/penglongli/gin-metrics/ginmetrics"

    "github.com/ThCompiler/go_game_constractor/pkg/logger"
    // Swagger docs.
    _ "github.com/evrone/go-clean-template/docs"
)

// NewMetricRouter -.
// Swagger spec:
// @title       Metric API
// @description Metric api
// @version     1.0
// @host        localhost:8081
// @BasePath    /
func NewMetricRouter(handler *gin.Engine, appHandler *gin.Engine) {
    m := ginmetrics.GetMonitor()
    // use metric middleware without expose metric path
    m.UseWithoutExposingEndpoint(appHandler)

    m.SetMetricPath("/metrics")
    // set metric path expose to metric router
    m.Expose(handler)
}

// NewRouter -.
// Swagger spec:
// @title       Go Clean Template API
// @description Lemonade game marusia skill
// @version     1.0
// @host        localhost:8080
// @BasePath    /v1
func NewRouter(
    handler *gin.Engine,
    l logger.Interface,
    opLemonade game.SceneDirectorConfig,
    opGarden game.SceneDirectorConfig,
    opBaseGarden game.SceneDirectorConfig,
    hub *hub.ScriptHub,
) {
    // Options
    corsConfig := cors.DefaultConfig()
    corsConfig.AllowOrigins = []string{"http://localhost", "https://skill-debugger.marusia.mail.ru"}
    corsConfig.AllowMethods = []string{"POST"}

    handler.Use(cors.New(corsConfig))
    handler.Use(gin.Recovery())
    handler.Use(http2.GinRequestLogger(l))

    // Swagger
    swaggerHandler := ginSwagger.DisablingWrapHandler(swaggerFiles.Handler, "DISABLE_SWAGGER_HTTP_HANDLER")
    handler.GET("/swagger/*any", swaggerHandler)

    //// K8s probe
    handler.GET("/healthz", func(c *gin.Context) {
        c.Status(http.StatusOK)
    })

    // Routers
    h := handler.Group("/v1")
    {
        newLemonadeSkillRoute(h, opLemonade, hub, l)
        newBotanicalGardenSkillRoute(h, opGarden, hub, l)
        newBotanicalGardenBaseSkillRoute(h, opBaseGarden, hub, l)
    }
}
