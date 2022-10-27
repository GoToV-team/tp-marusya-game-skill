package botanicalgardenscript

import (
	"strconv"

	"github.com/evrone/go-clean-template/pkg/convertor/words2num"
)

type NumberMatchers struct{}

func (nm NumberMatchers) Match(message string) (bool, string) {
	res, _ := words2num.Convert(message)
	return res != 0, strconv.FormatInt(res, 10)
}

var PositiveNumber = NumberMatchers{}
