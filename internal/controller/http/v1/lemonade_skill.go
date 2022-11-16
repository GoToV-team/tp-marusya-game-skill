package v1

import (
    "github.com/ThCompiler/go_game_constractor/director/scriptdirector"
    "github.com/ThCompiler/go_game_constractor/marusia/runner"
    "github.com/ThCompiler/go_game_constractor/marusia/webhook"
    "github.com/ThCompiler/go_game_constractor/pkg/logger/http"

    "github.com/ThCompiler/go_game_constractor/marusia"
    "github.com/gin-gonic/gin"

    "github.com/ThCompiler/go_game_constractor/pkg/logger"
)

type LemonadeSkillRoute struct {
    http.LogObject
    sdc  scriptdirector.SceneDirectorConfig
    shub runner.ScriptRunner
    wh   *marusia.Webhook
}

func newLemonadeSkillRoute(handler *gin.RouterGroup, sdc scriptdirector.SceneDirectorConfig,
    shub runner.ScriptRunner, l logger.Interface) {
    r := &LemonadeSkillRoute{
        LogObject: http.NewLogObject(l),
        sdc:       sdc,
        shub:      shub,
        wh:        webhook.NewDefaultMarusiaWebhook(l, shub, sdc),
    }

    handler.POST("/lemonade", r.wh.GinHandleFunc)
}

func newBotanicalGardenSkillRoute(handler *gin.RouterGroup, sdc scriptdirector.SceneDirectorConfig,
    shub runner.ScriptRunner, l logger.Interface) {
    r := &LemonadeSkillRoute{
        LogObject: http.NewLogObject(l),
        sdc:       sdc,
        shub:      shub,
        wh:        webhook.NewDefaultMarusiaWebhook(l, shub, sdc),
    }

    handler.POST("/garden", r.wh.GinHandleFunc)
}

func newBotanicalGardenBaseSkillRoute(handler *gin.RouterGroup, sdc scriptdirector.SceneDirectorConfig,
    shub runner.ScriptRunner, l logger.Interface) {
    r := &LemonadeSkillRoute{
        LogObject: http.NewLogObject(l),
        sdc:       sdc,
        shub:      shub,
        wh:        webhook.NewDefaultMarusiaWebhook(l, shub, sdc),
    }

    handler.POST("/botanic", r.wh.GinHandleFunc)
}
