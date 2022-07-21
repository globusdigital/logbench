package phuslog

import (
	"io"

	"github.com/globusdigital/logbench/benchmark"
	phuslog "github.com/phuslu/log"
)

func init() {
}

func newLogger(out io.ReadWriter) phuslog.Logger {
	// Initialize logger
	return phuslog.Logger{
		Level:  phuslog.InfoLevel,
		Writer: &phuslog.IOWriter{out},
	}
}

func newInfo(out io.ReadWriter) (benchmark.FnInfo, error) {
	l := newLogger(out)
	return func(msg string) {
		l.Info().Msg(msg)
	}, nil
}

func newInfoFmt(out io.ReadWriter) (benchmark.FnInfoFmt, error) {
	l := newLogger(out)
	return func(msg string, data int) {
		l.Info().Msgf(msg, data)
	}, nil
}

func newInfoWithErrorStack(out io.ReadWriter) (
	benchmark.FnInfoWithErrorStack, error,
) {
	l := newLogger(out)
	return func(msg string, err error) {
		l.Context = phuslog.NewContext(l.Context[:0]).Err(err).Value()
		l.Info().Msg(msg)
	}, nil
}

func newError(out io.ReadWriter) (benchmark.FnError, error) {
	l := newLogger(out)
	return func(msg string) {
		l.Error().Msg(msg)
	}, nil
}

func newInfoWith3(out io.ReadWriter) (benchmark.FnInfoWith3, error) {
	l := newLogger(out)
	return func(msg string, fields *benchmark.Fields3) {
		l.Context = phuslog.NewContext(l.Context[:0]).
			Str(fields.Name1, fields.Value1).
			Int(fields.Name2, fields.Value2).
			Float64(fields.Name3, fields.Value3).
			Value()
		l.Info().Msg(msg)
	}, nil
}

func newInfoWith10(out io.ReadWriter) (benchmark.FnInfoWith10, error) {
	l := newLogger(out)
	return func(msg string, fields *benchmark.Fields10) {
		l.Context = phuslog.NewContext(l.Context[:0]).
			Str(fields.Name1, fields.Value1).
			Str(fields.Name2, fields.Value2).
			Str(fields.Name3, fields.Value3).
			Bool(fields.Name4, fields.Value4).
			Str(fields.Name5, fields.Value5).
			Int(fields.Name6, fields.Value6).
			Float64(fields.Name7, fields.Value7).
			Strs(fields.Name8, fields.Value8).
			Ints(fields.Name9, fields.Value9).
			Floats64(fields.Name10, fields.Value10).
			Value()
		l.Info().Msg(msg)
	}, nil
}

func newInfoWith10Exist(out io.ReadWriter) (
	benchmark.FnInfoWith10Exist,
	error,
) {
	fields := benchmark.NewFields10()
	l := newLogger(out)
	l.Context = phuslog.NewContext(l.Context[:0]).
		Str(fields.Name1, fields.Value1).
		Str(fields.Name2, fields.Value2).
		Str(fields.Name3, fields.Value3).
		Bool(fields.Name4, fields.Value4).
		Str(fields.Name5, fields.Value5).
		Int(fields.Name6, fields.Value6).
		Float64(fields.Name7, fields.Value7).
		Strs(fields.Name8, fields.Value8).
		Ints(fields.Name9, fields.Value9).
		Floats64(fields.Name10, fields.Value10).
		Value()
	return func(msg string) {
		l.Info().Msg(msg)
	}, nil
}

// Setup initializes the phuslog based logger
func Setup() benchmark.Setup {
	return benchmark.Setup{
		Info:               newInfo,
		InfoFmt:            newInfoFmt,
		InfoWithErrorStack: newInfoWithErrorStack,
		Error:              newError,
		InfoWith3:          newInfoWith3,
		InfoWith10:         newInfoWith10,
		InfoWith10Exist:    newInfoWith10Exist,
	}
}
