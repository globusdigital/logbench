package zap

import (
	"fmt"
	"io"
	"logbench/benchmark"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func defaultConfig() zap.Config {
	return zap.Config{
		Level:    zap.NewAtomicLevelAt(zap.DebugLevel),
		Encoding: "json",
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:    "message",
			LevelKey:      "level",
			TimeKey:       "time",
			NameKey:       "name",
			CallerKey:     "caller",
			StacktraceKey: "stack",
			EncodeLevel: func(
				l zapcore.Level,
				enc zapcore.PrimitiveArrayEncoder,
			) {
				switch l {
				case zapcore.DebugLevel:
					enc.AppendString("debug")
				case zapcore.InfoLevel:
					enc.AppendString("info")
				case zapcore.WarnLevel:
					enc.AppendString("warning")
				case zapcore.ErrorLevel:
					enc.AppendString("error")
				case zapcore.DPanicLevel:
					enc.AppendString("dpanic")
				case zapcore.PanicLevel:
					enc.AppendString("panic")
				case zapcore.FatalLevel:
					enc.AppendString("fatal")
				}
			},
			EncodeTime: func(
				tm time.Time,
				enc zapcore.PrimitiveArrayEncoder,
			) {
				enc.AppendString(tm.Format(benchmark.TimeFormat))
			},
			EncodeCaller: zapcore.ShortCallerEncoder,
		},
	}
}

func newLogger(out io.ReadWriter, conf zap.Config) (*zap.Logger, error) {
	if err := zapSink.SetOut(out); err != nil {
		return nil, fmt.Errorf("setting sink output: %w", err)
	}

	if out == os.Stdout {
		conf.OutputPaths = []string{"stdout"}
	} else {
		conf.OutputPaths = []string{"memory://"}
	}

	l, err := conf.Build()
	if err != nil {
		return nil, fmt.Errorf("building zap config: %w", err)
	}
	return l, nil
}

func newInfo(out io.ReadWriter) (benchmark.FnInfo, error) {
	l, err := newLogger(out, defaultConfig())
	if err != nil {
		return nil, err
	}
	return func(msg string) {
		l.Info(msg)
	}, nil
}

func newInfoFmt(out io.ReadWriter) (benchmark.FnInfoFmt, error) {
	l, err := newLogger(out, defaultConfig())
	if err != nil {
		return nil, err
	}
	return func(msg string, data int) {
		l.Info(fmt.Sprintf(msg, data))
	}, nil
}

func newInfoWithErrorStack(out io.ReadWriter) (
	benchmark.FnInfoWithErrorStack,
	error,
) {
	l, err := newLogger(out, defaultConfig())
	if err != nil {
		return nil, err
	}
	return func(msg string, err error) {
		l.Info(msg, zap.Error(err))
	}, nil
}

func newError(out io.ReadWriter) (benchmark.FnError, error) {
	l, err := newLogger(out, defaultConfig())
	if err != nil {
		return nil, err
	}
	return func(msg string) {
		l.Error(msg)
	}, nil
}

func newInfoWith3(out io.ReadWriter) (benchmark.FnInfoWith3, error) {
	l, err := newLogger(out, defaultConfig())
	if err != nil {
		return nil, err
	}
	return func(msg string, fields *benchmark.Fields3) {
		l.Info(msg,
			zap.String(fields.Name1, fields.Value1),
			zap.Int(fields.Name2, fields.Value2),
			zap.Float64(fields.Name3, fields.Value3),
		)
	}, nil
}

func newInfoWith10(out io.ReadWriter) (benchmark.FnInfoWith10, error) {
	l, err := newLogger(out, defaultConfig())
	if err != nil {
		return nil, err
	}
	return func(msg string, fields *benchmark.Fields10) {
		l.Info(msg,
			zap.String(fields.Name1, fields.Value1),
			zap.String(fields.Name2, fields.Value2),
			zap.String(fields.Name3, fields.Value3),
			zap.String(fields.Name4, fields.Value4),
			zap.String(fields.Name5, fields.Value5),
			zap.Int(fields.Name6, fields.Value6),
			zap.Float64(fields.Name7, fields.Value7),
			zap.Strings(fields.Name8, fields.Value8),
			zap.Ints(fields.Name9, fields.Value9),
			zap.Float64s(fields.Name10, fields.Value10),
		)
	}, nil
}

// Setup defines the zap logger setup
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
