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

func newLevelValidator(expected string) func(interface{}) error {
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

func newTextValidator(expected string) func(interface{}) error {
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
