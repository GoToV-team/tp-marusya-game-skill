package lemonadescript

import (
	"context"
	"github.com/ThCompiler/go_game_constractor/director/matchers"
	"github.com/ThCompiler/go_game_constractor/director/scene"
	client "github.com/evrone/go-clean-template/pkg/grpc/client/lemonade"
	"log"
	"strconv"
)

var SessionToId = make(map[string]string)

type StartScene struct {
	Game         client.LemonadeGameClient
	IsStatistics bool
}

func (ss *StartScene) React(ctx *scene.Context) scene.Command {
	switch {
	case ctx.Request.NameMatched == matchers.AgreeMatchedString:
		ss.IsStatistics = false
	default:
		ss.IsStatistics = true
	}

	if ss.IsStatistics {
		return scene.StashScene
	}
	return scene.NoCommand
}

func (ss *StartScene) Next() scene.Scene {
	if ss.IsStatistics {
		return &StatisticsScene{ss.Game}
	}
	return &GetNameScene{ss.Game}
}

func (ss *StartScene) GetSceneInfo(_ *scene.Context) (scene.Info, bool) {
	return scene.Info{
		Text: scene.Text{
			BaseText:     StartText,
			TextToSpeech: StartTTS,
		},
		ExpectedMessages: []scene.MessageMatcher{
			matchers.NewSelectorMatcher([]string{"Хочу посмотреть статистику", "статистику"},
				"Хочу посмотреть статистику"),
			matchers.Agree},
		Buttons: []scene.Button{
			{
				Title: matchers.AgreeMatchedString,
			},
			{
				Title: "Хочу посмотреть статистику",
			},
		},
		Err: &scene.BaseSceneError{Scene: &ErrorScene{ss.Game, ""}},
	}, true
}

type GetNameScene struct {
	Game client.LemonadeGameClient
}

const nameParam = "Name"

func (gns *GetNameScene) React(ctx *scene.Context) scene.Command {
	ctx.Set(nameParam, ctx.Request.SearchedMessage)
	id, err := gns.Game.Create(ctx.Context, ctx.Request.SearchedMessage)
	log.Print(err)
	SessionToId[ctx.Info.SessionId] = id
	return scene.NoCommand
}

func (gns *GetNameScene) Next() scene.Scene {
	return &HelloScene{gns.Game, ""}
}

func (gns *GetNameScene) GetSceneInfo(_ *scene.Context) (scene.Info, bool) {
	return scene.Info{
		Text: scene.Text{
			BaseText:     GetNameText,
			TextToSpeech: GetNameTTS,
		},
		ExpectedMessages: []scene.MessageMatcher{matchers.FirstWordMatcher},
	}, true
}

type HelloScene struct {
	Game   client.LemonadeGameClient
	userId string
}

func (hs *HelloScene) React(_ *scene.Context) scene.Command {
	return scene.NoCommand
}

func (hs *HelloScene) Next() scene.Scene {
	balance, _ := hs.Game.GetBalance(context.Background(), SessionToId[hs.userId])
	return &DayInfo{hs.Game, firstDay, balance, "", 0}
}

func (hs *HelloScene) GetSceneInfo(ctx *scene.Context) (scene.Info, bool) {
	hs.userId = ctx.Info.SessionId
	userName := ctx.GetString(nameParam)
	return scene.Info{
		Text: scene.Text{
			BaseText:     GetHelloText(userName),
			TextToSpeech: GetHelloTTS(userName),
		},
		ExpectedMessages: []scene.MessageMatcher{},
		Err:              &matchers.PositiveNumberError,
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
		ExpectedMessages: []scene.MessageMatcher{matchers.PositiveNumberMatcher, matchers.PositiveNumberInWordsMatcher},
		Err:              &matchers.PositiveNumberError,
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
		ExpectedMessages: []scene.MessageMatcher{matchers.PositiveNumberMatcher, matchers.PositiveNumberInWordsMatcher},
		Err:              &matchers.PositiveNumberError,
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
		ExpectedMessages: []scene.MessageMatcher{matchers.PositiveNumberMatcher, matchers.PositiveNumberInWordsMatcher},
		Err:              &matchers.PositiveNumberError,
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
		return &EndGame{pi.Game, pi.balance, 0, true}
	}
	return &EndOfDay{pi.Game, pi.balance, pi.profit, pi.day}
}

func (pi *PriceInfo) GetSceneInfo(_ *scene.Context) (scene.Info, bool) {
	return scene.Info{
		Text: scene.Text{
			BaseText:     PriceText,
			TextToSpeech: PriceTTS,
		},
		ExpectedMessages: []scene.MessageMatcher{matchers.PositiveNumberMatcher, matchers.PositiveNumberInWordsMatcher},
		Err:              &matchers.PositiveNumberError,
	}, true
}

type EndGame struct {
	Game     client.LemonadeGameClient
	balance  int64
	isEnd    int64
	showSave bool
}

func (eg *EndGame) React(ctx *scene.Context) scene.Command {
	if ctx.Request.NameMatched == "Хочу посмотреть статистику" ||
		ctx.Request.SearchedMessage == "Хочу посмотреть статистику" {
		eg.isEnd = 0
		return scene.StashScene
	}
	if ctx.Request.NameMatched == "Сохранить" ||
		ctx.Request.SearchedMessage == "Сохранить" {
		eg.isEnd = 2
		return scene.NoCommand
	}
	eg.isEnd = 1
	return scene.FinishScene
}

func (eg *EndGame) Next() scene.Scene {
	if eg.isEnd == 1 {
		return &EndGame{eg.Game, eg.balance, 0, false}
	} else if eg.isEnd == 2 {
		return &SaveStatisticsScene{Game: eg.Game, balance: eg.balance}
	} else {
		return &StatisticsScene{Game: eg.Game}
	}
}

func (eg *EndGame) GetSceneInfo(_ *scene.Context) (scene.Info, bool) {
	if eg.showSave {
		return scene.Info{
			Text: scene.Text{
				BaseText:     GetEndGameText(eg.balance),
				TextToSpeech: GetEndGameTTS(eg.balance),
			},
			ExpectedMessages: []scene.MessageMatcher{
				matchers.NewSelectorMatcher([]string{"Хочу посмотреть статистику", "статистику"},
					"Хочу посмотреть статистику"),
				matchers.NewSelectorMatcher([]string{"Хочу сохранить результат", "сохранить"},
					"Сохранить"),
				matchers.AnyMatcher,
			},
			Buttons: []scene.Button{
				{
					Title: "Хочу посмотреть статистику",
				},
				{
					Title: "Сохранить",
				},
			},
		}, true
	}
	return scene.Info{
		Text: scene.Text{
			BaseText:     GetEndGameText(eg.balance),
			TextToSpeech: GetEndGameTTS(eg.balance),
		},
		ExpectedMessages: []scene.MessageMatcher{
			matchers.NewSelectorMatcher([]string{"Хочу посмотреть статистику", "статистику"},
				"Хочу посмотреть статистику"),
			matchers.AnyMatcher,
		},
		Buttons: []scene.Button{
			{
				Title: "Хочу посмотреть статистику",
			},
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

func (eod *EndOfDay) GetSceneInfo(ctx *scene.Context) (scene.Info, bool) {
	iceN := int64(ctx.GetInt(iceNumber))
	adjN := int64(ctx.GetInt(AdjNumber))
	glassN := int64(ctx.GetInt(glassNumber))
	return scene.Info{
		Text: scene.Text{
			BaseText:     GetEndOfDayText(glassN*10, iceN*50, adjN*10, eod.balance, eod.profit),
			TextToSpeech: GetEndOfDayTTS(glassN*10, iceN*50, adjN*10, eod.balance, eod.profit),
		},
		ExpectedMessages: []scene.MessageMatcher{matchers.Agree},
		Buttons: []scene.Button{
			{
				Title: matchers.AgreeMatchedString,
			},
		},
		Err: &scene.BaseSceneError{Scene: &ErrorScene{eod.Game, ""}},
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
	return scene.Info{Text: scene.Text{BaseText: GoodbyeText, TextToSpeech: GoodbyeTTS}}, true
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

type StatisticsScene struct {
	Game client.LemonadeGameClient
}

func (ss *StatisticsScene) React(ctx *scene.Context) scene.Command {
	return scene.ApplyStashedScene
}

func (ss *StatisticsScene) Next() scene.Scene {
	return &StatisticsScene{ss.Game}
}

func (ss *StatisticsScene) GetSceneInfo(ctx *scene.Context) (scene.Info, bool) {
	stats, err := ss.Game.GetResult(ctx.Context, SessionToId[ctx.Info.SessionId])
	print(err)
	statString := ""

	for _, stat := range stats {
		statString += stat.UserName + " - " + strconv.FormatInt(stat.Result, 10) + "\n"
	}

	return scene.Info{Text: scene.Text{
		BaseText:     GetStatisticsText(statString),
		TextToSpeech: GetStatisticsTTS(statString),
	},
		ExpectedMessages: []scene.MessageMatcher{matchers.Agree},
		Buttons: []scene.Button{
			{
				Title: matchers.AgreeMatchedString,
			},
		},
		Err: &scene.BaseSceneError{Scene: &ErrorScene{ss.Game, ""}},
	}, true
}

type SaveStatisticsScene struct {
	Game    client.LemonadeGameClient
	balance int64
	cont    bool
}

func (ss *SaveStatisticsScene) React(ctx *scene.Context) scene.Command {
	ss.cont = true
	if ctx.Request.NameMatched == matchers.AgreeMatchedString {
		ss.cont = false
		_ = ss.Game.SaveResult(ctx.Context, SessionToId[ctx.Info.SessionId], ss.balance)
	}
	return scene.NoCommand
}

func (ss *SaveStatisticsScene) Next() scene.Scene {
	return &EndGame{ss.Game, ss.balance, 0, ss.cont}
}

func (ss *SaveStatisticsScene) GetSceneInfo(_ *scene.Context) (scene.Info, bool) {
	return scene.Info{Text: scene.Text{
		BaseText:     GetSaveStatisticsText(ss.balance),
		TextToSpeech: GetSaveStatisticsTTS(ss.balance),
	},
		ExpectedMessages: []scene.MessageMatcher{
			matchers.Agree,
			matchers.NewSelectorMatcher([]string{"Неа", "Нет", "Не хочу"}, "Нет"),
		},
		Buttons: []scene.Button{
			{
				Title: matchers.AgreeMatchedString,
			},
			{
				Title: "Нет",
			},
		},
		Err: &scene.BaseSceneError{Scene: &ErrorScene{ss.Game, ""}},
	}, true
}
