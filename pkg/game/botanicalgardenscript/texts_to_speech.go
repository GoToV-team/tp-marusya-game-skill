package botanicalgardenscript

import "github.com/evrone/go-clean-template/pkg/stringutilits"

const (
	StartTTS = "Привет! Добро пожаловать в игру \"Ботанический сад\". " +
		"Вам предстоит вырастить и продать множество красивых растений. " +
		"У Вас есть семь дней, чтобы заработать как можно больше денег! " +
		"По утрам Вы будете узнавать прогноз погоды, заказывать рекламные стенды и закупать воду для полива растений в течение дня, " +
		"а также решать, сколько горшков с растениями выставить на продажу и за сколько каждый продать. " +
		"Каждый день Вы будете зарабатывать деньги с продаж, будьте внимательнее, непродуманные вложения могут привести к убыткам! " +
		"По истечении семи игровых дней Вы увидите, сколько денег Вы заработали. Поиграем?\n Один горшок с растением стоит 10 рублей. " +
		"Один стакан воды стоит 50 рублей, такого стакана хватит, чтобы полить 3 растения.\n " +
		"Один стенд стоит 10 рублей\n Вы готовы поиграть?"

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
