package benchmark

import (
	"errors"
	"fmt"
	"io"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

const (
	// LevelInfo represents the name of the info-level
	LevelInfo = "info"

	// LevelError represents the name of the error-level
	LevelError = "error"

	// LevelDebug represents the name of the debug-level
	LevelDebug = "debug"

	// FieldTime represents the name of the time field
	FieldTime = "time"

	// FieldLevel represents the name of the level field
	FieldLevel = "level"

	// FieldError represents the name of the error field
	FieldError = "error"

	// FieldMessage represents the name of the message field
	FieldMessage = "message"

	// LogOperationInfo represents the name of an info-log operation
	LogOperationInfo = "info"

	// LogOperationInfoFmt represents the name of an info-log operation
	// involving formatting
	LogOperationInfoFmt = "info_fmt"

	// LogOperationInfoWithErrorStack represents the name of an info-log
	// operation involving a stack-traced error value
	LogOperationInfoWithErrorStack = "info_with_error_stack"

	// LogOperationInfoWith3 represents the name of an info-log operation
	// involving 3 newly appended fields
	LogOperationInfoWith3 = "info_with_3"

	// LogOperationInfoWith10 represents the name of an info-log operation
	// involving 10 newly appended fields
	LogOperationInfoWith10 = "info_with_10"

	// LogOperationError represents the name of an error-log operation
	LogOperationError = "error"

	// TimeFormat defines the time logging format
	TimeFormat = "2006-01-02T15:04:05.999999999-07:00"
)

// Fields3 is a list of 3 fields and their according values
type Fields3 struct {
	Name1 string
	Name2 string
	Name3 string

	Value1 string
	Value2 int
	Value3 float64
}

// Fields10 is a list of 10 fields and their according values
type Fields10 struct {
	Name1  string
	Name2  string
	Name3  string
	Name4  string
	Name5  string
	Name6  string
	Name7  string
	Name8  string
	Name9  string
	Name10 string

	Value1  string
	Value2  string
	Value3  string
	Value4  string
	Value5  string
	Value6  int
	Value7  float64
	Value8  []string
	Value9  []int
	Value10 []float64
}

// NewFields3 creates a new instance of a set of 3 fields
func NewFields3() *Fields3 {
	return &Fields3{
		Name1: "field1", Value1: "some textual value",
		Name2: "field_2_int", Value2: 42,
		Name3: "field_3_float_64", Value3: 42.5,
	}
}

// NewFields10 creates a new instance of a set of 10 fields
func NewFields10() *Fields10 {
	return &Fields10{
		Name1: "field1", Value1: "",
		Name2: "field2", Value2: "some textual value",
		Name3: "field3", Value3: "and another textual value",
		Name5: "field5", Value5: "an even longer textual value",
		Name6: "field_6_int", Value6: 42,
		Name7: "field_7_float_64", Value7: 42.5,
		Name8: "field_8_multipleStrings", Value8: []string{
			"first",
			"second",
			"third",
		},
		Name9: "field_9_multipleIntegers", Value9: []int{
			0, 1, 2, 3, 4, 5, 6, 7, 8, 9,
		},
		Name10: "field_9_multipleIntegers", Value10: []float64{
			11.5, 24.9, 99.99, 50.5001, 1000.11,
		},
	}
}

// Setup defines the callback functions for all benchmarked cases
type Setup struct {
	Info               func(msg string)
	InfoFmt            func(msg string, data int)
	Error              func(msg string)
	InfoWithErrorStack func(msg string, err error)
	InfoWith3          func(msg string, fields *Fields3)
	InfoWith10         func(msg string, fields *Fields10)
}

// SetupInit defines the setup initialization function
type SetupInit func(out io.ReadWriter) (Setup, error)

func initSetup(
	out io.ReadWriter,
	setup SetupInit,
	operation string,
) (func(), error) {
	if out == nil {
		out = os.Stdout
	}

	st, err := setup(out)
	if err != nil {
		return nil, err
	}

	switch operation {
	case LogOperationInfo:
		return func() { st.Info("information") }, nil
	case LogOperationInfoFmt:
		return func() { st.InfoFmt("information", 42) }, nil
	case LogOperationInfoWithErrorStack:
		err := errors.New("error with stack trace")
		return func() { st.InfoWithErrorStack("information", err) }, nil
	case LogOperationError:
		return func() { st.Error("error message") }, nil
	case LogOperationInfoWith3:
		fields := NewFields3()
		return func() { st.InfoWith3("information", fields) }, nil
	case LogOperationInfoWith10:
		fields := NewFields10()
		return func() { st.InfoWith10("information", fields) }, nil
	}
	return nil, fmt.Errorf("unsupported operation: %q", operation)
}

// New creates a new benchmark instance
func New(
	out io.ReadWriter,
	operation string,
	setup SetupInit,
) (*Benchmark, error) {
	fn, err := initSetup(out, setup, operation)
	if err != nil {
		return nil, fmt.Errorf("init: %w", err)
	}

	return &Benchmark{
		writeLog: fn,
	}, nil
}

// Benchmark is a log benchmark
type Benchmark struct {
	writeLog func()
}

// Statistics are the statistics of the execution of a benchmark
type Statistics struct {
	TotalLogsWritten uint64
	TotalTime        time.Duration
}

// Run runs the benchmark
func (bench *Benchmark) Run(
	target uint64,
	concurrentWriters uint,
	stopped func() bool,
) Statistics {
	if stopped == nil {
		stopped = func() bool { return false }
	}

	// Execute benchmark
	start := time.Now()
	logsWritten := uint64(0)

	wg := sync.WaitGroup{}
	wg.Add(int(concurrentWriters))
	for wk := uint(0); wk < concurrentWriters; wk++ {
		go func() {
			defer wg.Done()
			for {
				if stopped() {
					break
				}
				if atomic.AddUint64(&logsWritten, 1) > target {
					break
				}
				// Write a log
				bench.writeLog()
			}
		}()
	}
	wg.Wait()

	timeTotal := time.Since(start)

	stats := Statistics{
		TotalLogsWritten: atomic.LoadUint64(&logsWritten),
		TotalTime:        timeTotal,
	}
	stats.TotalLogsWritten -= uint64(concurrentWriters)

	return stats
}
