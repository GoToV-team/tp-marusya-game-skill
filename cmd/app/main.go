package main

import (
	"github.com/evrone/go-clean-template/config"
	"github.com/evrone/go-clean-template/internal/app"
	"github.com/evrone/go-clean-template/pkg/convertor/words"
	"github.com/evrone/go-clean-template/pkg/convertor/words/languages"
	"log"
)

func main() {
	// Configuration
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	_ = words.LoadWordsConstants(languages.Russia, cfg.ResourcesDir)

	// Run
	app.Run(cfg)
}
