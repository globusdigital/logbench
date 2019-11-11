package main

import (
	"fmt"
	"reflect"
	"time"
)

// FV represents a mapping between field names
// and the according list validators
type FV map[string]func(interface{}) error

func validateTime(actual interface{}) error {
	val, ok := actual.(string)
	if !ok {
		return fmt.Errorf(
			"unexpected field type (expected: string; got: %s)",
			reflect.TypeOf(actual),
		)
	}
	_, err := time.Parse(time.RFC3339, val)
	return err
}

func newValidatorLevel(expected string) func(interface{}) error {
	return func(actual interface{}) error {
		val, ok := actual.(string)
		if !ok {
			return fmt.Errorf(
				"unexpected field type (expected: string; got: %s)",
				reflect.TypeOf(actual),
			)
		}
		if val != expected {
			return fmt.Errorf(
				"mismatching level: (expected: %q, got: %q)",
				expected,
				actual,
			)
		}
		return nil
	}
}

func newValidatorText(expected string) func(interface{}) error {
	return func(actual interface{}) error {
		val, ok := actual.(string)
		if !ok {
			return fmt.Errorf(
				"unexpected field type (expected: string; got: %s)",
				reflect.TypeOf(actual),
			)
		}
		if val != expected {
			return fmt.Errorf(
				"mismatching text: (expected: %q, got: %q)",
				expected,
				actual,
			)
		}
		return nil
	}
}

func newValidatorInt(expected int) func(interface{}) error {
	return func(actual interface{}) error {
		val, ok := actual.(float64)
		if !ok {
			return fmt.Errorf(
				"unexpected field type (expected: int; got: %s)",
				reflect.TypeOf(actual),
			)
		}
		if val != float64(expected) {
			return fmt.Errorf(
				"mismatching int: (expected: %d, got: %d)",
				expected,
				actual,
			)
		}
		return nil
	}
}

func newValidatorFloat64(expected float64) func(interface{}) error {
	return func(actual interface{}) error {
		val, ok := actual.(float64)
		if !ok {
			return fmt.Errorf(
				"unexpected field type (expected: float64; got: %s)",
				reflect.TypeOf(actual),
			)
		}
		if val != expected {
			return fmt.Errorf(
				"mismatching float64: (expected: %f, got: %f)",
				expected,
				actual,
			)
		}
		return nil
	}
}

func newValidatorStrings(expected []string) func(interface{}) error {
	return func(actual interface{}) error {
		val, ok := actual.([]interface{})
		if !ok {
			return fmt.Errorf(
				"unexpected field type (expected: []string; got: %s)",
				reflect.TypeOf(actual),
			)
		}
		for i, val := range val {
			val, ok := val.(string)
			if !ok {
				return fmt.Errorf(
					"unexpected array-field item type "+
						"(expected: string; got: %s)",
					reflect.TypeOf(val),
				)
			}

			if expected[i] != val {
				return fmt.Errorf(
					"mismatching array item: (expected: %s, got: %s)",
					expected,
					actual,
				)
			}
		}
		return nil
	}
}

func newValidatorInts(expected []int) func(interface{}) error {
	return func(actual interface{}) error {
		val, ok := actual.([]interface{})
		if !ok {
			return fmt.Errorf(
				"unexpected field type (expected: []int; got: %s)",
				reflect.TypeOf(actual),
			)
		}
		for i, val := range val {
			val, ok := val.(float64)
			if !ok {
				return fmt.Errorf(
					"unexpected array-field item type "+
						"(expected: float64; got: %s)",
					reflect.TypeOf(val),
				)
			}

			if float64(expected[i]) != val {
				return fmt.Errorf(
					"mismatching array item: (expected: %d, got: %f)",
					expected,
					actual,
				)
			}
		}
		return nil
	}
}

func newValidatorFloat64s(expected []float64) func(interface{}) error {
	return func(actual interface{}) error {
		val, ok := actual.([]interface{})
		if !ok {
			return fmt.Errorf(
				"unexpected field type (expected: []interface{}; got: %s)",
				reflect.TypeOf(actual),
			)
		}
		for i, val := range val {
			val, ok := val.(float64)
			if !ok {
				return fmt.Errorf(
					"unexpected array-field item type "+
						"(expected: float64; got: %s)",
					reflect.TypeOf(val),
				)
			}

			if float64(expected[i]) != val {
				return fmt.Errorf(
					"mismatching array item: (expected: %f, got: %f)",
					expected,
					actual,
				)
			}
		}
		return nil
	}
}
