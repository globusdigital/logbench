// Package log wraps the zerolog library making its API safer and easier to use
package log

import (
	"io"
	"runtime"
	"strconv"

	"github.com/globusdigital/logbench/zerowrap/log/bufpool"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

// Log is a logger instance
type Log struct {
	zl   zerolog.Logger
	cx   zerolog.Context
	isCx bool
}

func (l Log) logger() zerolog.Logger {
	if l.isCx {
		return l.cx.Logger()
	}
	return l.zl
}

func (l Log) ctx() zerolog.Context {
	if l.isCx {
		return l.cx
	}
	return l.zl.With()
}

func logFromCtx(ctx zerolog.Context) Log {
	return Log{
		isCx: true,
		cx:   ctx,
	}
}

// Info logs an info-level message
func (l Log) Info(msg string) {
	lg := l.logger()
	lg.Info().Msg(msg)
}

// Error logs an error-level message
func (l Log) Error(msg string) {
	lg := l.logger()
	lg.Error().Msg(msg)
}

// Debug logs a debug-level message
func (l Log) Debug(msg string) {
	lg := l.logger()
	lg.Debug().Msg(msg)
}

// Warn logs a warn-level message
func (l Log) Warn(msg string) {
	lg := l.logger()
	lg.Warn().Msg(msg)
}

// Fatal logs a fatal-level message
func (l Log) Fatal(msg string) {
	lg := l.logger()
	lg.Fatal().Msg(msg)
}

// WithErr appends the "error" and "error.stack" fields to the logger context
func (l Log) WithErr(err error) Log {
	if err == nil {
		return l
	}

	ctx := l.ctx()

	stackTrace, _ := findErrCause(err)
	if len(stackTrace) > 0 {
		// Append original stack trace
		ctx = ctx.
			Str("error", err.Error()).
			Bytes("error.stack", stringifyStackTrace(stackTrace))
	} else {
		// Automatically append a stack trace
		ctx = ctx.Err(err)
	}

	return logFromCtx(ctx)
}

// Str appends a field to the logger context
func (l Log) Str(fieldName string, value string) Log {
	return logFromCtx(l.ctx().Str(fieldName, value))
}

// Bool appends a field to the logger context
func (l Log) Bool(fieldName string, value bool) Log {
	return logFromCtx(l.ctx().Bool(fieldName, value))
}

// Uint appends a field to the logger context
func (l Log) Uint(fieldName string, value uint) Log {
	return logFromCtx(l.ctx().Uint(fieldName, value))
}

// Int appends a field to the logger context
func (l Log) Int(fieldName string, value int) Log {
	return logFromCtx(l.ctx().Int(fieldName, value))
}

// Int8 appends a field to the logger context
func (l Log) Int8(fieldName string, value int8) Log {
	return logFromCtx(l.ctx().Int8(fieldName, value))
}

// Uint16 appends a field to the logger context
func (l Log) Uint16(fieldName string, value uint16) Log {
	return logFromCtx(l.ctx().Uint16(fieldName, value))
}

// Int16 appends a field to the logger context
func (l Log) Int16(fieldName string, value int16) Log {
	return logFromCtx(l.ctx().Int16(fieldName, value))
}

// Uint32 appends a field to the logger context
func (l Log) Uint32(fieldName string, value uint32) Log {
	return logFromCtx(l.ctx().Uint32(fieldName, value))
}

// Int32 appends a field to the logger context
func (l Log) Int32(fieldName string, value int32) Log {
	return logFromCtx(l.ctx().Int32(fieldName, value))
}

// Uint64 appends a field to the logger context
func (l Log) Uint64(fieldName string, value uint64) Log {
	return logFromCtx(l.ctx().Uint64(fieldName, value))
}

// Int64 appends a field to the logger context
func (l Log) Int64(fieldName string, value int64) Log {
	return logFromCtx(l.ctx().Int64(fieldName, value))
}

// Float32 appends a field to the logger context
func (l Log) Float32(fieldName string, value float32) Log {
	return logFromCtx(l.ctx().Float32(fieldName, value))
}

// Float64 appends a field to the logger context
func (l Log) Float64(fieldName string, value float64) Log {
	return logFromCtx(l.ctx().Float64(fieldName, value))
}

// Bytes appends a field to the logger context
func (l Log) Bytes(fieldName string, value []byte) Log {
	return logFromCtx(l.ctx().Bytes(fieldName, value))
}

// Strs appends a field to the logger context
func (l Log) Strs(fieldName string, value []string) Log {
	return logFromCtx(l.ctx().Strs(fieldName, value))
}

// Bools appends a field to the logger context
func (l Log) Bools(fieldName string, value []bool) Log {
	return logFromCtx(l.ctx().Bools(fieldName, value))
}

// Uints appends a field to the logger context
func (l Log) Uints(fieldName string, value []uint) Log {
	return logFromCtx(l.ctx().Uints(fieldName, value))
}

// Ints appends a field to the logger context
func (l Log) Ints(fieldName string, value []int) Log {
	return logFromCtx(l.ctx().Ints(fieldName, value))
}

// Int8s appends a field to the logger context
func (l Log) Int8s(fieldName string, value []int8) Log {
	return logFromCtx(l.ctx().Ints8(fieldName, value))
}

// Uint16s appends a field to the logger context
func (l Log) Uint16s(fieldName string, value []uint16) Log {
	return logFromCtx(l.ctx().Uints16(fieldName, value))
}

// Int16s appends a field to the logger context
func (l Log) Int16s(fieldName string, value []int16) Log {
	return logFromCtx(l.ctx().Ints16(fieldName, value))
}

// Uint32s appends a field to the logger context
func (l Log) Uint32s(fieldName string, value []uint32) Log {
	return logFromCtx(l.ctx().Uints32(fieldName, value))
}

// Int32s appends a field to the logger context
func (l Log) Int32s(fieldName string, value []int32) Log {
	return logFromCtx(l.ctx().Ints32(fieldName, value))
}

// Uint64s appends a field to the logger context
func (l Log) Uint64s(fieldName string, value []uint64) Log {
	return logFromCtx(l.ctx().Uints64(fieldName, value))
}

// Int64s appends a field to the logger context
func (l Log) Int64s(fieldName string, value []int64) Log {
	return logFromCtx(l.ctx().Ints64(fieldName, value))
}

// Float32s appends a field to the logger context
func (l Log) Float32s(fieldName string, value []float32) Log {
	return logFromCtx(l.ctx().Floats32(fieldName, value))
}

// Float64s appends a field to the logger context
func (l Log) Float64s(fieldName string, value []float64) Log {
	return logFromCtx(l.ctx().Floats64(fieldName, value))
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
	buf := bufpool.BufferPool.Get()
	defer bufpool.BufferPool.Put(buf)
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

// NewLog creates a new logger instance
func NewLog(out io.ReadWriter) Log {
	// Initialize logger
	return Log{
		zl: zerolog.New(out).With().Timestamp().Logger(),
	}
}
