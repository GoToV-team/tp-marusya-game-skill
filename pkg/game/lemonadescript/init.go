package lemonadescript

import (
	"github.com/evrone/go-clean-template/pkg/game"
	"github.com/evrone/go-clean-template/pkg/grpc/client"
)

func NewLemonadeScript(client client.LemonadeGameClient) game.SceneDirectorConfig {
	return game.SceneDirectorConfig{
		StartScene:   &StartScene{client},
		GoodbyeScene: &InitGoodByeScene{client},
		EndCommand:   "Пока",
	}
}
