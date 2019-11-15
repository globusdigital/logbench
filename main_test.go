package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/globusdigital/logbench/benchmark"
	"github.com/stretchr/testify/require"
)

// SyncBuffer is a thread-safe buffer implementing the io.ReadWriter interface
type SyncBuffer struct {
	m sync.Mutex
	b bytes.Buffer
}

func (b *SyncBuffer) Read(p []byte) (n int, err error) {
	b.m.Lock()
	defer b.m.Unlock()
	return b.b.Read(p)
}

func (b *SyncBuffer) Write(p []byte) (n int, err error) {
	b.m.Lock()
	defer b.m.Unlock()
	return b.b.Write(p)
}

func (b *SyncBuffer) String() string {
	b.m.Lock()
	defer b.m.Unlock()
	return b.b.String()
}

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

func expectBool(actual interface{}) (bool, error) {
	val, ok := actual.(bool)
	if !ok {
		return false, fmt.Errorf(
			"unexpected field type (expected: bool; got: %s)",
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

func newValidatorBool(expected bool) func(interface{}) error {
	return func(actual interface{}) error {
		val, err := expectBool(actual)
		if err != nil {
			return err
		}
		if val != expected {
			return fmt.Errorf("mismatching bool: (expected: %t)", expected)
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

func TestFormat(t *testing.T) {
	fields3 := benchmark.NewFields3()
	fields10 := benchmark.NewFields10()

	fieldValidators := map[string]FV{
		benchmark.LogOperationInfo: {
			benchmark.FieldTime:    validateTime,
			benchmark.FieldLevel:   newValidatorLevel(benchmark.LevelInfo),
			benchmark.FieldMessage: newValidatorText("information"),
		},
		benchmark.LogOperationInfoFmt: {
			benchmark.FieldTime:    validateTime,
			benchmark.FieldLevel:   newValidatorLevel(benchmark.LevelInfo),
			benchmark.FieldMessage: newValidatorText("information 42"),
		},
		benchmark.LogOperationInfoWithErrorStack: {
			benchmark.FieldTime:    validateTime,
			benchmark.FieldLevel:   newValidatorLevel(benchmark.LevelInfo),
			benchmark.FieldMessage: newValidatorText("information"),
			benchmark.FieldError:   newValidatorText("error with stack trace"),
		},
		benchmark.LogOperationError: {
			benchmark.FieldTime:    validateTime,
			benchmark.FieldLevel:   newValidatorLevel(benchmark.LevelError),
			benchmark.FieldMessage: newValidatorText("error message"),
		},
		benchmark.LogOperationInfoWith3: {
			benchmark.FieldTime:    validateTime,
			benchmark.FieldLevel:   newValidatorLevel(benchmark.LevelInfo),
			benchmark.FieldMessage: newValidatorText("information"),
			fields3.Name1:          newValidatorText(fields3.Value1),
			fields3.Name2:          newValidatorInt(fields3.Value2),
			fields3.Name3:          newValidatorFloat64(fields3.Value3),
		},
		benchmark.LogOperationInfoWith10: {
			benchmark.FieldTime:    validateTime,
			benchmark.FieldLevel:   newValidatorLevel(benchmark.LevelInfo),
			benchmark.FieldMessage: newValidatorText("information"),
			fields10.Name1:         newValidatorText(fields10.Value1),
			fields10.Name2:         newValidatorText(fields10.Value2),
			fields10.Name3:         newValidatorText(fields10.Value3),
			fields10.Name4:         newValidatorBool(fields10.Value4),
			fields10.Name5:         newValidatorText(fields10.Value5),
			fields10.Name6:         newValidatorInt(fields10.Value6),
			fields10.Name7:         newValidatorFloat64(fields10.Value7),
			fields10.Name8:         newValidatorStrings(fields10.Value8),
			fields10.Name9:         newValidatorInts(fields10.Value9),
			fields10.Name10:        newValidatorFloat64s(fields10.Value10),
		},
		benchmark.LogOperationInfoWith10Exist: {
			benchmark.FieldTime:    validateTime,
			benchmark.FieldLevel:   newValidatorLevel(benchmark.LevelInfo),
			benchmark.FieldMessage: newValidatorText("information"),
			fields10.Name1:         newValidatorText(fields10.Value1),
			fields10.Name2:         newValidatorText(fields10.Value2),
			fields10.Name3:         newValidatorText(fields10.Value3),
			fields10.Name4:         newValidatorBool(fields10.Value4),
			fields10.Name5:         newValidatorText(fields10.Value5),
			fields10.Name6:         newValidatorInt(fields10.Value6),
			fields10.Name7:         newValidatorFloat64(fields10.Value7),
			fields10.Name8:         newValidatorStrings(fields10.Value8),
			fields10.Name9:         newValidatorInts(fields10.Value9),
			fields10.Name10:        newValidatorFloat64s(fields10.Value10),
		},
	}

	for loggerName, initFn := range setups {
		t.Run(loggerName, func(t *testing.T) {
			for operationName, validators := range fieldValidators {
				t.Run(operationName, func(t *testing.T) {
					buf := new(SyncBuffer)
					bench, err := benchmark.New(buf, operationName, initFn)
					require.NoError(t, err)
					stats := bench.Run(1, 1, nil)
					require.Equal(t, uint64(1), stats.TotalLogsWritten)

					var fields map[string]interface{}
					dec := json.NewDecoder(buf)
					require.NoError(
						t,
						dec.Decode(&fields),
						"decoding JSON for logger %q",
						loggerName,
					)

					for requiredField, validate := range validators {
						require.Contains(t, fields, requiredField)
						require.NoError(
							t,
							validate(fields[requiredField]),
							"invalid value for field %q of logger %q",
							requiredField,
							loggerName,
						)
					}
				})
			}
		})
	}
}
