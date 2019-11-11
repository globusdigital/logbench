package main

import (
	"fmt"
	"reflect"
	"time"
)

// FV represents a mapping between field names
// and the according list validators
type FV map[string]func(interface{}) error

func expectString(actual interface{}) (string, error) {
	val, ok := actual.(string)
	if !ok {
		return "", fmt.Errorf(
			"unexpected field type (expected: string; got: %s)",
			reflect.TypeOf(actual),
		)
	}
	return val, nil
}

func expectFloat64(actual interface{}) (float64, error) {
	val, ok := actual.(float64)
	if !ok {
		return 0, fmt.Errorf(
			"unexpected field type (expected: int; got: %s)",
			reflect.TypeOf(actual),
		)
	}
	return val, nil
}

func validateTime(actual interface{}) error {
	val, err := expectString(actual)
	if err != nil {
		return err
	}

	_, err = time.Parse(time.RFC3339, val)
	return err
}

func newValidatorLevel(expected string) func(interface{}) error {
	return func(actual interface{}) error {
		val, err := expectString(actual)
		if err != nil {
			return err
		}
		if val != expected {
			return fmt.Errorf(
				"mismatching level: (expected: %q, got: %q)",
				expected,
				val,
			)
		}
		return nil
	}
}

func newValidatorText(expected string) func(interface{}) error {
	return func(actual interface{}) error {
		val, err := expectString(actual)
		if err != nil {
			return err
		}
		if val != expected {
			return fmt.Errorf(
				"mismatching text: (expected: %q, got: %q)",
				expected,
				val,
			)
		}
		return nil
	}
}

func newValidatorInt(expected int) func(interface{}) error {
	return func(actual interface{}) error {
		val, err := expectFloat64(actual)
		if err != nil {
			return err
		}
		if val != float64(expected) {
			return fmt.Errorf(
				"mismatching int: (expected: %d, got: %f)",
				expected,
				val,
			)
		}
		return nil
	}
}

func newValidatorFloat64(expected float64) func(interface{}) error {
	return func(actual interface{}) error {
		val, err := expectFloat64(actual)
		if err != nil {
			return err
		}
		if val != expected {
			return fmt.Errorf(
				"mismatching float64: (expected: %f, got: %f)",
				expected,
				val,
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
			val, err := expectString(val)
			if err != nil {
				return err
			}
			expected := expected[i]
			if expected != val {
				return fmt.Errorf(
					"mismatching array item: (expected: %s, got: %s)",
					expected,
					val,
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
			val, err := expectFloat64(val)
			if err != nil {
				return err
			}
			expected := expected[i]

			if float64(expected) != val {
				return fmt.Errorf(
					"mismatching array item: (expected: %d, got: %f)",
					expected,
					val,
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
			val, err := expectFloat64(val)
			if err != nil {
				return err
			}
			expected := expected[i]

			if float64(expected) != val {
				return fmt.Errorf(
					"mismatching array item: (expected: %f, got: %f)",
					expected,
					val,
				)
			}
		}
		return nil
	}
}
