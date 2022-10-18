package stringutilits

// TODO
//import (
//	"regexp"
//	"strconv"
//	"strings"
//)
//
//type Options struct{}
//
//type ConvertedNumber interface {
//	int64
//	string
//}
//
//// Convert converts number into the words representation.
//func Convert(number interface{}, options Options) ([]string, error) {
//	numberString := convertNumberToString(number)
//	if numberString == "" {
//		return []string{}, nil
//	}
//	// Обработать введенное число
//
//	numberArray := splitNumberToArray(numberString)
//	/*// Собрать конечный словесный результат
//	const convertedNumberString = combineResultData(numberArray, appliedOptions)
//	return convertedNumberString, nil*/
//	return numberArray, nil
//}
//
//func convertNumberToString(number interface{}) string {
//	switch number.(type) {
//	case string:
//		return number.(string)
//	case int:
//		return strconv.FormatInt(int64(number.(int)), 10)
//	case int8:
//		return strconv.FormatInt(int64(number.(int8)), 10)
//	case int16:
//		return strconv.FormatInt(int64(number.(int16)), 10)
//	case int32:
//		return strconv.FormatInt(int64(number.(int32)), 10)
//	case int64:
//		return strconv.FormatInt(number.(int64), 10)
//	case uint:
//		return strconv.FormatUint(uint64(number.(uint)), 10)
//	case uint8:
//		return strconv.FormatUint(uint64(number.(uint8)), 10)
//	case uint16:
//		return strconv.FormatUint(uint64(number.(uint16)), 10)
//	case uint32:
//		return strconv.FormatUint(uint64(number.(uint32)), 10)
//	case uint64:
//		return strconv.FormatUint(number.(uint64), 10)
//	case float32:
//		return strconv.FormatFloat(float64(number.(float32)), 'f', -1, 32)
//	case float64:
//		return strconv.FormatFloat(number.(float64), 'f', -1, 64)
//	default:
//		return ""
//	}
//}
//
//const (
//	maxIntegerPartLength = 306
//
//	countHumberPart = 4
//	numberSign      = 0
//	leftPartNumber  = 1
//	dividerNumber   = 2
//	rightPartNumber = 3
//
//	decimalNumber    = ","
//	fractionalNumber = "/"
//)
//
//func splitNumberToArray(number string) []string {
//	numberArray := make([]string, 4)
//	numberArray[dividerNumber] = decimalNumber
//	numberArray[numberSign] = "+"
//
//	// Убрать из строки всё лишнее
//	cleanNumber := clearFromString(number, `[^\d\.\,\/\-]`)
//	if len(cleanNumber) < 1 {
//		cleanNumber = "0"
//	}
//
//	if strings.ContainsRune(cleanNumber, '-') {
//		numberArray[numberSign] = "-"
//	}
//	// Удалить все знаки минуса
//	cleanNumber = clearFromString(cleanNumber, `[\-]`)
//
//	leftPart := ""
//	rightPart := ""
//	found := false
//	// Добавить разделитель числа в массив и разделить число
//	if leftPart, rightPart, found = strings.Cut(cleanNumber, ","); found {
//		numberArray[dividerNumber] = decimalNumber
//	} else if leftPart, rightPart, found = strings.Cut(cleanNumber, "."); found {
//		numberArray[dividerNumber] = decimalNumber
//	} else if leftPart, rightPart, found = strings.Cut(cleanNumber, "/"); found {
//		numberArray[dividerNumber] = fractionalNumber
//	}
//
//	// Удалить все разделители числа
//	leftPart = clearFromString(leftPart, `[\, \.\/]`)
//	rightPart = clearFromString(rightPart, `[\, \.\/]`)
//
//	// Убрать лишние нули из целой части
//	numberArray[leftPartNumber] = clearFromString(leftPart, `^0+/`)
//
//	numberArray[rightPartNumber] = rightPart
//	// Убрать лишние нули из дробной части
//	if numberArray[dividerNumber] == decimalNumber {
//		numberArray[rightPartNumber] = clearFromString(rightPart, `^0+/`)
//	}
//
//	// Заменить пустые значения на ноль
//	if numberArray[rightPartNumber] == "" {
//		numberArray[rightPartNumber] = "0"
//	}
//	if numberArray[leftPartNumber] == "" {
//		numberArray[leftPartNumber] = "0"
//	}
//
//	if len(numberArray[leftPartNumber]) > maxIntegerPartLength {
//		// Убрать лишнюю целую часть числа
//		numberArray[leftPartNumber] = numberArray[leftPartNumber][0:maxIntegerPartLength]
//	}
//
//	if len(numberArray[rightPartNumber]) > maxIntegerPartLength {
//		// Убрать лишнюю десятичную часть числа
//		numberArray[rightPartNumber] = numberArray[rightPartNumber][0:maxIntegerPartLength]
//	}
//	return numberArray
//}
//
//func clearFromString(str string, regex string) string {
//	return replaceInString(str, regex, ``)
//}
//
//func replaceInString(str string, regex string, repl string) string {
//	return regexp.MustCompile(regex).ReplaceAllString(str, repl)
//}
