package lemonadescript

import (
	"github.com/evrone/go-clean-template/pkg/stringutilits"
	"unicode"
)

const (
	StartText = "Привет! Добро пожаловать в игру \"Лимонадная стойка\". Правила игры очень просты." +
		" Вы решили открыть ларёк для продажи лимонада. У Вас есть семь дней, чтобы заработать как можно больше денег!" +
		" Каждый игровой день с утра Вы узнаёте прогноз погоды, заказываете рекламные стенды и лёд, а также решаете, " +
		"сколько стаканов лимонада сделать и за сколько каждый продать. На основе Ваших решений и погоды, за день " +
		"Вы получаете прибыль которая отражается на вашем балансе. По истечении семи игровых дней " +
		"Вы увидите, сколько денег Вы заработали. Поиграем?\n Один стакан стоит 10 рублей. " +
		"Один кубик льда стоит 50 рублей\nОдин стенд стоит 10 рублей\n Поиграем?"

	GetNameText = "Перед началом игры мне надо узнать Ваше имя. Как Вас зовут?"

	HelloText = "{playerName}, начнём игру! Обращайте особое внимание на погоду, она играет большую роль в Вашей прибыли. " +
		"Когда погода плохая, не рассчитывайте продать столько же прохладного лимонада, сколько продали бы в солнечный день. " +
		"Устанавливайте более высокие цены в жаркие дни, чтобы увеличить прибыль, даже если вы продаете меньше лимонада."

	DayInfoText = "День {day}. У вас {balance} рублей.\n" +
		"{weather}\n" +
		"Сколько Вы бы хотели заготовить бумажных стаканчиков для лимонада?\n Цена одного стакана 10 рублей."

	IceInfoText = "А сколько кубиков льда купить?\n" +
		" Прохладный лимонад особенно актуален в жаркую погоду. Цена одного кубика льда 50 рублей."

	AdjInfoText = "Рекламные стенды делают вашу лимонадную стойку более популярной и повышают продажи.\n" +
		"Сколько Вы бы хотели поставить рекламных стендов? Один стенд стоит 10 рублей."

	PriceText = "Какую Вы назначите цену за один стакан лимонада?"

	EndOfDayText = "Рабочий день окончен!\n" +
		"За сегодня Вы потратили {glassPrice} рублей на заготовленные бумажные стаканчики, " +
		"{icePrice} рублей на покупку кубиков льда и {adjPrice} рублей - на рекламные стенды и заработали {profit} рублей," +
		" теперь Ваш баланс: {balance} рублей. Продолжим?"

	GoodbyeText = "Пока"

	EndGameText = "Игра окончена, Ваш баланс по истечению недели: {balance}"

	sunnyWeatherText  = "Сегодня солнечно"
	hotWeatherText    = "Сегодня очень жарко"
	cloudyWeatherText = "Сегодня облачно.  Вероятность выпадения осадков: {rainChance}. \n " +
		"Будьте внимательнее, если начнется сильный дождь, Ваши рекламные стенды могут быть разрушены"
)

func GetHelloText(playerName string) string {
	r := []rune(playerName)
	r[0] = unicode.ToUpper(r[0])
	s := string(r)
	return stringutilits.StringFormat(HelloText,
		"playerName", s,
	)
}

func GetDayInfoText(day uint64, balance int64, weather string, chance int64) string {
	return stringutilits.StringFormat(DayInfoText,
		"day", day,
		"balance", balance,
		"weather", getWeatherText(weather, chance),
	)
}

func GetEndOfDayText(glassPrice int64, icePrice int64, adjPrice int64, balance int64, profit int64) string {
	return stringutilits.StringFormat(EndOfDayText,
		"glassPrice", glassPrice,
		"icePrice", icePrice,
		"adjPrice", adjPrice,
		"profit", profit,
		"balance", balance,
	)
}

func GetEndGameText(balance int64) string {
	return stringutilits.StringFormat(EndGameText,
		"balance", balance,
	)
}

func getWeatherText(weather string, chance int64) string {
	switch weather {
	case cloudy:
		return stringutilits.StringFormat(cloudyWeatherText,
			"rainChance", chance,
		)
	case sunny:
		return sunnyWeatherText
	case hot:
		return hotWeatherText
	}
	return ""
}
