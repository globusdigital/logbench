package logrus

import (
	"io"
	"logbench/benchmark"

	"github.com/sirupsen/logrus"
)

func newLogger(out io.ReadWriter) *logrus.Logger {
	l := logrus.New()
	l.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: benchmark.TimeFormat,
		FieldMap: logrus.FieldMap{
			"msg": "message",
		},
	})
	l.SetOutput(out)
	l.SetLevel(logrus.InfoLevel)
	return l
}

func newInfo(out io.ReadWriter) (benchmark.FnInfo, error) {
	l := newLogger(out)
	return func(msg string) {
		l.Info(msg)
	}, nil
}

func newInfoFmt(out io.ReadWriter) (benchmark.FnInfoFmt, error) {
	l := newLogger(out)
	return func(msg string, data int) {
		l.Infof(msg, data)
	}, nil
}

func newInfoWithErrorStack(out io.ReadWriter) (
	benchmark.FnInfoWithErrorStack,
	error,
) {
	l := newLogger(out)
	return func(msg string, err error) {
		l.WithError(err).Info(msg)
	}, nil
}

func newError(out io.ReadWriter) (benchmark.FnError, error) {
	l := newLogger(out)
	return func(msg string) {
		l.Error(msg)
	}, nil
}

func newInfoWith3(out io.ReadWriter) (benchmark.FnInfoWith3, error) {
	l := newLogger(out)
	return func(msg string, fields *benchmark.Fields3) {
		l.WithFields(logrus.Fields{
			fields.Name1: fields.Value1,
			fields.Name2: fields.Value2,
			fields.Name3: fields.Value3,
		}).Info(msg)
	}, nil
}

func newInfoWith10(out io.ReadWriter) (benchmark.FnInfoWith10, error) {
	l := newLogger(out)
	return func(msg string, fields *benchmark.Fields10) {
		l.WithFields(logrus.Fields{
			fields.Name1:  fields.Value1,
			fields.Name2:  fields.Value2,
			fields.Name3:  fields.Value3,
			fields.Name4:  fields.Value4,
			fields.Name5:  fields.Value5,
			fields.Name6:  fields.Value6,
			fields.Name7:  fields.Value7,
			fields.Name8:  fields.Value8,
			fields.Name9:  fields.Value9,
			fields.Name10: fields.Value10,
		}).Info(msg)
	}, nil
}

func newInfoWith10Exist(out io.ReadWriter) (
	benchmark.FnInfoWith10Exist,
	error,
) {
	l := newLogger(out)
	fields := benchmark.NewFields10()
	e := l.WithFields(logrus.Fields{
		fields.Name1:  fields.Value1,
		fields.Name2:  fields.Value2,
		fields.Name3:  fields.Value3,
		fields.Name4:  fields.Value4,
		fields.Name5:  fields.Value5,
		fields.Name6:  fields.Value6,
		fields.Name7:  fields.Value7,
		fields.Name8:  fields.Value8,
		fields.Name9:  fields.Value9,
		fields.Name10: fields.Value10,
	})
	return func(msg string) {
		e.Info(msg)
	}, nil
}

// Setup defines the logrus logger setup
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
