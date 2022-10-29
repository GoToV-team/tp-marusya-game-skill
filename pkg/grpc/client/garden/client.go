package garden

import (
	"context"

	proto "github.com/evrone/go-clean-template/internal/generated/delivery/protobuf"
	"google.golang.org/grpc"
)

type GardenGameClient interface {
	Create(ctx context.Context) (string, error)
	RandomWeather(ctx context.Context, userID string) (Weather, error)
	GetBalance(ctx context.Context, userID string) (int64, error)
	Calculate(ctx context.Context, data *DayParams) (DayResult, error)
}

type GardenGame struct {
	client proto.BotanicalGardenGameClient
}

func NewBotanicalGardenGame(con *grpc.ClientConn) *GardenGame {
	client := proto.NewLemonadeGameClient(con)
	return &GardenGame{
		client: client,
	}
}

func (l *GardenGame) Create(ctx context.Context) (string, error) {
	res, err := l.client.Create(ctx, &proto.Nothing{})
	if err != nil {
		return "", err
	}
	return res.Id, err
}

func (l *GardenGame) RandomWeather(ctx context.Context, userID string) (Weather, error) {
	protoGameID := &proto.GameID{Id: userID}
	res, err := l.client.RandomWeather(ctx, protoGameID)
	if err != nil {
		return Weather{}, err
	}
	return Weather{
		Wtype:      res.WeatherName,
		RainChance: res.RainChance,
	}, err
}

func (l *GardenGame) GetBalance(ctx context.Context, userID string) (int64, error) {
	protoGameID := &proto.GameID{Id: userID}
	res, err := l.client.GetBalance(ctx, protoGameID)
	if err != nil {
		return 0, err
	}
	return res.Balance, nil
}

func (l *GardenGame) Calculate(ctx context.Context, userID string, data *DayParams) (DayResult, error) {
	protoCalculateData := &proto.CalculateRequest{
		Game:        &proto.GameID{Id: userID},
		CupsAmount:  data.CupsAmount,
		IceAmount:   data.IceAmount,
		StandAmount: data.StandAmount,
		Price:       data.Price,
	}
	res, err := l.client.Calculate(ctx, protoCalculateData)
	if err != nil {
		return DayResult{}, err
	}
	return DayResult{
		Balance: res.Balance,
		Profit:  res.Profit,
		Day:     res.Day,
	}, nil
}
