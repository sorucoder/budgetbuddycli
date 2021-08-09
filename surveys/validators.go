package surveys

import (
	"errors"
	"fmt"
	"math"

	"github.com/AlecAivazis/survey/v2"
	"github.com/sorucoder/budgetbuddy/budget/quantity"
)

var (
	// Invalid input errors
	errNotBounded        = errors.New("both bounds are infinite")
	errInvalidLowerBound = errors.New("lower bound is invalid")
	errInvalidUpperBound = errors.New("upper bound is invalid")
	errMismatchedBounds  = errors.New("lower bound exceeds upper bound")

	// Non-verbose user errors
	errNotInteger    = errors.New("Value must be an integer.")
	errNotNumber     = errors.New("Value must be a number.")
	errNotMoney      = errors.New("Value must be a monetary value.")
	errNotPercentage = errors.New("Value must be a percentage.")
)

// integerValidator validates that a quantity.Integer was given
func integerValidator(answer interface{}) error {
	if _, err := quantity.NewInteger(answer); err != nil {
		return errNotInteger
	}
	return nil
}

// boundedIntegerValidator returns a survey.Validator that validates that a quantity.Integer, such that its value is within the given bounds inclusively, was given.
// To specify a lower or upper bound as infinite, use nil.
// This function panics if both bounds are infinite, the lower or upper bounds cannot be transformed into quantity.Integer, or the lower bound
// is greater than the upper bound.
func boundedIntegerValidator(lowerBoundValue interface{}, upperBoundValue interface{}) survey.Validator {
	var lowerBound, upperBound quantity.Integer

	// Evaluate lower bound
	if lowerBoundValue != nil {
		if value, err := quantity.NewInteger(lowerBoundValue); err == nil {
			lowerBound = value
		} else {
			panic(errInvalidLowerBound)
		}
	} else {
		lowerBound = quantity.Integer(math.Inf(-1))
	}

	// Evaluate upper bound
	if upperBoundValue != nil {
		if value, err := quantity.NewInteger(upperBoundValue); err == nil {
			upperBound = value
		} else {
			panic(errInvalidUpperBound)
		}
	} else {
		upperBound = quantity.Integer(math.Inf(1))
	}

	// Check bounds
	if lowerBound > upperBound {
		panic(errMismatchedBounds)
	} else if lowerBound.IsInf(-1) && upperBound.IsInf(1) {
		panic(errNotBounded)
	}

	// Generate out of bounds error
	var errOutOfBounds error
	if upperBound.IsInf(1) {
		errOutOfBounds = fmt.Errorf(`Value must be greater than %s.`, lowerBound)
	} else if lowerBound.IsInf(-1) {
		errOutOfBounds = fmt.Errorf(`Value must be less than %s.`, upperBound)
	} else {
		errOutOfBounds = fmt.Errorf(`Value must be between %s and %s.`, lowerBound, upperBound)
	}

	return func(answer interface{}) error {
		if integer, err := quantity.NewInteger(answer); err == nil {
			if integer < lowerBound || integer > upperBound {
				return errOutOfBounds
			}
		} else {
			return errNotInteger
		}
		return nil
	}
}

// numberValidator validates that a quantity.Number was given
func numberValidator(answer interface{}) error {
	if _, err := quantity.NewNumber(answer); err != nil {
		return errNotNumber
	}
	return nil
}

// boundedNumberValidator returns a survey.Validator that validates that a quantity.Number, such that its value is within the given bounds inclusively, was given.
// To specify a lower or upper bound as infinite, use nil.
// This function panics if both bounds are infinite, the lower or upper bounds cannot be transformed into quantity.Number, or the lower bound
// is greater than the upper bound.
func boundedNumberValidator(lowerBoundValue interface{}, upperBoundValue interface{}) survey.Validator {
	var lowerBound, upperBound quantity.Number

	// Evalulate lower bound
	if lowerBoundValue != nil {
		if value, err := quantity.NewNumber(lowerBoundValue); err == nil {
			lowerBound = value
		} else {
			panic(errInvalidLowerBound)
		}
	} else {
		lowerBound = quantity.Number(math.Inf(-1))
	}

	// Evaluate upper bound
	if upperBoundValue != nil {
		if value, err := quantity.NewNumber(upperBoundValue); err == nil {
			upperBound = value
		} else {
			panic(errInvalidUpperBound)
		}
	} else {
		upperBound = quantity.Number(math.Inf(1))
	}

	// Check bounds
	if lowerBound > upperBound {
		panic(errMismatchedBounds)
	} else if lowerBound.IsInf(-1) && upperBound.IsInf(1) {
		panic(errNotBounded)
	}

	// Generate out of bounds error
	var errOutOfBounds error
	if upperBound.IsInf(1) {
		errOutOfBounds = fmt.Errorf(`Value must be greater than %s.`, lowerBound)
	} else if lowerBound.IsInf(-1) {
		errOutOfBounds = fmt.Errorf(`Value must be less than %s.`, upperBound)
	} else {
		errOutOfBounds = fmt.Errorf(`Value must be between %s and %s.`, lowerBound, upperBound)
	}

	return func(answer interface{}) error {
		if number, err := quantity.NewNumber(answer); err == nil {
			if number < lowerBound || number > upperBound {
				return errOutOfBounds
			}
		} else {
			return errNotNumber
		}
		return nil
	}
}

// moneyValidator validates that a quantity.Money was given
func moneyValidator(answer interface{}) error {
	if _, err := quantity.NewMoney(answer); err != nil {
		return errNotMoney
	}
	return nil
}

// boundedMoneyValidator returns a survey.Validator that validates that a quantity.Money, such that its value is within the given bounds inclusively, was given.
// To specify a lower or upper bound as infinite, use nil.
// This function panics if both bounds are infinite, the lower or upper bounds cannot be transformed into quantity.Money, or the lower bound
// is greater than the upper bound.
func boundedMoneyValidator(lowerBoundValue interface{}, upperBoundValue interface{}) survey.Validator {
	var lowerBound, upperBound quantity.Money

	// Evaluate lower bound
	if lowerBoundValue != nil {
		if value, err := quantity.NewMoney(lowerBoundValue); err == nil {
			lowerBound = value
		} else {
			panic(errInvalidLowerBound)
		}
	} else {
		lowerBound = quantity.Money(math.Inf(-1))
	}

	// Evalulate upper bound
	if upperBoundValue != nil {
		if value, err := quantity.NewMoney(upperBoundValue); err == nil {
			upperBound = value
		} else {
			panic(errInvalidUpperBound)
		}
	} else {
		upperBound = quantity.Money(math.Inf(1))
	}

	// Check bounds
	if lowerBound > upperBound {
		panic(errMismatchedBounds)
	} else if lowerBound.IsInf(-1) && upperBound.IsInf(1) {
		panic(errNotBounded)
	}

	// Generate out of bounds error
	var errOutOfBounds error
	if upperBound.IsInf(1) {
		errOutOfBounds = fmt.Errorf(`Value must be greater than %s.`, lowerBound)
	} else if lowerBound.IsInf(-1) {
		errOutOfBounds = fmt.Errorf(`Value must be less than %s.`, upperBound)
	} else {
		errOutOfBounds = fmt.Errorf(`Value must be between %s and %s.`, lowerBound, upperBound)
	}

	return func(answer interface{}) error {
		if money, err := quantity.NewMoney(answer); err == nil {
			if money < lowerBound || money > upperBound {
				return errOutOfBounds
			}
		} else {
			return errNotMoney
		}
		return nil
	}
}

// percentageValidator validates that a quantity.Percentage was given
func percentageValidator(answer interface{}) error {
	if _, err := quantity.NewPercentage(answer); err != nil {
		return err
		// return errNotPercentage
	}
	return nil
}

// boundedPercentageValidator returns a survey.Validator that validates that a quantity.Percentage, such that its value is within the given bounds inclusively, was given.
// To specify a lower or upper bound as infinite, use nil.
// This function panics if both bounds are infinite, the lower or upper bounds cannot be transformed into quantity.Percentage, or the lower bound
// is greater than the upper bound.
func boundedPercentageValidator(lowerBoundValue interface{}, upperBoundValue interface{}) survey.Validator {
	var lowerBound, upperBound quantity.Percentage

	// Evaluate lower bound
	if lowerBoundValue != nil {
		if value, err := quantity.NewPercentage(lowerBoundValue); err == nil {
			lowerBound = value
		} else {
			panic(errInvalidLowerBound)
		}
	} else {
		lowerBound = quantity.Percentage(math.Inf(-1))
	}

	// Evaluate upper bound
	if upperBoundValue != nil {
		if value, err := quantity.NewPercentage(upperBoundValue); err == nil {
			upperBound = value
		} else {
			panic(errInvalidLowerBound)
		}
	} else {
		upperBound = quantity.Percentage(math.Inf(1))
	}

	// Check bounds
	if lowerBound > upperBound {
		panic(errMismatchedBounds)
	} else if lowerBound.IsInf(-1) && upperBound.IsInf(1) {
		panic(errNotBounded)
	}

	// Generate out of bounds error
	var errOutOfBounds error
	if upperBound.IsInf(1) {
		errOutOfBounds = fmt.Errorf(`Value must be greater than %s.`, lowerBound)
	} else if lowerBound.IsInf(-1) {
		errOutOfBounds = fmt.Errorf(`Value must be less than %s.`, upperBound)
	} else {
		errOutOfBounds = fmt.Errorf(`Value must be between %s and %s.`, lowerBound, upperBound)
	}

	return func(answer interface{}) error {
		if percentage, err := quantity.NewPercentage(answer); err == nil {
			if percentage < lowerBound || percentage > upperBound {
				return errOutOfBounds
			}
		} else {
			return errNotPercentage
		}
		return nil
	}
}
