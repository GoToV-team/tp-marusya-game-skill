package words

import (
	"github.com/evrone/go-clean-template/pkg/num2words/constants"
)

type SlashNumberUnitPrefixes struct {
	Units    [constants.CountDigits]string `yaml:"units"`
	Tens     [constants.CountDigits]string `yaml:"tens"`
	Dozens   [constants.CountDigits]string `yaml:"dozens"`
	Hundreds [constants.CountDigits]string `yaml:"hundreds"`
}
