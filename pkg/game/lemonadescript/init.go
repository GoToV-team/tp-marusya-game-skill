package lemonadescript

import (
    game "github.com/ThCompiler/go_game_constractor/director/scriptdirector"
    client "github.com/evrone/go-clean-template/pkg/grpc/client/lemonade"
)

func NewLemonadeScript(clt client.LemonadeGameClient) game.SceneDirectorConfig {
    return game.SceneDirectorConfig{
        StartScene:   &StartScene{clt, true},
        GoodbyeScene: &InitGoodByeScene{clt},
        EndCommand:   "Пока",
    }
}
