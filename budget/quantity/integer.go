package quantity

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
)

var (
	integerRegexp        = regexp.MustCompile(`^([+\-]?\d+)$`)
	integerGroupedRegexp = regexp.MustCompile(`^([+\-]?\d{1,3}(?:,\d{3})*)$`)
)

// Integer describes a human-friendly mathematical integer
type Integer float64

// NewInteger transforms the given value into an Integer, if possible; otherwise, this returns an error
func NewInteger(value interface{}) (Integer, error) {
	switch typedValue := value.(type) {
	case Quantity:
		integerValue, _ := math.Modf(typedValue.ValueOf())
		return Integer(integerValue), nil
	case nil:
		return Integer(math.NaN()), nil
	case int:
		return Integer(typedValue), nil
	case int64:
		return Integer(typedValue), nil
	case float64:
		typedValue, _ = math.Modf(typedValue)
		return Integer(typedValue), nil
	case string:
		if integerRegexp.MatchString(typedValue) {
			typedValue = integerRegexp.FindStringSubmatch(typedValue)[0]
			if integer, err := strconv.ParseInt(typedValue, 10, 64); err == nil {
				return Integer(integer), nil
			} else {
				return 0, fmt.Errorf(`failed to parse string %s as budget.Integer: %w`, typedValue, err)
			}
		} else if integerGroupedRegexp.MatchString(typedValue) {
			typedValue = strings.ReplaceAll(integerRegexp.FindStringSubmatch(typedValue)[0], ",", "")
			if integer, err := strconv.ParseInt(typedValue, 10, 64); err == nil {
				return Integer(integer), nil
			} else {
				return 0, fmt.Errorf(`failed to parse string %s as budget.Integer: %w`, typedValue, err)
			}
		} else {
			return 0, fmt.Errorf(`failed to parse string %s as budget.Number: invalid format`, typedValue)
		}
	default:
		return Integer(math.NaN()), fmt.Errorf(`failed to parse %[1]T %[1]v as budget.Integer: invalid type`, value)
	}
}

// MakeInteger transforms the given value into an Integer, if possible; otherwise, this panics
func MakeInteger(value interface{}) Integer {
	if integer, err := NewInteger(value); err == nil {
		return integer
	} else {
		panic(err)
	}
}

// ValueOf implements Quantity for Integer
func (integer Integer) ValueOf() float64 {
	return float64(integer)
}

// IsInf is a wrapper of math.IsInf
func (integer Integer) IsInf(sign int) bool {
	return math.IsInf(integer.ValueOf(), sign)
}

// IsNaN is a wrapper of math.IsNaN
func (integer Integer) IsNaN() bool {
	return math.IsNaN(integer.ValueOf())
}

// Ordinal returns the ordinal number that corresponds to the value of Integer. Returns an empty string if negative.
func (integer Integer) Ordinal() string {
	var builder strings.Builder

	integerValue := integer.ValueOf()
	if integerValue >= 0 {
		builder.WriteString(strconv.FormatFloat(integerValue, 'f', 0, 64))
		if integerValue > 10 && integerValue < 20 {
			builder.WriteString("th")
		} else {
			switch math.Mod(integerValue, 10) {
			case 1:
				builder.WriteString("st")
			case 2:
				builder.WriteString("nd")
			case 3:
				builder.WriteString("rd")
			default:
				builder.WriteString("th")
			}
		}
	}

	return builder.String()
}

// String implements fmt.Stringer for Integer
func (integer Integer) String() string {
	var builder strings.Builder

	switch integerValue := integer.ValueOf(); {
	case math.IsInf(integerValue, 1):
		builder.WriteString("∞")
	case math.IsInf(integerValue, -1):
		builder.WriteString("-∞")
	case math.IsNaN(integerValue):
		builder.WriteString("?")
	default:
		integer := strconv.FormatFloat(integerValue, 'f', 0, 64)
		for index, integerRune := range integer {
			if index > 0 && (len(integer)-index)%3 == 0 {
				builder.WriteRune(',')
			}
			builder.WriteRune(integerRune)
		}
	}

	return builder.String()
}

// WriteAnswer implements survey.core.Settable for Integer
func (integer *Integer) WriteAnswer(field string, value interface{}) error {
	if integerValue, err := NewInteger(value); err == nil {
		*integer = integerValue
	} else {
		return err
	}
	return nil
}
