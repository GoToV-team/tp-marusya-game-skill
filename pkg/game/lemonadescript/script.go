package lemonadescript

import (
	"github.com/evrone/go-clean-template/pkg/game/matchers"
	"github.com/evrone/go-clean-template/pkg/game/scene"
	"github.com/evrone/go-clean-template/pkg/grpc/client"
	"log"
	"strconv"
)

var SessionToId = make(map[string]string)

type StartScene struct {
	Game client.LemonadeGameClient
}

func (ss *StartScene) React(ctx *scene.Context) scene.Command {
	id, err := ss.Game.Create(ctx.Context)
	log.Print(err)
	SessionToId[ctx.Info.SessionId] = id
	return scene.NoCommand
}

func (ss *StartScene) Next() scene.Scene {
	return &GetNameScene{ss.Game}
}

func (ss *StartScene) GetSceneInfo(_ *scene.Context) (scene.Info, bool) {
	return scene.Info{
		Text: scene.Text{
			BaseText:     StartText,
			TextToSpeech: StartTTS,
		},
		ExpectedMessages: []scene.MessageMatcher{matchers.Agree},
		Buttons: []scene.Button{
			{
				Title: matchers.AgreeString,
			},
		},
	}, true
}

type GetNameScene struct {
	Game client.LemonadeGameClient
}

const nameParam = "Name"

func (gns *GetNameScene) React(ctx *scene.Context) scene.Command {
	ctx.Set(nameParam, ctx.Request.SearchedMessage)
	return scene.NoCommand
}

func (gns *GetNameScene) Next() scene.Scene {
	return &HelloScene{gns.Game}
}

func (gns *GetNameScene) GetSceneInfo(_ *scene.Context) (scene.Info, bool) {
	return scene.Info{
		Text: scene.Text{
			BaseText:     GetNameText,
			TextToSpeech: GetNameTTS,
		},
		ExpectedMessages: []scene.MessageMatcher{matchers.FirstWord},
	}, true
}

type HelloScene struct {
	Game client.LemonadeGameClient
}

func (hs *HelloScene) React(_ *scene.Context) scene.Command {
	return scene.NoCommand
}

func (hs *HelloScene) Next() scene.Scene {
	return &DayInfo{hs.Game, firstDay, startBalance, "", 0}
}

func (hs *HelloScene) GetSceneInfo(ctx *scene.Context) (scene.Info, bool) {
	userName := ctx.GetString(nameParam)
	return scene.Info{
		Text: scene.Text{
			BaseText:     GetHelloText(userName),
			TextToSpeech: GetHelloTTS(userName),
		},
		ExpectedMessages: []scene.MessageMatcher{},
	}, false
}

type DayInfo struct {
	Game    client.LemonadeGameClient
	Day     uint64
	Balance int64
	Weather string
	Chance  int64
}

const glassNumber = "glassNumber"

func (gns *DayInfo) React(ctx *scene.Context) scene.Command {
	number, _ := strconv.Atoi(ctx.Request.SearchedMessage)
	ctx.Set(glassNumber, number)
	return scene.NoCommand
}

func (gns *DayInfo) Next() scene.Scene {
	return &IceInfo{gns.Game}
}

func (gns *DayInfo) GetSceneInfo(ctx *scene.Context) (scene.Info, bool) {
	weather, err := gns.Game.RandomWeather(ctx.Context, SessionToId[ctx.Info.SessionId])
	log.Print(err)
	gns.Weather = weather.Wtype
	gns.Chance = weather.RainChance

	return scene.Info{
		Text: scene.Text{
			BaseText:     GetDayInfoText(gns.Day, gns.Balance, gns.Weather, gns.Chance),
			TextToSpeech: GetDayInfoTTS(gns.Day, gns.Balance, gns.Weather, gns.Chance),
		},
		ExpectedMessages: []scene.MessageMatcher{matchers.PositiveNumberMatcher},
	}, true
}

type IceInfo struct {
	Game client.LemonadeGameClient
}

const iceNumber = "iceNumber"

func (ii *IceInfo) React(ctx *scene.Context) scene.Command {
	number, _ := strconv.Atoi(ctx.Request.SearchedMessage)
	ctx.Set(iceNumber, number)
	return scene.NoCommand
}

func (ii *IceInfo) Next() scene.Scene {
	return &AdjInfo{ii.Game}
}

func (ii *IceInfo) GetSceneInfo(_ *scene.Context) (scene.Info, bool) {
	return scene.Info{
		Text: scene.Text{
			BaseText:     IceInfoText,
			TextToSpeech: IceInfoTTS,
		},
		ExpectedMessages: []scene.MessageMatcher{matchers.PositiveNumberMatcher},
	}, true
}

type AdjInfo struct {
	Game client.LemonadeGameClient
}

const AdjNumber = "adjNumber"

func (ai *AdjInfo) React(ctx *scene.Context) scene.Command {
	number, _ := strconv.Atoi(ctx.Request.SearchedMessage)
	ctx.Set(AdjNumber, number)
	return scene.NoCommand
}

func (ai *AdjInfo) Next() scene.Scene {
	return &PriceInfo{ai.Game, 0, 0, 0}
}

func (ai *AdjInfo) GetSceneInfo(_ *scene.Context) (scene.Info, bool) {
	return scene.Info{
		Text: scene.Text{
			BaseText:     AdjInfoText,
			TextToSpeech: AdjInfoTTS,
		},
		ExpectedMessages: []scene.MessageMatcher{matchers.PositiveNumberMatcher},
	}, true
}

type PriceInfo struct {
	Game    client.LemonadeGameClient
	day     uint64
	profit  int64
	balance int64
}

const Price = "price"

func (pi *PriceInfo) React(ctx *scene.Context) scene.Command {
	number, _ := strconv.Atoi(ctx.Request.SearchedMessage)
	ctx.Set(Price, number)
	iceN := ctx.GetInt(iceNumber)
	adjN := ctx.GetInt(AdjNumber)
	glassN := ctx.GetInt(glassNumber)

	res, err := pi.Game.Calculate(ctx.Context, SessionToId[ctx.Info.SessionId], &client.DayParams{
		CupsAmount:  int64(glassN),
		IceAmount:   int64(iceN),
		StandAmount: int64(adjN),
		Price:       int64(number),
	})

	log.Print(err)

	pi.day = uint64(res.Day)
	pi.balance = res.Balance
	pi.profit = res.Profit

	return scene.NoCommand
}

func (pi *PriceInfo) Next() scene.Scene {
	if pi.day == maxDay {
		return &EndGame{pi.Game, pi.balance}
	}
	return &EndOfDay{pi.Game, pi.balance, pi.profit, pi.day}
}

func (pi *PriceInfo) GetSceneInfo(_ *scene.Context) (scene.Info, bool) {
	return scene.Info{
		Text: scene.Text{
			BaseText:     PriceText,
			TextToSpeech: PriceTTS,
		},
		ExpectedMessages: []scene.MessageMatcher{matchers.PositiveNumberMatcher},
	}, true
}

type EndGame struct {
	Game    client.LemonadeGameClient
	balance int64
}

func (eg *EndGame) React(_ *scene.Context) scene.Command {
	return scene.FinishScene
}

func (eg *EndGame) Next() scene.Scene {
	return &EndGame{eg.Game, eg.balance}
}

func (eg *EndGame) GetSceneInfo(_ *scene.Context) (scene.Info, bool) {
	return scene.Info{
		Text: scene.Text{
			BaseText:     GetEndGameText(eg.balance),
			TextToSpeech: GetEndGameTTS(eg.balance),
		},
	}, true
}

type EndOfDay struct {
	Game    client.LemonadeGameClient
	balance int64
	profit  int64
	day     uint64
}

func (eod *EndOfDay) React(_ *scene.Context) scene.Command {
	return scene.NoCommand
}

func (eod *EndOfDay) Next() scene.Scene {
	return &DayInfo{eod.Game, eod.day + 1, eod.balance, "", 0}
}

func (eod *EndOfDay) GetSceneInfo(_ *scene.Context) (scene.Info, bool) {
	return scene.Info{
		Text: scene.Text{
			BaseText:     GetEndOfDayText(eod.balance, eod.profit),
			TextToSpeech: GetEndOfDayTTS(eod.balance, eod.profit),
		},
		ExpectedMessages: []scene.MessageMatcher{matchers.Agree},
		Buttons: []scene.Button{
			{
				Title: matchers.AgreeString,
			},
		},
	}, true
}

type InitGoodByeScene struct {
	Game client.LemonadeGameClient
}

func (igs *InitGoodByeScene) React(_ *scene.Context) scene.Command {
	return scene.FinishScene
}

func (igs *InitGoodByeScene) Next() scene.Scene {
	return &InitGoodByeScene{igs.Game}
}

func (igs *InitGoodByeScene) GetSceneInfo(_ *scene.Context) (scene.Info, bool) {
	return scene.Info{}, true
}

type ErrorScene struct {
	Game     client.LemonadeGameClient
	userText string
}

func (es *ErrorScene) React(ctx *scene.Context) scene.Command {
	es.userText = ctx.Request.SearchedMessage
	return scene.ApplyStashedScene
}

func (es *ErrorScene) Next() scene.Scene {
	return &ErrorScene{es.Game, es.userText}
}

func (es *ErrorScene) GetSceneInfo(_ *scene.Context) (scene.Info, bool) {
	return scene.Info{Text: scene.Text{
		BaseText:     "Я не знаю такую команду " + es.userText,
		TextToSpeech: "Я не знаю такую команду" + es.userText,
	}}, true
}
