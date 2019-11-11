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

// Setup initializes the zap based logger
func Setup(out io.ReadWriter) (benchmark.Setup, error) {
	if err := zapSink.SetOut(out); err != nil {
		return benchmark.Setup{}, fmt.Errorf("setting sink output: %w", err)
	}

	var outputs []string
	if out == os.Stdout {
		outputs = []string{"stdout"}
	} else {
		outputs = []string{"memory://"}
	}

	conf := zap.Config{
		Level:       zap.NewAtomicLevelAt(zap.DebugLevel),
		OutputPaths: outputs,
		Encoding:    "json",
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

	l, err := conf.Build()
	if err != nil {
		return benchmark.Setup{}, fmt.Errorf("building zap config: %w", err)
	}

	// Choose log function
	return benchmark.Setup{
		Info: func(msg string) {
			l.Info(msg)
		},
		InfoFmt: func(msg string, data int) {
			l.Info(fmt.Sprintf(msg, data))
		},
		InfoWithErrorStack: func(msg string, err error) {
			l.Info(msg, zap.Error(err))
		},
		Error: func(msg string) {
			l.Error(msg)
		},
		InfoWith3: func(msg string, fields *benchmark.Fields3) {
			l.Info(msg,
				zap.String(fields.Name1, fields.Value1),
				zap.Int(fields.Name2, fields.Value2),
				zap.Float64(fields.Name3, fields.Value3),
			)
		},
		InfoWith10: func(msg string, fields *benchmark.Fields10) {
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
		},
	}, nil
}
