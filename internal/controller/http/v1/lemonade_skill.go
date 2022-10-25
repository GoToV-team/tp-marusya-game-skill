package v1

import (
	game "github.com/ThCompiler/go_game_constractor/director"
	"github.com/ThCompiler/go_game_constractor/director/scene"
	"github.com/ThCompiler/go_game_constractor/marusia"
	"github.com/gin-gonic/gin"
	"time"

	"github.com/evrone/go-clean-template/pkg/logger"
)

const RequestTime = 60 * time.Second

type LemonadeSkillRoute struct {
	sdc  game.SceneDirectorConfig
	shub marusia.ScriptRunner
	l    logger.Interface
	wh   *marusia.Webhook
}

func newLemonadeSkillRoute(handler *gin.RouterGroup, sdc game.SceneDirectorConfig,
	shub marusia.ScriptRunner, l logger.Interface) {
	r := &LemonadeSkillRoute{
		sdc:  sdc,
		shub: shub,
		l:    l,
	}
	r.initWebhook()

	handler.POST("/lemonade", r.wh.HandleFunc)
}

func toMarusiaButtons(buttons []scene.Button) []marusia.Button {
	res := make([]marusia.Button, 0)
	for _, button := range buttons {
		res = append(res, marusia.Button{
			Title:   button.Title,
			URL:     button.URL,
			Payload: button.Payload,
		})
	}
	return res
}

func (ls *LemonadeSkillRoute) initWebhook() {
	ls.wh = marusia.NewWebhook(ls.l)

	ls.wh.OnEvent(func(r marusia.Request) (resp marusia.Response, err error) {
		err = nil

		if r.Request.Command == marusia.OnStart || r.Request.Command == "debug" {
			ls.shub.AttachDirector(r.Session.SessionID, game.NewScriptDirector(ls.sdc))
		}

		ans := ls.shub.RunScene(r)

		ticker := time.NewTicker(RequestTime)
		select {
		case answer, ok := <-ans:
			if ok {
				resp.Text = answer.Text.BaseText
				resp.TTS = answer.Text.TextToSpeech
				resp.EndSession = answer.IsEndOfScript
				resp.Buttons = toMarusiaButtons(answer.Buttons)

				if answer.IsEndOfScript {
					answer.WorkedDirector.Close()
				}
			} else {
				err = BadDirectorAnswer
			}
			break
		case <-ticker.C:
			err = TooLongRunning
			break
		}
		ticker.Stop()

		if err != nil {
			ls.l.Error(err)
		}
		return
	})
}
