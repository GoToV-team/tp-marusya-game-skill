package lemonadescript

import "github.com/evrone/go-clean-template/pkg/stringutilits"

const (
	StartTTS = "Привет! Добро пожаловать в игру \"Лимонадная стойка\". " +
		"Правила игры очень просты. Вы решили открыть ларёк для продажи лимонада. " +
		"Каждый игровой день с утра Вы узнаёте прогноз погоды, заказываете рекламные стенды и лёд," +
		" а также решаете, сколько стаканов лимонада сделать и за сколько каждый продать. \n" +
		"На основе Ваших решений и погоды, за день Вы получаете прибыль которая отражается на вашем балансе. " +
		"Всего в игре семь дней. На этом с правилами всё.\n Один стакан стоит десять рублей.\n" +
		"Один кубик льда стоит пятьдесят рублей\nОдин стенд стоит десять рублей\n Поиграем?\n"

	GetNameTTS  = GetNameText
	HelloTTS    = HelloText
	DayInfoTTS  = DayInfoText
	IceInfoTTS  = IceInfoText
	AdjInfoTTS  = AdjInfoText
	PriceTTS    = PriceText
	EndOfDayTTS = EndOfDayText
	GoodbyeTTS  = GoodbyeText
	EndGameTTS  = EndGameText

	sunnyWeatherTTS  = sunnyWeatherText
	hotWeatherTTS    = hotWeatherText
	cloudyWeatherTTS = cloudyWeatherText
)

func GetHelloTTS(playerName string) string {
	return stringutilits.StringFormat(HelloTTS,
		"playerName", playerName,
	)
}

func GetDayInfoTTS(day uint64, balance int64, weather string, chance int64) string {
	return stringutilits.StringFormat(DayInfoTTS,
		"day", day,
		"balance", balance,
		"weather", getWeatherTTS(weather, chance),
	)
}

func GetEndOfDayTTS(glassPrice int64, icePrice int64, adjPrice int64, balance int64, profit int64) string {
	return stringutilits.StringFormat(EndOfDayText,
		"glassPrice", glassPrice,
		"icePrice", icePrice,
		"adjPrice", adjPrice,
		"profit", profit,
		"balance", balance,
	)
}

func GetEndGameTTS(balance int64) string {
	return stringutilits.StringFormat(EndGameTTS,
		"balance", balance,
	)
}

func getWeatherTTS(weather string, chance int64) string {
	switch weather {
	case cloudy:
		return stringutilits.StringFormat(cloudyWeatherTTS,
			"rainChance", chance,
		)
	case sunny:
		return sunnyWeatherTTS
	case hot:
		return hotWeatherTTS
	}
	return ""
}
