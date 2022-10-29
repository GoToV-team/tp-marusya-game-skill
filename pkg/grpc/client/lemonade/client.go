package lemonade

import (
	"context"

	proto "github.com/evrone/go-clean-template/internal/generated/delivery/protobuf"
	"google.golang.org/grpc"
)

type LemonadeGameClient interface {
	Create(ctx context.Context, username string) (string, error)
	RandomWeather(ctx context.Context, userID string) (Weather, error)
	GetBalance(ctx context.Context, userID string) (int64, error)
	Calculate(ctx context.Context, userID string, data *DayParams) (DayResult, error)
	SaveResult(ctx context.Context, userID string, result int64) error
	GetResult(ctx context.Context, userID string) ([]StatResult, error)
}

type LemonadeGame struct {
	client proto.LemonadeGameClient
}

func NewLemonadeGame(con *grpc.ClientConn) *LemonadeGame {
	client := proto.NewLemonadeGameClient(con)
	return &LemonadeGame{
		client: client,
	}
}

func (l *LemonadeGame) Create(ctx context.Context, username string) (string, error) {
	var res, err = l.client.Create(ctx, &proto.User{
		Username: username,
	})
	if err != nil {
		return "", err
	}
	return res.Id, err
}

func (l *LemonadeGame) RandomWeather(ctx context.Context, userID string) (Weather, error) {
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

func (l *LemonadeGame) GetBalance(ctx context.Context, userID string) (int64, error) {
	protoGameID := &proto.GameID{Id: userID}
	res, err := l.client.GetBalance(ctx, protoGameID)
	if err != nil {
		return 0, err
	}
	return res.Balance, nil
}

func (l *LemonadeGame) Calculate(ctx context.Context, userID string, data *DayParams) (DayResult, error) {
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

func (l *LemonadeGame) SaveResult(ctx context.Context, userID string, result int64) error {
	protoResultData := &proto.SaveResultMessage{
		ID: &proto.GameID{
			Id: userID,
		},
		Result: result,
	}
	_, err := l.client.SaveResult(ctx, protoResultData)

	return err
}
func (l *LemonadeGame) GetResult(ctx context.Context, userID string) ([]StatResult, error) {
	protoGameID := &proto.GameID{Id: userID}

	res, err := l.client.GetResult(ctx, protoGameID)
	if err != nil {
		return []StatResult{}, err
	}
	return exported(res.Results), nil
}

func exported(data []*proto.Result) []StatResult {
	res := make([]StatResult, 0, len(data))
	for _, val := range data {
		if val != nil {
			res = append(res, StatResult{
				Result: val.Result,
			})
		}
	}
	return res
}
