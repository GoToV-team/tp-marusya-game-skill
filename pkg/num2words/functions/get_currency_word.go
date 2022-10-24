package functions

import (
	"github.com/evrone/go-clean-template/pkg/num2words/constants"
	"github.com/evrone/go-clean-template/pkg/num2words/currency"
	"github.com/evrone/go-clean-template/pkg/num2words/objects"
	"github.com/evrone/go-clean-template/pkg/num2words/words/declension"
)

func GetCurrencyWord(
	currencyObject currency.CustomCurrency, numberPart constants.NumberType,
	scaleNameForm objects.ScaleForm, lastScaleIsZero bool,
	curc currency.Currency, decl declension.Declension,
) string {
	declensionsObject := currencyObject.FractionalPartNameDeclensions
	if numberPart == constants.DecimalNumber {
		declensionsObject = currencyObject.CurrencyNameDeclensions
	}

	scaleForm := 1
	if scaleNameForm == 0 {
		scaleForm = 0
	}
	currentDeclension := decl

	// Если падеж "именительный" или "винительный" и множественное число
	if (decl == declension.NOMINATIVE || decl == declension.ACCUSATIVE) && scaleNameForm >= 1 {
		// Если валюта указана как "number"
		if curc == currency.NUMBER {
			scaleForm = 1
		} else {
			scaleForm = 1
			if scaleNameForm == 1 {
				scaleForm = 0
			}
		}
		// Использовать родительный падеж.
		currentDeclension = declension.GENITIVE
	}
	// Если последний класс числа равен "000"
	if lastScaleIsZero {
		scaleForm = 1
		// Всегда родительный падеж и множественное число
		currentDeclension = declension.GENITIVE

	}

	return declensionsObject[currentDeclension][scaleForm]
}
