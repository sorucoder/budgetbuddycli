package quantity

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
)

var (
	moneyRegexp        = regexp.MustCompile(`^([+\-]?)\$(\d+(?:\.\d{2})?)$`)
	moneyGroupedRegexp = regexp.MustCompile(`^([+\-]?)\$(\d{1,3}(?:,\d{3})*(?:\.\d{2})?)$`)
)

// Money describes a human-friendly monetary value.
type Money float64

// NewMoney transforms the given value into a Money, if possible; otherwise, this returns an error
func NewMoney(value interface{}) (Money, error) {
	switch moneyValue := value.(type) {
	case Quantity:
		return Money(moneyValue.ValueOf()), nil
	case nil:
		return Money(math.NaN()), nil
	case int:
		return Money(float64(moneyValue)), nil
	case float64:
		return Money(moneyValue), nil
	case string:
		if moneyRegexp.MatchString(moneyValue) {
			moneyValue = strings.Join(moneyRegexp.FindStringSubmatch(moneyValue)[1:], "")
			if money, err := strconv.ParseFloat(moneyValue, 64); err == nil {
				return Money(money), nil
			} else {
				return Money(math.NaN()), fmt.Errorf(`failed to parse string %s as budget.Money: %w`, moneyValue, err)
			}
		} else if moneyGroupedRegexp.MatchString(moneyValue) {
			moneyValue = strings.ReplaceAll(strings.Join(moneyGroupedRegexp.FindStringSubmatch(moneyValue)[1:], ""), ",", "")
			if money, err := strconv.ParseFloat(moneyValue, 64); err == nil {
				return Money(money), nil
			} else {
				return Money(math.NaN()), fmt.Errorf(`failed to parse string %s as budget.Money: %w`, moneyValue, err)
			}
		} else {
			return Money(math.NaN()), fmt.Errorf(`failed to parse string %s as budget.Money: invalid format`, moneyValue)
		}
	default:
		return Money(math.NaN()), fmt.Errorf(`failed to parse %[1]T %[1]v as budget.Money: invalid type`, value)
	}
}

// MakeMoney transforms the given value into a Money; otherwise, this panics
func MakeMoney(value interface{}) Money {
	if money, err := NewMoney(value); err == nil {
		return money
	} else {
		panic(err)
	}
}

// ValueOf implements Quantity for Money
func (money Money) ValueOf() float64 {
	return float64(money)
}

// IsInf is a wrapper of math.IsInf
func (money Money) IsInf(sign int) bool {
	return math.IsInf(money.ValueOf(), sign)
}

// IsNaN is a wrapper of math.IsNaN
func (money Money) IsNaN() bool {
	return math.IsNaN(money.ValueOf())
}

// Dollars returns the value of whole dollars
func (money Money) Dollars() Money {
	dollarValue, _ := math.Modf(money.ValueOf())
	return Money(dollarValue)
}

// Cents returns the value of fractional dollars
func (money Money) Cents() Money {
	_, centValue := math.Modf(money.ValueOf())
	return Money(centValue)
}

// String implements fmt.Stringer for Money
func (money Money) String() string {
	var builder strings.Builder

	switch moneyValue := money.ValueOf(); {
	case math.IsNaN(moneyValue):
		builder.WriteString("$?")
	case math.IsInf(moneyValue, 1):
		builder.WriteString("$∞")
	case math.IsInf(moneyValue, -1):
		builder.WriteString("-$∞")
	default:
		dollarValue, centValue := math.Modf(moneyValue)
		dollarString := strconv.FormatFloat(dollarValue, 'f', 0, 64)
		centString := strings.TrimPrefix(strconv.FormatFloat(math.Abs(centValue), 'f', 2, 64), "0")

		builder.WriteRune('$')
		for index, dollarRune := range dollarString {
			if index > 0 && (len(dollarString)-index)%3 == 0 {
				builder.WriteRune(',')
			}
			builder.WriteRune(dollarRune)
		}
		builder.WriteString(centString)
	}

	return builder.String()
}

// WriteAnswer implements survey.core.Settable for Money
func (money *Money) WriteAnswer(field string, value interface{}) error {
	if moneyValue, err := NewMoney(value); err == nil {
		*money = moneyValue
	} else {
		return err
	}
	return nil
}
