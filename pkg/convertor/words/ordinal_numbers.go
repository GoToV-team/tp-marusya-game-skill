package words

import (
	"github.com/evrone/go-clean-template/pkg/convertor/constants"
	"github.com/evrone/go-clean-template/pkg/convertor/words/declension"
	"github.com/evrone/go-clean-template/pkg/convertor/words/genders"
)

type declensionOrdinalNumbersT map[declension.Declension][constants.CountScaleNumberNameForms]string

type genderOrdinalNumbersT map[genders.Gender]declensionOrdinalNumbersT

type OrdinalNumbers struct {
	Units    []genderOrdinalNumbersT `yaml:"units,omitempty"`
	Tens     []genderOrdinalNumbersT `yaml:"tens,omitempty"`
	Dozens   []genderOrdinalNumbersT `yaml:"dozens,omitempty"`
	Hundreds []genderOrdinalNumbersT `yaml:"hundreds,omitempty"`
}
