package lemonadescript

import "github.com/evrone/go-clean-template/pkg/stringutilits"

const (
	StartTTS    = StartText
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

func GetEndOfDayTTS(balance int64, profit int64) string {
	return stringutilits.StringFormat(EndOfDayTTS,
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
