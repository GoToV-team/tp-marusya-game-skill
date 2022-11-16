package main

import (
    "github.com/ThCompiler/go_game_constractor/pkg/convertor/words"
    "github.com/ThCompiler/go_game_constractor/pkg/convertor/words/languages"
    "github.com/evrone/go-clean-template/config"
    "github.com/evrone/go-clean-template/internal/app"
    "log"
)

func main() {
    // Configuration
    cfg, err := config.NewConfig()
    if err != nil {
        log.Fatalf("Config error: %s", err)
    }

    _ = words.LoadWordsConstants(languages.Russia, "")

    // Run
    app.Run(cfg)
}
