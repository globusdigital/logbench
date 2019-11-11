package zerolog

import (
	"io"
	"logbench/benchmark"

	"github.com/rs/zerolog"
)

func newLogger(out io.ReadWriter) zerolog.Logger {
	// Initialize logger
	return zerolog.New(out).With().Timestamp().Logger()
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
		l := l.With().Err(err).Logger()
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
		l := l.With().
			Str(fields.Name1, fields.Value1).
			Int(fields.Name2, fields.Value2).
			Float64(fields.Name3, fields.Value3).
			Logger()
		l.Info().Msg(msg)
	}, nil
}

func newInfoWith10(out io.ReadWriter) (benchmark.FnInfoWith10, error) {
	l := newLogger(out)
	return func(msg string, fields *benchmark.Fields10) {
		l := l.With().
			Str(fields.Name1, fields.Value1).
			Str(fields.Name2, fields.Value2).
			Str(fields.Name3, fields.Value3).
			Str(fields.Name4, fields.Value4).
			Str(fields.Name5, fields.Value5).
			Int(fields.Name6, fields.Value6).
			Float64(fields.Name7, fields.Value7).
			Strs(fields.Name8, fields.Value8).
			Ints(fields.Name9, fields.Value9).
			Floats64(fields.Name10, fields.Value10).
			Logger()
		l.Info().Msg(msg)
	}, nil
}

// Setup initializes the zerolog based logger
func Setup() benchmark.Setup {
	return benchmark.Setup{
		Info:               newInfo,
		InfoFmt:            newInfoFmt,
		InfoWithErrorStack: newInfoWithErrorStack,
		Error:              newError,
		InfoWith3:          newInfoWith3,
		InfoWith10:         newInfoWith10,
	}
}
