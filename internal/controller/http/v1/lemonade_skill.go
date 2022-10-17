package v1

import (
	"github.com/evrone/go-clean-template/pkg/game"
	"github.com/evrone/go-clean-template/pkg/game/scene"
	"github.com/evrone/go-clean-template/pkg/marusia"
	"github.com/gin-gonic/gin"
	"time"

	"github.com/evrone/go-clean-template/pkg/logger"
)

const RequestTime = 60 * time.Second

type LemonadeSkillRoute struct {
	sdc  game.SceneDirectorConfig
	shub *game.ScriptHub
	l    logger.Interface
	wh   *marusia.Webhook
}

func newLemonadeSkillRoute(handler *gin.RouterGroup, sdc game.SceneDirectorConfig, shub *game.ScriptHub, l logger.Interface) {
	r := &LemonadeSkillRoute{
		sdc:  sdc,
		shub: shub,
		l:    l,
	}
	r.initWebhook()

	handler.POST("/lemonade", r.wh.HandleFunc)
}

type myPayload struct {
	Text string
	marusia.DefaultPayload
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

	ls.wh.OnEvent(func(r marusia.Request) (resp marusia.Response) {
		if r.Request.Command == marusia.OnStart || r.Request.Command == "debug" {
			ls.shub.RegisterClient(game.NewClient(ls.shub, r.Session.SessionID, game.NewScriptDirector(ls.sdc)))
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
					answer.WorkedClient.Close()
				}
			}
			break
		case <-ticker.C:
			ls.l.Error("Too long run scene")
			resp.Text = "Too long run scene"
			break
		}
		ticker.Stop()
		return
	})
}
