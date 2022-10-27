package botanicalgardenscript

import (
	game "github.com/ThCompiler/go_game_constractor/director"
	"github.com/evrone/go-clean-template/pkg/grpc/client"
)

func NewBotanicalGardenScript(client client.LemonadeGameClient) game.SceneDirectorConfig {
	return game.SceneDirectorConfig{
		StartScene:   &StartScene{client},
		GoodbyeScene: &InitGoodByeScene{client},
		EndCommand:   "Пока",
	}
}
