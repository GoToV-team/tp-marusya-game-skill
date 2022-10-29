package botanicalgardenscript

import (
	game "github.com/ThCompiler/go_game_constractor/director"
	client "github.com/evrone/go-clean-template/pkg/grpc/client/garden"
)

func NewBotanicalGardenScript(client client.GardenGameClient) game.SceneDirectorConfig {
	return game.SceneDirectorConfig{
		StartScene:   &StartScene{client},
		GoodbyeScene: &InitGoodByeScene{client},
		EndCommand:   "Пока",
	}
}
