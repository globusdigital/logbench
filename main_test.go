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
	fieldValidators := map[string]FV{
		benchmark.LogOperationInfo: {
			benchmark.FieldTime:  validateTime,
			benchmark.FieldLevel: newLevelValidator(benchmark.LevelInfo),
		},
		benchmark.LogOperationInfoFmt: {
			benchmark.FieldTime:  validateTime,
			benchmark.FieldLevel: newLevelValidator(benchmark.LevelInfo),
		},
		benchmark.LogOperationInfoWithErrorStack: {
			benchmark.FieldTime:  validateTime,
			benchmark.FieldLevel: newLevelValidator(benchmark.LevelInfo),
			benchmark.FieldError: newTextValidator("error with stack trace"),
		},
		benchmark.LogOperationError: {
			benchmark.FieldTime:    validateTime,
			benchmark.FieldLevel:   newLevelValidator(benchmark.LevelError),
			benchmark.FieldMessage: newTextValidator("error message"),
		},
		benchmark.LogOperationInfoWith3: {
			benchmark.FieldTime:  validateTime,
			benchmark.FieldLevel: newLevelValidator(benchmark.LevelInfo),
		},
		benchmark.LogOperationInfoWith10: {
			benchmark.FieldTime:  validateTime,
			benchmark.FieldLevel: newLevelValidator(benchmark.LevelInfo),
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
