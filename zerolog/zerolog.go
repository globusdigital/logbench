package zerolog

import (
	"io"
	"logbench/benchmark"

	"github.com/rs/zerolog"
)

// Setup initializes the zerolog based logger
func Setup(out io.ReadWriter) (benchmark.Setup, error) {
	// Initialize logger
	l := zerolog.New(out).With().Timestamp().Logger()

	// Choose log function
	return benchmark.Setup{
		Info: func(msg string) {
			l.Info().Msg(msg)
		},
		InfoFmt: func(msg string, data int) {
			l.Info().Msgf(msg, data)
		},
		InfoWithErrorStack: func(msg string, err error) {
			l := l.With().Err(err).Logger()
			l.Info().Msg(msg)
		},
		Error: func(msg string) {
			l.Error().Msg(msg)
		},
		InfoWith3: func(msg string, fields *benchmark.Fields3) {
			l := l.With().
				Str(fields.Name1, fields.Value1).
				Int(fields.Name2, fields.Value2).
				Float64(fields.Name3, fields.Value3).
				Logger()
			l.Info().Msg(msg)
		},
		InfoWith10: func(msg string, fields *benchmark.Fields10) {
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
		},
	}, nil
}
