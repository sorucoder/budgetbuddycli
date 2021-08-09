package quantity

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
)

var (
	numberRegexp        = regexp.MustCompile(`^([+\-]?\d+(?:\.\d+)?)$`)
	numberGroupedRegexp = regexp.MustCompile(`^([+\-]?\d{1,3}(?:,\d{3})*)(?:\.\d+)?$`)
)

// Number describes a human-friendly mathematical real number
type Number float64

// NewNumber transforms the given value into a Number, if possible; otherwise this returns an error
func NewNumber(value interface{}) (Number, error) {
	switch numberValue := value.(type) {
	case Quantity:
		return Number(numberValue.ValueOf()), nil
	case nil:
		return 0, nil
	case int:
		return Number(numberValue), nil
	case float64:
		return Number(numberValue), nil
	case string:
		if numberRegexp.MatchString(numberValue) {
			numberValue = numberRegexp.FindStringSubmatch(numberValue)[0]
			if number, err := strconv.ParseFloat(numberValue, 64); err == nil {
				return Number(number), err
			} else {
				return Number(math.NaN()), fmt.Errorf(`failed to parse string %s as budget.Number: %w`, numberValue, err)
			}
		} else if numberGroupedRegexp.MatchString(numberValue) {
			numberValue = strings.ReplaceAll(numberRegexp.FindStringSubmatch(numberValue)[0], ",", "")
			if number, err := strconv.ParseFloat(numberValue, 64); err == nil {
				return Number(number), nil
			} else {
				return Number(math.NaN()), fmt.Errorf(`failed to parse string %s as budget.Number: %w`, numberValue, err)
			}
		} else {
			return Number(math.NaN()), fmt.Errorf(`failed to parse string %s as budget.Number: invalid format`, numberValue)
		}
	case Integer:
		return Number(numberValue), nil
	default:
		return Number(math.NaN()), fmt.Errorf(`failed to parse %[1]T %[1]v as budget.Number: invalid type`, value)
	}
}

// MakeNumber transforms the given value into a Number, if possible; otherwise, this panics
func MakeNumber(value interface{}) Number {
	if number, err := NewNumber(value); err == nil {
		return number
	} else {
		panic(err)
	}
}

// ValueOf implements Quantity for Number
func (number Number) ValueOf() float64 {
	return float64(number)
}

// IsInf is a wrapper of math.IsInf
func (number Number) IsInf(sign int) bool {
	return math.IsInf(number.ValueOf(), sign)
}

// IsNaN is a wrapper of math.IsNaN
func (number Number) IsNaN() bool {
	return math.IsNaN(number.ValueOf())
}

// String implements fmt.Stringer for Number
func (number Number) String() string {
	var builder strings.Builder

	switch numberValue := number.ValueOf(); {
	case math.IsInf(numberValue, 1):
		builder.WriteString("∞")
	case math.IsInf(numberValue, -1):
		builder.WriteString("-∞")
	case math.IsNaN(numberValue):
		builder.WriteString("?")
	default:
		integerValue, fractionalValue := math.Modf(numberValue)
		integerString := strconv.FormatFloat(integerValue, 'f', 0, 64)
		fractionString := strings.TrimPrefix(strconv.FormatFloat(math.Abs(fractionalValue), 'f', -1, 64), "0")

		for index, integerRune := range integerString {
			if (len(integerString)-index)%3 == 0 {
				builder.WriteRune(',')
			}
			builder.WriteRune(integerRune)
		}
		builder.WriteString(fractionString)
	}

	return builder.String()
}

// WriteAnswer implements survey.core.Settable for Number
func (number *Number) WriteAnswer(field string, value interface{}) error {
	if numberValue, err := NewNumber(value); err == nil {
		*number = numberValue
	} else {
		return err
	}
	return nil
}
