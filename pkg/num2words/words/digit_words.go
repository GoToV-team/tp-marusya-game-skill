package words

import (
	"github.com/evrone/go-clean-template/pkg/num2words/constants"
	"github.com/evrone/go-clean-template/pkg/num2words/words/declension"
	"github.com/evrone/go-clean-template/pkg/num2words/words/genders"
	"gopkg.in/yaml.v3"
)

type gendersWordT map[genders.Gender]string

type digitT struct {
	word           string
	wordWithGender gendersWordT
}

func (d *digitT) WithGender() bool {
	return d.word == ""
}

func (d *digitT) GetGendersWord() gendersWordT {
	return d.wordWithGender
}

func (d *digitT) GetWord() string {
	return d.word
}

func (d *digitT) UnmarshalYAML(n *yaml.Node) error {
	var err error
	if err = n.Decode(&d.word); err == nil {
		return nil
	}
	d.word = ""

	if err = n.Decode(&d.wordWithGender); err == nil {
		return nil
	}
	d.wordWithGender = nil

	return err
}

type DeclensionNumbers map[declension.Declension][constants.CountDigits]digitT

type DigitWords struct {
	Units    DeclensionNumbers `yaml:"units"`
	Tens     DeclensionNumbers `yaml:"tens"`
	Dozens   DeclensionNumbers `yaml:"dozens"`
	Hundreds DeclensionNumbers `yaml:"hundreds"`
}
