package botanicalgardenscript

import (
	"unicode"

	"github.com/evrone/go-clean-template/pkg/stringutilits"
)

const (
	StartText = "Привет! Добро пожаловать в игру \"Ботанический сад\". Вам предстоит вырастить и продать множество красивых растений. " +
		"У Вас есть семь дней, чтобы заработать как можно больше денег! " +
		"По утрам Вы будете узнавать прогноз погоды, заказывать рекламные стенды и закупать воду для полива растений в течение дня, " +
		"а также решать, сколько горшков с растениями выставить на продажу и за сколько каждый продать. " +
		"Каждый день Вы будете зарабатывать деньги с продаж, будьте внимательнее, непродуманные вложения могут привести к убыткам! " +
		"По истечении семи игровых дней Вы увидите, сколько денег Вы заработали. Поиграем?\n Один горшок с растением стоит 10 рублей. " +
		"Один стакан воды стоит 50 рублей, такого стакана хватит, чтобы полить 3 растения.\n " +
		"Один стенд стоит 10 рублей\n Вы готовы поиграть?"

	GetNameText = "Перед началом игры мне надо узнать Ваше имя. Как Вас зовут?"

	HelloText = "{playerName}, начнём игру! Обращайте особое внимание на погоду, она играет большую роль в Вашей прибыли. " +
		"Когда погода плохая, не рассчитывайте продать столько же прохладного лимонада, сколько продали бы в солнечный день. " +
		"Устанавливайте более высокие цены в жаркие дни, чтобы увеличить прибыль, даже если вы продаете меньше лимонада."

	DayInfoText = "День {day}. У вас {balance} рублей.\n" +
		"{weather}\n" +
		"Сколько горшков с растениями Вы хотели бы продать сегодня?\n Цена одного стакана 10 рублей."

	IceInfoText = "А сколько стаканов воды заготовить?\n" +
		"Тщательный полив очень важен для растений в жаркую погоду. Цена одного стакана воды 50 рублей, такого стакана хватит для полива 3 растений"

	AdjInfoText = "Рекламные стенды делают вашу лимонадную стойку более популярной и повышают продажи.\n" +
		"Сколько Вы бы хотели поставить рекламных стендов? Один стенд стоит 10 рублей."

	PriceText = "Какую Вы назначите цену за одно растение?"

	EndOfDayText = "Рабочий день окончен!\n" +
		"За сегодня Вы потратили {glassPrice} рублей на заготовленные горшки с растениями, " +
		"{icePrice} рублей покупку стаканов с водой и {adjPrice} рублей - на рекламные стенды.\n" +
		"Ваш доход составил {profit} рублей, теперь Ваш баланс: {balance} рублей. Продолжим?"

	GoodbyeText = "Пока"

	EndGameText = "Игра окончена, Ваш баланс по истечению недели: {balance}"

	sunnyWeatherText  = "Сегодня солнечно"
	hotWeatherText    = "Сегодня очень жарко"
	cloudyWeatherText = "Сегодня облачно.  Вероятность выпадения осадков: {rainChance}.\n" +
		"Будьте осторожнее, сильный дождь может разрушить Ваши рекламные стенды."
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
