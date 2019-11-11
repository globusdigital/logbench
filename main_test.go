package main

import (
	"bytes"
	"encoding/json"
	"logbench/benchmark"
	"sync"
	"testing"

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

func TestFormat(t *testing.T) {
	fields3 := benchmark.NewFields3()
	fields10 := benchmark.NewFields10()

	fieldValidators := map[string]FV{
		benchmark.LogOperationInfo: {
			benchmark.FieldTime:  validateTime,
			benchmark.FieldLevel: newValidatorLevel(benchmark.LevelInfo),
		},
		benchmark.LogOperationInfoFmt: {
			benchmark.FieldTime:  validateTime,
			benchmark.FieldLevel: newValidatorLevel(benchmark.LevelInfo),
		},
		benchmark.LogOperationInfoWithErrorStack: {
			benchmark.FieldTime:  validateTime,
			benchmark.FieldLevel: newValidatorLevel(benchmark.LevelInfo),
			benchmark.FieldError: newValidatorText("error with stack trace"),
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
			fields10.Name4:         newValidatorText(fields10.Value4),
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
