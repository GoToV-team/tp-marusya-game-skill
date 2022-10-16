package v1

import (
	"github.com/evrone/go-clean-template/pkg/game"
	"github.com/evrone/go-clean-template/pkg/game/scene"
	"github.com/evrone/go-clean-template/pkg/marusia"
	"github.com/gin-gonic/gin"

	"github.com/evrone/go-clean-template/pkg/logger"
)

type LemonadeSkillRoute struct {
	op game.Operator
	l  logger.Interface
	wh *marusia.Webhook
}

func newLemonadeSkillRoute(handler *gin.RouterGroup, op game.Operator, l logger.Interface) {
	r := &LemonadeSkillRoute{
		op: op,
		l:  l,
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
		answer := ls.op.RunScene(r.Request.Command, r.Request.OriginalUtterance, r.Request.Payload)

		resp.Text = answer.Text.BaseText
		resp.TTS = answer.Text.TextToSpeech
		resp.EndSession = answer.IsEndOfScript
		resp.Buttons = toMarusiaButtons(answer.Buttons)
		return
	})
}
