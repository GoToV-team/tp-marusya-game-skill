package main

import (
	"github.com/evrone/go-clean-template/config"
	num2words2 "github.com/evrone/go-clean-template/pkg/num2words"
	"github.com/evrone/go-clean-template/pkg/num2words/words"
	"github.com/evrone/go-clean-template/pkg/num2words/words/languages"
	"log"
)

func main() {
	// Configuration
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	/*// Run
	app.Run(cfg)*/

	err = words.LoadWordsConstants(languages.Russia, cfg.App.ResourcesDir)
	print(err)
	res, _ := num2words2.Convert("2/2052", num2words2.DefaultOption)
	print(res)
}
