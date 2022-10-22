package words

import (
	"github.com/evrone/go-clean-template/pkg/num2words/constants"
	"github.com/evrone/go-clean-template/pkg/num2words/words/declension"
)

type declensionUnitName map[declension.Declension][constants.CountScaleNumberNameForms]string

type UnitScalesNames struct {
	Thousands      declensionUnitName `yaml:"thousands"`
	OtherEnding    declensionUnitName `yaml:"otherEnding"`
	OtherBeginning []string           `yaml:"otherBeginning"`
}
