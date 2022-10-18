package client

type Weather struct {
	Wtype      string
	RainChance int64
}

type DayParams struct {
	CupsAmount  int64
	IceAmount   int64
	StandAmount int64
	Price       int64
}

type DayResult struct {
	Balance int64
	Profit  int64
	Day     int64
}
