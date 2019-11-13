// Package log wraps the zerolog library making its API safer and easier to use
package log

import (
	"io"
	"runtime"
	"strconv"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

// Log is a logger instance
type Log struct{ zl zerolog.Logger }

// Context is a logger context
type Context struct{ cx zerolog.Context }

// Info logs an info-level message
func (l Log) Info(msg string) {
	l.zl.Info().Msg(msg)
}

// Error logs an error-level message
func (l Log) Error(msg string) {
	l.zl.Error().Msg(msg)
}

// Debug logs a debug-level message
func (l Log) Debug(msg string) {
	l.zl.Debug().Msg(msg)
}

// Warn logs a warn-level message
func (l Log) Warn(msg string) {
	l.zl.Warn().Msg(msg)
}

// Fatal logs a fatal-level message
func (l Log) Fatal(msg string) {
	l.zl.Fatal().Msg(msg)
}

// WithErr appends the "error" and "error.stack" fields to the logger context
func (l Log) WithErr(err error) Context {
	return Context{l.zl.With()}.WithErr(err)
}

// Str appends a field to the logger context
func (l Log) Str(fieldName string, value string) Context {
	return Context{l.zl.With()}.Str(fieldName, value)
}

// Bool appends a field to the logger context
func (l Log) Bool(fieldName string, value bool) Context {
	return Context{l.zl.With()}.Bool(fieldName, value)
}

// Uint appends a field to the logger context
func (l Log) Uint(fieldName string, value uint) Context {
	return Context{l.zl.With()}.Uint(fieldName, value)
}

// Int appends a field to the logger context
func (l Log) Int(fieldName string, value int) Context {
	return Context{l.zl.With()}.Int(fieldName, value)
}

// Uint32 appends a field to the logger context
func (l Log) Uint32(fieldName string, value uint32) Context {
	return Context{l.zl.With()}.Uint32(fieldName, value)
}

// Int32 appends a field to the logger context
func (l Log) Int32(fieldName string, value int32) Context {
	return Context{l.zl.With()}.Int32(fieldName, value)
}

// Uint64 appends a field to the logger context
func (l Log) Uint64(fieldName string, value uint64) Context {
	return Context{l.zl.With()}.Uint64(fieldName, value)
}

// Int64 appends a field to the logger context
func (l Log) Int64(fieldName string, value int64) Context {
	return Context{l.zl.With()}.Int64(fieldName, value)
}

// Float32 appends a field to the logger context
func (l Log) Float32(fieldName string, value float32) Context {
	return Context{l.zl.With()}.Float32(fieldName, value)
}

// Float64 appends a field to the logger context
func (l Log) Float64(fieldName string, value float64) Context {
	return Context{l.zl.With()}.Float64(fieldName, value)
}

// Bytes appends a field to the logger context
func (l Log) Bytes(fieldName string, value []byte) Context {
	return Context{l.zl.With()}.Bytes(fieldName, value)
}

// Strs appends a field to the logger context
func (l Log) Strs(fieldName string, value []string) Context {
	return Context{l.zl.With()}.Strs(fieldName, value)
}

// Bools appends a field to the logger context
func (l Log) Bools(fieldName string, value []bool) Context {
	return Context{l.zl.With()}.Bools(fieldName, value)
}

// Uints appends a field to the logger context
func (l Log) Uints(fieldName string, value []uint) Context {
	return Context{l.zl.With()}.Uints(fieldName, value)
}

// Ints appends a field to the logger context
func (l Log) Ints(fieldName string, value []int) Context {
	return Context{l.zl.With()}.Ints(fieldName, value)
}

// Int8s appends a field to the logger context
func (l Log) Int8s(fieldName string, value []int8) Context {
	return Context{l.zl.With()}.Int8s(fieldName, value)
}

// Uint16s appends a field to the logger context
func (l Log) Uint16s(fieldName string, value []uint16) Context {
	return Context{l.zl.With()}.Uint16s(fieldName, value)
}

// Int16s appends a field to the logger context
func (l Log) Int16s(fieldName string, value []int16) Context {
	return Context{l.zl.With()}.Int16s(fieldName, value)
}

// Uint32s appends a field to the logger context
func (l Log) Uint32s(fieldName string, value []uint32) Context {
	return Context{l.zl.With()}.Uint32s(fieldName, value)
}

// Int32s appends a field to the logger context
func (l Log) Int32s(fieldName string, value []int32) Context {
	return Context{l.zl.With()}.Int32s(fieldName, value)
}

// Uint64s appends a field to the logger context
func (l Log) Uint64s(fieldName string, value []uint64) Context {
	return Context{l.zl.With()}.Uint64s(fieldName, value)
}

// Int64s appends a field to the logger context
func (l Log) Int64s(fieldName string, value []int64) Context {
	return Context{l.zl.With()}.Int64s(fieldName, value)
}

// Float32s appends a field to the logger context
func (l Log) Float32s(fieldName string, value []float32) Context {
	return Context{l.zl.With()}.Float32s(fieldName, value)
}

// Float64s appends a field to the logger context
func (l Log) Float64s(fieldName string, value []float64) Context {
	return Context{l.zl.With()}.Float64s(fieldName, value)
}

// Str appends a field to the logger context
func (c Context) Str(fieldName string, value string) Context {
	return Context{c.cx.Str(fieldName, value)}
}

// Bool appends a field to the logger context
func (c Context) Bool(fieldName string, value bool) Context {
	return Context{c.cx.Bool(fieldName, value)}
}

// Uint appends a field to the logger context
func (c Context) Uint(fieldName string, value uint) Context {
	return Context{c.cx.Uint(fieldName, value)}
}

// Int appends a field to the logger context
func (c Context) Int(fieldName string, value int) Context {
	return Context{c.cx.Int(fieldName, value)}
}

// Uint32 appends a field to the logger context
func (c Context) Uint32(fieldName string, value uint32) Context {
	return Context{c.cx.Uint32(fieldName, value)}
}

// Int32 appends a field to the logger context
func (c Context) Int32(fieldName string, value int32) Context {
	return Context{c.cx.Int32(fieldName, value)}
}

// Uint64 appends a field to the logger context
func (c Context) Uint64(fieldName string, value uint64) Context {
	return Context{c.cx.Uint64(fieldName, value)}
}

// Int64 appends a field to the logger context
func (c Context) Int64(fieldName string, value int64) Context {
	return Context{c.cx.Int64(fieldName, value)}
}

// Float32 appends a field to the logger context
func (c Context) Float32(fieldName string, value float32) Context {
	return Context{c.cx.Float32(fieldName, value)}
}

// Float64 appends a field to the logger context
func (c Context) Float64(fieldName string, value float64) Context {
	return Context{c.cx.Float64(fieldName, value)}
}

// Bytes appends a field to the logger context
func (c Context) Bytes(fieldName string, value []byte) Context {
	return Context{c.cx.Bytes(fieldName, value)}
}

// Strs appends a field to the logger context
func (c Context) Strs(fieldName string, value []string) Context {
	return Context{c.cx.Strs(fieldName, value)}
}

// Bools appends a field to the logger context
func (c Context) Bools(fieldName string, value []bool) Context {
	return Context{c.cx.Bools(fieldName, value)}
}

// Uints appends a field to the logger context
func (c Context) Uints(fieldName string, value []uint) Context {
	return Context{c.cx.Uints(fieldName, value)}
}

// Ints appends a field to the logger context
func (c Context) Ints(fieldName string, value []int) Context {
	return Context{c.cx.Ints(fieldName, value)}
}

// Int8s appends a field to the logger context
func (c Context) Int8s(fieldName string, value []int8) Context {
	return Context{c.cx.Ints8(fieldName, value)}
}

// Uint16s appends a field to the logger context
func (c Context) Uint16s(fieldName string, value []uint16) Context {
	return Context{c.cx.Uints16(fieldName, value)}
}

// Int16s appends a field to the logger context
func (c Context) Int16s(fieldName string, value []int16) Context {
	return Context{c.cx.Ints16(fieldName, value)}
}

// Uint32s appends a field to the logger context
func (c Context) Uint32s(fieldName string, value []uint32) Context {
	return Context{c.cx.Uints32(fieldName, value)}
}

// Int32s appends a field to the logger context
func (c Context) Int32s(fieldName string, value []int32) Context {
	return Context{c.cx.Ints32(fieldName, value)}
}

// Uint64s appends a field to the logger context
func (c Context) Uint64s(fieldName string, value []uint64) Context {
	return Context{c.cx.Uints64(fieldName, value)}
}

// Int64s appends a field to the logger context
func (c Context) Int64s(fieldName string, value []int64) Context {
	return Context{c.cx.Ints64(fieldName, value)}
}

// Float32s appends a field to the logger context
func (c Context) Float32s(fieldName string, value []float32) Context {
	return Context{c.cx.Floats32(fieldName, value)}
}

// Float64s appends a field to the logger context
func (c Context) Float64s(fieldName string, value []float64) Context {
	return Context{c.cx.Floats64(fieldName, value)}
}

// WithErr appends the "error" and "error.stack" fields to the logger context
func (c Context) WithErr(err error) Context {
	if err == nil {
		return c
	}
	stackTrace, _ := findErrCause(err)
	if len(stackTrace) > 0 {
		// Append original stack trace
		ctx := c.cx.
			Str("error", err.Error()).
			Bytes("error.stack", stringifyStackTrace(stackTrace))
		return Context{ctx}
	}

	// Automatically append a stack trace
	return Context{c.cx.Err(err)}
}

// Info logs an info-level message
func (c Context) Info(msg string) {
	l := c.cx.Logger()
	l.Info().Msg(msg)
}

// Error logs an error-level message
func (c Context) Error(msg string) {
	l := c.cx.Logger()
	l.Error().Msg(msg)
}

// Debug logs a debug-level message
func (c Context) Debug(msg string) {
	l := c.cx.Logger()
	l.Debug().Msg(msg)
}

// Warn logs a warn-level message
func (c Context) Warn(msg string) {
	l := c.cx.Logger()
	l.Warn().Msg(msg)
}

// Fatal logs a fatal-level message
func (c Context) Fatal(msg string) {
	l := c.cx.Logger()
	l.Fatal().Msg(msg)
}

// findErrCause iteratively tries to find the root cause error
// and determine the deepest stack-trace available.
// It supports both Go 1.13 wrapped errors and github.com/pkg/errors
func findErrCause(err error) (stackTrace errors.StackTrace, causeErr error) {
	type wrapper interface {
		error
		Unwrap() error
	}
	type tracer interface {
		error
		StackTrace() errors.StackTrace
	}
	type causer interface {
		error
		Cause() error
	}

	for {
		if er, ok := err.(tracer); ok {
			stackTrace = er.StackTrace()
		}

		// Try to unwrap wrapped errors and find the cause
		switch er := err.(type) {
		case nil:
			// No more errors in the chain
			return
		case wrapper:
			// Unwrap Go 1.13 wrapped errors and continue
			err = er.Unwrap()
		case causer:
			// Go directly to the root cause error (github.com/pkg/errors)
			err = er.Cause()
		case error:
			// Last error in the chain
			causeErr = er
			return
		}
	}
}

// stringifyStackTrace stringifies the given error stack trace
func stringifyStackTrace(trace errors.StackTrace) []byte {
	buf := bufferPool.Get()
	defer bufferPool.Put(buf)
	for ix, frame := range trace {
		pc := uintptr(frame) - 1
		fn := runtime.FuncForPC(pc)
		file, line := fn.FileLine(pc)
		fnName := fn.Name()
		buf.WriteString(file)
		buf.WriteRune(':')
		buf.WriteString(strconv.Itoa(line))
		buf.WriteRune('\n')
		buf.WriteRune('\t')
		buf.WriteString(fnName)
		if ix+1 < len(trace) {
			// Not the last entry
			buf.WriteRune('\n')
		}
	}
	return buf.Bytes()
}

// NewLogger creates a new logger instance
func NewLogger(out io.ReadWriter) Log {
	// Initialize logger
	return Log{zerolog.New(out).With().Timestamp().Logger()}
}
