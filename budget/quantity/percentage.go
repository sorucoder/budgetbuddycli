package quantity

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
)

var (
	percentageRegexp        = regexp.MustCompile(`^([+\-]?\d+(?:\.\d+)?)%$`)
	percentageGroupedRegexp = regexp.MustCompile(`^([+\-]?\d{1,3}(?:,\d{3})*(?:\.\d+)?)%$`)
)

// Percentage describes a human-friendly percentage
type Percentage float64

// NewPercentage transforms the given value into a Percentage, if possible; otherwise, this returns an error
func NewPercentage(value interface{}) (Percentage, error) {
	switch percentageValue := value.(type) {
	case Quantity:
		return Percentage(percentageValue.ValueOf() / 100), nil
	case nil:
		return Percentage(math.NaN()), nil
	case int:
		return Percentage(float64(percentageValue) / 100), nil
	case float64:
		return Percentage(percentageValue / 100), nil
	case string:
		if percentageRegexp.MatchString(percentageValue) {
			percentageValue = strings.TrimSuffix(percentageRegexp.FindStringSubmatch(percentageValue)[0], "%")
			if percentage, err := strconv.ParseFloat(percentageValue, 64); err == nil {
				return Percentage(float64(percentage) / 100), nil
			} else {
				return Percentage(math.NaN()), fmt.Errorf(`failed to parse string %s as budget.Percentage: %w`, percentageValue, err)
			}
		} else if percentageGroupedRegexp.MatchString(percentageValue) {
			percentageValue = strings.ReplaceAll(strings.TrimSuffix(percentageGroupedRegexp.FindStringSubmatch(percentageValue)[0], "%"), ",", "")
			if percentage, err := strconv.ParseFloat(percentageValue, 64); err == nil {
				return Percentage(float64(percentage) / 100), nil
			} else {
				return Percentage(math.NaN()), fmt.Errorf(`failed to parse string %s as budget.Percentage: %w`, percentageValue, err)
			}
		} else {
			return Percentage(math.NaN()), fmt.Errorf(`failed to parse string %s as budget.Percentage: invalid format`, percentageValue)
		}
	default:
		return Percentage(math.NaN()), fmt.Errorf(`failed to parse %[1]T %[1]v as budget.Percentage: invalid type`, value)
	}
}

// MakePercentage transforms the given value into a Percentage, if possible; otherwise, this panics
func MakePercentage(value interface{}) Percentage {
	if percentage, err := NewPercentage(value); err == nil {
		return percentage
	} else {
		panic(err)
	}
}

// ValueOf implements Quantity for Percentage
func (percentage Percentage) ValueOf() float64 {
	return float64(percentage)
}

// IsInf is a wrapper of math.IsInf
func (percentage Percentage) IsInf(sign int) bool {
	return math.IsInf(percentage.ValueOf(), sign)
}

// IsNaN is a wrapper of math.IsNaN
func (percentage Percentage) IsNaN() bool {
	return math.IsNaN(percentage.ValueOf())
}

// String implements fmt.Stringer for Percentage
func (percentage Percentage) String() string {
	var builder strings.Builder

	switch percentageValue := percentage.ValueOf(); {
	case math.IsNaN(percentageValue):
		builder.WriteString("?%")
	case math.IsInf(percentageValue, 1):
		builder.WriteString("∞%")
	case math.IsInf(percentageValue, -1):
		builder.WriteString("-∞%")
	default:
		integerValue, fractionalValue := math.Modf(percentageValue * 100)
		integerString := strconv.FormatFloat(integerValue, 'f', 0, 64)
		fractionString := strconv.FormatFloat(math.Abs(fractionalValue), 'f', -1, 64)

		for index, integerRune := range integerString {
			if (len(integerString)-index)%3 == 0 {
				builder.WriteRune(',')
			}
			builder.WriteRune(integerRune)
		}
		builder.WriteString(fractionString)
		builder.WriteRune('%')
	}

	return builder.String()
}

// WriteAnswer implements survey.core.Settable for Percentage
func (percentage *Percentage) WriteAnswer(field string, value interface{}) error {
	if percentageValue, err := NewPercentage(value); err == nil {
		*percentage = percentageValue
	} else {
		return err
	}
	return nil
}
