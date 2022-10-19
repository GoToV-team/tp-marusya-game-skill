package lemonadescript

import (
	"github.com/evrone/go-clean-template/pkg/stringutilits"
	"unicode"
)

const (
	StartText = "Привет! Добро пожаловать в игру \"Лимонадная стойка\". " +
		"Правила игры очень просты. Вы решили открыть ларёк для продажи лимонада. " +
		"Каждый игровой день с утра Вы узнаёте прогноз погоды, заказываете рекламные стенды и лёд," +
		" а также решаете, сколько стаканов лимонада сделать и за сколько каждый продать. \n" +
		"На основе Ваших решений и погоды, за день Вы получаете прибыль которая отражается на вашем балансе. " +
		"Всего в игре семь дней. На этом с правилами всё.\nПоиграем?\n"

	GetNameText = "Перед началом игры мне надо узнать Ваше имя. Как Вас зовут?"

	HelloText = "{playerName}, начнём игру!"

	DayInfoText = "День {day}. У вас {balance} рублей.\n" +
		"{weather}\n" +
		"Сколько вы бы хотели заготовить стаканом лимонада?"

	IceInfoText = "А сколько кубиков льда купить?"

	AdjInfoText = "Сколько Вы хотите поставить рекламных стендов?"

	PriceText = "Какую Вы назначите цену за один стакан лимонада?"

	EndOfDayText = "Рабочий день окончен!\n" +
		" За сегодняшний день Вы заработали: {profit} рублей, Ваш баланс: {balance} рублей. \n Продолжим?"

	GoodbyeText = "Пока"

	EndGameText = "Игра окончена, Ваш баланс по истечению недели: {balance}"

	sunnyWeatherText  = "Сегодня солнечно"
	hotWeatherText    = "Сегодня очень жарко"
	cloudyWeatherText = "Сегодня облачно.  Вероятность выпадения осадков: {rainChance}. Будьте внимательнее"
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

func GetEndOfDayText(balance int64, profit int64) string {
	return stringutilits.StringFormat(EndOfDayText,
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
