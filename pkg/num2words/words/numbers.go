package words

import (
	"fmt"
	"github.com/evrone/go-clean-template/pkg/num2words/currency"
	"github.com/evrone/go-clean-template/pkg/num2words/words/languages"
	"github.com/ilyakaznacheev/cleanenv"
)

var WordConstantsForNumbers wordsConstantsForNumbers

type CurrencyWords struct {
	Currencies map[currency.Currency]currency.CustomCurrency `yaml:"currencies"`
}

type Sign struct {
	Minus string `yaml:"minus"`
}

type wordsConstantsForNumbers struct {
	UnitScalesNames         UnitScalesNames         `yaml:"unitScalesNames"`
	SlashNumberUnitPrefixes SlashNumberUnitPrefixes `yaml:"slashNumberUnitPrefixes"`
	DigitWords              DigitWords              `yaml:"digitWords"`
	CurrenciesStrings       CurrencyWords           `yaml:"currenciesStrings"`
	FractionalUnit          FractionalUnit          `yaml:"fractionalUnit"`
	OrdinalNumbers          OrdinalNumbers          `yaml:"ordinalNumbers"`
	Sign                    Sign                    `yaml:"sign"`
}

func LoadWordsConstants(lang languages.Language, resourcesDirPath string) error {
	WordConstantsForNumbers = wordsConstantsForNumbers{}

	// Sign
	err := cleanenv.ReadConfig(resourcesDirPath+"/"+string(lang)+"/sign.yml", &WordConstantsForNumbers.Sign)
	if err != nil {
		return fmt.Errorf("error load %s sign: %w", lang, err)
	}

	// Digit words
	err = cleanenv.ReadConfig(resourcesDirPath+"/"+string(lang)+"/digit_words.yml",
		&WordConstantsForNumbers.DigitWords)
	if err != nil {
		return fmt.Errorf("error load %s digit words: %w", lang, err)
	}

	// Fractional unit
	err = cleanenv.ReadConfig(resourcesDirPath+"/"+string(lang)+"/fractional_unit.yml",
		&WordConstantsForNumbers.FractionalUnit)
	if err != nil {
		return fmt.Errorf("error load %s fractional unit: %w", lang, err)
	}

	// Unit scales names
	err = cleanenv.ReadConfig(resourcesDirPath+"/"+string(lang)+"/unit_scales_names.yml",
		&WordConstantsForNumbers.UnitScalesNames)
	if err != nil {
		return fmt.Errorf("error load %s unit scales names: %w", lang, err)
	}

	// Slash number unit prefixes
	err = cleanenv.ReadConfig(resourcesDirPath+"/"+string(lang)+"/slash_number_unit_prefixes.yml",
		&WordConstantsForNumbers.SlashNumberUnitPrefixes)
	if err != nil {
		return fmt.Errorf("error load %s slash number unit prefixes: %w", lang, err)
	}

	// Currencies strings
	err = cleanenv.ReadConfig(resourcesDirPath+"/"+string(lang)+"/currencies_strings.yml",
		&WordConstantsForNumbers.CurrenciesStrings)
	if err != nil {
		return fmt.Errorf("error load %s currencies strings: %w", lang, err)
	}

	// Ordinal numbers
	err = cleanenv.ReadConfig(resourcesDirPath+"/"+string(lang)+"/ordinal_numbers.yml",
		&WordConstantsForNumbers.OrdinalNumbers)
	if err != nil {
		return fmt.Errorf("error load %s ordinal numbers: %w", lang, err)
	}

	return nil
}
